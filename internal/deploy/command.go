package deploy

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mazxaxz/remitly-cli/pkg/optional"
	"github.com/mazxaxz/remitly-cli/pkg/remitly"
)

const (
	version = "1.0.0"
)

type cmdContext struct {
	app, revision string
	count         optional.Integer
	timeout       int
}

func NewCmd() *cobra.Command {
	var c cmdContext

	cmd := cobra.Command{
		Use:     "deploy",
		Version: version,
		Short:   "A subcommand used deployment",
		Long: `
A subcommand for deploying specified version of 
the application to the remote cloud.

Subcommand uses:
	'REMITLY_CONTEXT' - environment variable (optional, default: default)
	'./contexts.yml | $HOME/.remitly/contexts.yml' - created by 'remitly contexts --init' (required)
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if err := loadSettings(cmd, args); err != nil {
				return err
			}
			if err := c.scanFlags(cmd, args); err != nil {
				return err
			}
			return nil
		},
		RunE: c.run,
	}

	cmd.Flags().StringVarP(&c.app, "application", "a", "", "Application name to be deployed (required)")
	cmd.MarkFlagRequired("application")
	cmd.Flags().StringVar(&c.revision, "revision", "", "The version of the application to to deploy (required)")
	cmd.MarkFlagRequired("revision")

	cmd.Flags().IntVar(&c.count.Value, "replica-count", 0, "The number of instances of this version of the app to deploy (optional, default: same as previous version)")
	cmd.Flags().IntVarP(&c.timeout, "wait", "w", 360, "The time in seconds to wait for successful deployment (optional, default: 360)")

	return &cmd
}

func (c *cmdContext) scanFlags(cmd *cobra.Command, _ []string) error {
	c.count.Specified = cmd.Flag("replica-count").Changed
	return nil
}

func (c *cmdContext) run(cmd *cobra.Command, _ []string) error {
	profile := viper.GetString("PROFILE")
	if profile == "" {
		return ErrProfileVariableNotSet
	}
	pc, err := profileContextFrom(viper.AllSettings(), profile)
	if err != nil {
		return err
	}

	u, err := url.Parse(pc.http.url)
	if err != nil {
		return errors.Wrapf(err, "could not parse: '%s' url", pc.http.url)
	}
	remitlyClient := remitly.NewClient(u, pc.http.username)

	timeout, cancel := context.WithTimeout(cmd.Context(), time.Duration(c.timeout)*time.Second)
	defer cancel()

	loadBalancerName := fmt.Sprintf("%s-lb", c.app)

	original, err := snapshot(timeout, remitlyClient, loadBalancerName)
	if err != nil {
		return err
	}

	replicas := c.count.Value
	if len(original.instances) == 0 {
		if c.count.Specified && c.count.Value <= 0 {
			if c.count.Value == 0 {
				log.WithContext(cmd.Context()).WithField("replica-count", c.count.Value).
					Info("specified replica count is zero or negative, skipping")
				return nil
			}
		}
	} else {
		if !c.count.Specified {
			replicas = len(original.instances)
		}
		// TODO same version deployed
	}

	rb, err := deploy(timeout, remitlyClient, loadBalancerName, c.revision, replicas)
	if !rb && err != nil {
		return err
	}
	if rb && err != nil {
		log.WithContext(cmd.Context()).WithError(err).Error("an error has occurred while deploying")
		log.WithContext(cmd.Context()).WithField("snapshot", original).Info("rolling back...")
		if err := rollback(cmd.Context(), remitlyClient, original); err != nil {
			return errors.Wrap(err, "an error has occurred while rolling back")
		}
		log.WithContext(cmd.Context()).Info("rolling back succeeded")
		return err
	}

	result := make(chan Code)
	go c.orchestrate(timeout, remitlyClient, loadBalancerName, result)
	code := <-result

	if code == CodeSuccess {
		f := log.Fields{"app": c.app, "version": c.revision}
		log.WithContext(cmd.Context()).WithFields(f).Info("successfully deployed application")
		return nil
	}

	switch code {
	case CodeError:
		log.WithContext(cmd.Context()).Error("an error has occurred while orchestrating")
	case CodeTimeout:
		log.WithContext(cmd.Context()).Error("timeout exceeded")
	case CodeUnhealthy:
		log.WithContext(cmd.Context()).Error("service unhealthy")
	}

	log.WithContext(cmd.Context()).WithField("snapshot", original).Info("rolling back...")
	if err := rollback(cmd.Context(), remitlyClient, original); err != nil {
		return errors.Wrap(err, "an error has occurred while rolling back")
	}
	return ErrFailedDeployment
}

func deploy(ctx context.Context, rc remitly.Clienter, lb, version string, replicas int) (rollback bool, _ error) {
	if replicas <= 0 {
		ss, err := snapshot(ctx, rc, lb)
		if err != nil {
			return false, err
		}
		for _, ins := range ss.instances {
			if err := rc.DeleteInstance(ctx, lb, ins.ID); err != nil {
				f := log.Fields{"name": lb, "id": ins.ID}
				log.WithContext(ctx).WithFields(f).WithError(err).Warn("could not remove instance, skipping")
			}
		}
		return false, nil
	}

	for i := 0; i < replicas; i++ {
		if _, err := rc.CreateInstance(ctx, lb, version); err != nil {
			return true, err
		}
	}
	return false, nil
}

func rollback(ctx context.Context, rc remitly.Clienter, original Snapshot) error {
	current, err := snapshot(ctx, rc, original.loadBalancer)
	if err != nil {
		return err
	}

	create := len(original.instances)
	remove := make([]string, 0)
	for _, instance := range current.instances {
		if len(original.instances) == 0 {
			remove = append(remove, instance.ID)
			continue
		}
		if exists(original.instances, instance.ID) {
			create--
		} else {
			remove = append(remove, instance.ID)
		}
	}

	if len(original.instances) > 0 {
		previousVersion := original.instances[0].Version
		for i := 0; i < create; i++ {
			if _, err := rc.CreateInstance(ctx, original.loadBalancer, previousVersion); err != nil {
				return err
			}
		}
	}

	for _, ID := range remove {
		if err := rc.DeleteInstance(ctx, original.loadBalancer, ID); err != nil {
			return err
		}
	}

	return nil
}

func exists(src []remitly.Instance, ID string) bool {
	for _, ins := range src {
		if ins.ID == ID {
			return true
		}
	}
	return false
}
