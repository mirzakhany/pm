package main

import (
	"context"
	"flag"
	"github.com/mirzakhany/pm/pkg/config"
	"github.com/mirzakhany/pm/pkg/db"
	"github.com/mirzakhany/pm/pkg/grpcgw"
	"github.com/mirzakhany/pm/pkg/log"
	"github.com/mirzakhany/pm/services"
	"os"
	"os/signal"
	"syscall"
)

var (
	debugMode = flag.Bool("config", false, "run in debug mode")
)

func cliContext() context.Context {
	signals := []os.Signal{syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGABRT}
	var sig = make(chan os.Signal, len(signals))
	ctx, cancel := context.WithCancel(context.Background())
	signal.Notify(sig, signals...)
	go func() {
		<-sig
		cancel()
	}()
	return ctx
}

func main() {
	flag.Parse()

	ctx := cliContext()
	err := log.Init(ctx, *debugMode)
	if err != nil {
		panic(err)
	}

	err = config.Init("config", "yaml", "")
	if err != nil {
		panic(err)
	}

	database, err := db.Init(ctx)
	if err != nil {
		panic(err)
	}

	err = services.Setup(database)
	if err != nil {
		panic(err)
	}

	err = grpcgw.Init()
	if err != nil {
		panic(err)
	}

	err = grpcgw.Serve(ctx)
	if err != nil {
		log.Error("Serve failed with an error", log.Err(err))
		panic(err)
	}
	os.Exit(0)
}
