package main

import (
	"context"
	"log"
	"os"

	"github.com/yaruz/app/internal/pkg/config"

	commonApp "github.com/yaruz/app/internal/app"
	"github.com/yaruz/app/internal/app/cli"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Get()
	if err != nil {
		log.Fatalln("Can not load the config")
	}
	app := cli.New(commonApp.New(ctx, *cfg), *cfg)

	if err := app.Run(); err != nil {
		log.Fatalf("Error while cli application is running: %s", err.Error())
		os.Exit(1)
	}
}
