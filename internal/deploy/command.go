package deploy

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type cmdContext struct {
	app, version string
	count        int32
	timeout      int32
}

func NewCmd() *cobra.Command {
	var c cmdContext
	cmd := cobra.Command{
		Use:   "deploy",
		Short: "A subcommand used deployment",
		Long: `
A subcommand for deploying specified version of 
the application to the remote cloud.

Subcommand uses:
	'REMITLY_CONTEXT' - environment variable (optional, default: default)
	'./contexts.yml | $HOME/.remitly/contexts.yml' - file defining usable contexts (required)
`,
		Args: args,
		Run:  c.run,
	}

	cmd.Flags().StringVarP(&c.app, "application", "a", "", "Application name to be deployed (required)")
	cmd.MarkFlagRequired("application")
	cmd.Flags().StringVarP(&c.version, "version", "", "", "The version of the application to to deploy (required)")
	cmd.MarkFlagRequired("version")

	cmd.Flags().Int32Var(&c.count, "count", 0, "The number of instances of this version of the app to deploy (optional, default: 0 - same as previous version)")
	cmd.Flags().Int32VarP(&c.timeout, "wait", "w", 360, "The time in seconds to wait for successful deployment (optional, default: 360)")

	return &cmd
}

func init() {
	cobra.OnInitialize(loadContexts)
}

func args(cmd *cobra.Command, args []string) error {
	fmt.Println("deploy - args")
	return nil
}

func (c *cmdContext) run(cmd *cobra.Command, args []string) {
	_ = viper.AllSettings()
	fmt.Println("deploy - run")
}
