package main

import (
	"context"
	"time"

	"github.com/mazxaxz/remitly-cli/cmd/root"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	root.Execute(ctx)
}
