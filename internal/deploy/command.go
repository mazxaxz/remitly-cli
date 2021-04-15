package deploy

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	version = "1.0.0"
)

type cmdContext struct {
	app, revision  string
	count          int
	countSpecified bool
	timeout        int
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
			if err := loadContexts(cmd, args); err != nil {
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

	cmd.Flags().IntVar(&c.count, "count", 0, "The number of instances of this version of the app to deploy (optional, default: same as previous version)")
	cmd.Flags().IntVarP(&c.timeout, "wait", "w", 360, "The time in seconds to wait for successful deployment (optional, default: 360)")

	return &cmd
}

func (c *cmdContext) scanFlags(cmd *cobra.Command, _ []string) error {
	c.countSpecified = cmd.Flag("count").Changed
	return nil
}

func (c *cmdContext) run(cmd *cobra.Command, _ []string) error {
	cc := viper.GetString("CONTEXT")
	x := viper.AllSettings()
	fmt.Println(cc)
	fmt.Println(x)
	fmt.Println("deploy - run")
	return nil
}
