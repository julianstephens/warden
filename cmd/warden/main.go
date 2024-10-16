package main

import (
	"strings"

	"github.com/alecthomas/kong"
	"github.com/julianstephens/warden/internal/backend"
	"github.com/julianstephens/warden/internal/crypto"
)

type Globals struct {
	Debug   bool        `short:"D" help:"Enable debug mode"`
	Version VersionFlag `name:"version" help:"Print version information and quit"`
}

type CLI struct {
	Globals
	Init InitCmd `cmd:"" help:"Create a new encrypted backup store."`
}

func main() {
	cli := CLI{
		Globals: Globals{
			Version: VersionFlag(Version),
		},
	}

	ctx := kong.Parse(&cli,
		kong.Name("warden"),
		kong.Description("A CLI for encrypted backups"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version":       Version,
			"backendTypes":  strings.Join(backend.BackendTypes, ","),
			"defaultParams": crypto.DefaultParams.String(),
		})
	err := ctx.Run(&cli.Globals)
	ctx.FatalIfErrorf(err)
}
