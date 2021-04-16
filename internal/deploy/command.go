package deploy

import (
	"context"
	"net/url"
	"time"

	"github.com/pkg/errors"
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

	c.count = optional.Integer{}
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

	deploy(timeout, remitlyClient)

	return nil
}

func deploy(ctx context.Context, rc remitly.Clienter) {

}

func rollback(ctx context.Context, rc remitly.Clienter) {

}
