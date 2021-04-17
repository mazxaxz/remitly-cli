package initialize

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	version = "1.0.0"

	contextsPath     = "$HOME/.remitly"
	contextsTemplate = `
contexts:
  - name: %s
    http:
      url: %s
      username: %s
`
)

type cmdContext struct {
	name     string
	url      string
	username string
}

func NewCmd() *cobra.Command {
	var c cmdContext
	cmd := cobra.Command{
		Use:     "initialize",
		Version: version,
		Short:   "A subcommand for managing contexts",
		RunE:    c.run,
	}

	cmd.Flags().StringVarP(&c.name, "name", "n", "default", "profile name of initialized context")
	cmd.Flags().StringVar(&c.url, "url", "default", "url for initialized context")
	cmd.Flags().StringVar(&c.username, "username", "default", "username for initialized context")

	return &cmd
}

func (c *cmdContext) run(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	if c.name == "" || c.url == "" || c.username == "" {
		return ErrFlagsNotSpecified
	}
	yml := fmt.Sprintf(contextsTemplate, c.name, c.url, c.username)
	if err := touch(contextsPath, "contexts", yml); err != nil {
		log.WithContext(ctx).WithError(err).Error("could not initialize contexts")
		return err
	}
	log.WithContext(ctx).Infof("contexts were intialized at '%s'", contextsPath)
	return nil
}

func touch(path, filename, content string) error {
	normalizedPath := os.ExpandEnv(path)
	fullFilePath := fmt.Sprintf("%s/%s.yml", normalizedPath, filename)

	if _, err := os.Stat(normalizedPath); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(normalizedPath, os.FileMode(0755)); err != nil {
				return errors.Wrapf(err, "could not create directory: '%s'", normalizedPath)
			}
		} else {
			return errors.Wrapf(err, "an error occurred while checking if '%s' exists", normalizedPath)
		}
	}

	if _, err := os.Stat(fullFilePath); err != nil {
		if os.IsNotExist(err) {
			f, err := os.OpenFile(fullFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.FileMode(0755))
			if err != nil {
				return errors.Wrapf(err, "could not create file: '%s'", fullFilePath)
			}
			if _, err := f.WriteString(content); err != nil {
				return errors.Wrapf(err, "could not write into file: '%s'", fullFilePath)
			}
			defer f.Close()
		} else {
			return errors.Wrapf(err, "an error occurred while checking if '%s' exists", fullFilePath)
		}
	}
	return nil
}
