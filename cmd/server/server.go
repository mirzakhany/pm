package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"proj/pkg/config"
	"proj/pkg/db"
	"proj/pkg/grpcgw"
	"proj/pkg/log"
	"proj/services"
	"syscall"
)

var (
	debugMode bool
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

	flag.BoolVar(&debugMode, "debug", false, "run in debug mode")

	ctx := cliContext()
	err := log.Init(ctx, debugMode)
	if err != nil {
		panic(err)
	}

	err = config.Init("config", "yaml", "")
	if err != nil {
		panic(err)
	}

	database, err := db.Init()
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
