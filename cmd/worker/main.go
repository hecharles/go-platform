package main

import (
	"flag"
	log "log/slog"
	"os"
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

	log.Info("Starting worker")

	return exitCodeSuccess, nil

}
