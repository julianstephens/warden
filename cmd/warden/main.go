package main

import (
	"context"
	"strings"

	"github.com/alecthomas/kong"

	"github.com/julianstephens/warden/internal/backend/common"
	"github.com/julianstephens/warden/internal/crypto"
)

type Globals struct {
	Debug   bool        `help:"Enable debug mode"`
	Version VersionFlag `name:"version" help:"Print version information and quit"`
}

type CLI struct {
	Globals
	Init   InitCmd   `cmd:"" help:"Create a new encrypted backup store."`
	Show   ShowCmd   `cmd:"" help:"Print resource information."`
	Backup BackupCmd `cmd:"" help:"Create a new backup of a directory"`
}

func main() {
	ctx := context.Background()

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
	kongCtx.FatalIfErrorf(kongCtx.Run(&cli.Globals))
}
