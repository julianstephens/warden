package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"

	"github.com/julianstephens/warden/internal/backend/common"
	"github.com/julianstephens/warden/internal/crypto"
	"github.com/julianstephens/warden/internal/warden"
)

type Globals struct {
	Debug   debugFlag   `help:"Enable debug mode"`
	Version VersionFlag `name:"version" help:"Print version information and quit"`
}

type CLI struct {
	Globals
	Init   InitCmd   `cmd:"" help:"Create a new encrypted backup store."`
	Show   ShowCmd   `cmd:"" help:"Print resource information."`
	Backup BackupCmd `cmd:"" help:"Create a new backup of a directory."`
}

type debugFlag bool

func (d debugFlag) BeforeApply() error {
	warden.SetLog(warden.NewLog(os.Stdout, zerolog.DebugLevel, time.RFC3339))
	return nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	defer func() {
		signal.Stop(sigChan)
		cancel()
	}()

	go gracefulShutdown(sigChan, ctx, cancel)

	if err := run(ctx); err != nil {
		fmt.Println()
		warden.Log.Error().Err(err).Send()
		os.Exit(warden.ExitCodeErr)
	}
}

func run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			kongCtx, cli := createApp(ctx)
			return kongCtx.Run(&cli.Globals)
		}
	}
}

func createApp(ctx context.Context) (*kong.Context, CLI) {
	warden.SetLog(warden.NewLog(os.Stderr, zerolog.ErrorLevel, time.RFC1123))

	cli := CLI{
		Globals: Globals{
			Version: VersionFlag(Version),
		},
	}

	kongCtx := kong.Parse(&cli,
		kong.Name("warden"),
		kong.Description("A CLI for encrypted backups"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version":        Version,
			"backendTypes":   strings.Join(common.BackendTypes, ","),
			"defaultParams":  crypto.DefaultParams.String(),
			"defaultBackend": common.LocalStorage.String(),
			"resources":      strings.Join(common.Resources, ","),
		},
		kong.Bind(ctx))
	kongCtx.BindTo(ctx, (*context.Context)(nil))
	return kongCtx, cli
}

func gracefulShutdown(sigChan chan os.Signal, ctx context.Context, cancel context.CancelFunc) {
	select {
	case <-sigChan:
		fmt.Println("Interrupted. Ctrl/Cmd+C again to force...")
		cancel()
		fmt.Println("Cleaning up...")
		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancelShutdown()

		err := warden.Cleanup(shutdownCtx)
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(warden.ExitCodeInterrupt)
	case <-ctx.Done():
	}
	<-sigChan
	os.Exit(warden.ExitCodeInterrupt)
}
