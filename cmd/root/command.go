package root

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/mazxaxz/remitly-cli/internal/deploy"
	"github.com/mazxaxz/remitly-cli/internal/initialize"
)

func Execute(ctx context.Context) {
	cmd := &cobra.Command{
		Use:   "remitly [COMMAND]",
		Short: "Command Line Interface (CLI)",
		Long:  "A simple exec created for recruitment purposes",
	}
	// subcommands
	cmd.AddCommand(initialize.NewCmd())
	cmd.AddCommand(deploy.NewCmd())

	now := time.Now()
	defer func() {
		took := time.Since(now)
		log.WithField("milliseconds", took.Milliseconds()).Info("command finished")
	}()

	if err := cmd.ExecuteContext(ctx); err != nil {
		log.WithError(err).Errorln("a runtime error has occurred")
		log.Exit(1)
	}
}
