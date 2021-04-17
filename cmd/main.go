package main

import (
	"context"
	"time"

	"github.com/mazxaxz/remitly-cli/cmd/root"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.TraceLevel)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	root.Execute(ctx)
}
