package root

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cliContext "github.com/mazxaxz/remitly-cli/internal/context"
	"github.com/mazxaxz/remitly-cli/internal/deploy"
)

func Execute(ctx context.Context) {
	cmd := &cobra.Command{
		Use:   "remitly [COMMAND]",
		Short: "Command Line Interface (CLI)",
		Long:  "A simple exec created for recruitment purposes",
	}
	// subcommands
	cmd.AddCommand(cliContext.NewCmd())
	cmd.AddCommand(deploy.NewCmd())

	if err := cmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
