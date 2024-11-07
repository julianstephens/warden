package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/julianstephens/warden/internal/warden"
)

type BackupCmd struct {
	CommonFlags
	Dir    string `arg:"" type:"existingDir" help:"Path to the directory to backup"`
	DryRun bool   `short:"d" help:"Print backup results with no write."`
}

func (c *BackupCmd) Run(ctx context.Context, globals *Globals) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT)
	defer func() {
		fmt.Println()
		warden.Printf("Ctrl/Cmd+C again to quit...")
		<-sigs
		fmt.Println("Interrupted. Stopping.")
		os.Exit(warden.ExitCodeInterrupt)
	}()

	// TODO: do backup

	return nil
}
