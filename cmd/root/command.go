package root

import (
	"context"

	log "github.com/sirupsen/logrus"
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
		log.WithError(err).Errorln("a runtime error has occurred")
		log.Exit(1)
	}
}
