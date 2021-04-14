package context

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	contextsPath     = "$HOME/.remitly"
	contextsFileName = "contexts"
)

const contextDefaultYml = `
contexts:
  - name: default
    http:
      url: https://XXXX/
      username: XXXX
`

type cmdContext struct {
	init bool
	use  string
}

func NewCmd() *cobra.Command {
	var c cmdContext
	cmd := cobra.Command{
		Use:   "context",
		Short: "A subcommand for managing contexts",
		Run:   c.run,
	}

	cmd.Flags().StringVar(&c.use, "use", "default", "Sets REMITLY_CONTEXT env variable for config usage")
	cmd.Flags().BoolVar(&c.init, "init", false, "Initializes new context file")

	return &cmd
}

func (c *cmdContext) run(cmd *cobra.Command, args []string) {
	if c.init {
		// TODO: set them as env
		viper.SetConfigName(contextsFileName)
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath(contextsPath)

		// maybe just write it into file
		if err := touch(); err != nil {

		}
		if err := viper.ReadConfig(bytes.NewBuffer([]byte(contextDefaultYml))); err != nil {
			fmt.Println(err)
			// log
		}
		if err := viper.WriteConfig(); err != nil {
			fmt.Println(err)
			// log
		}
	}

	if c.use == "" {
		// log
	}
	// ...
}

func touch() error {
	fullContextsPath := strings.Replace(contextsPath, "$HOME", os.Getenv("HOME"), -1)
	fullFilePath := fmt.Sprintf("%s/%s.yml", fullContextsPath, contextsFileName)

	if _, err := os.Stat(fullContextsPath); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(fullContextsPath, os.FileMode(0755)); err != nil {
				return errors.Wrapf(err, "could not create directory: '%s'", fullContextsPath)
			}
			f, err := os.Create(fullFilePath)
			if err != nil {
				return errors.Wrapf(err, "could not create file: '%s'", fullFilePath)
			}
			defer f.Close()
			return nil
		} else {
			return errors.Wrapf(err, "an error occurred while checking if '%s' exists", fullContextsPath)
		}
	}

	if _, err := os.Stat(fullFilePath); err != nil {
		if os.IsNotExist(err) {
			f, err := os.OpenFile(fullFilePath, os.O_CREATE, os.FileMode(0755))
			if err != nil {
				return errors.Wrapf(err, "could not create file: '%s'", fullFilePath)
			}
			defer f.Close()
		} else {
			return errors.Wrapf(err, "an error occurred while checking if '%s' exists", fullFilePath)
		}
	}
	return nil
}
