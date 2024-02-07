package main

import (
	"context"
	"flag"
	log "log/slog"
	"os"

	"go-platform/internal/web"
	"go-platform/pkg/server"
)

type args struct {
	config string
}

func main() {

	args := args{}
	flag.StringVar(&args.config, "config", "", "configuration path")
	flag.Parse()

	exitCode, err := run(args)
	if err != nil {
		log.Error("%v", err)
	}
	os.Exit(exitCode)

}

func run(arg args) (int, error) {

	const (
		exitCodeSuccess = 0
		exitCodeError   = 1
	)

	var (
		httpServer *server.Server
	)

	config := server.Config{
		Address: "127.0.0.1:5000",
	}

	httpServer, err := server.New(config)
	if err != nil {
		log.Error("failed to create http server: %v", err)
		return exitCodeError, err
	}

	webApp := web.New()

	webApp.AttachHandler(httpServer)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := httpServer.Start(ctx); err != nil {
		return exitCodeError, err
	}

	// wait until cancel signal is received, shutdown the server
	<-ctx.Done()
	if err := httpServer.Shutdown(ctx); err != nil {
		return exitCodeError, err
	}

	return exitCodeSuccess, nil

}
