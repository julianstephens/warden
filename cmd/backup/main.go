package main

import "github.com/alecthomas/kong"

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
			Version: Version,
		},
	}

	ctx := kong.Parse(&cli,
		kong.Name("backup"),
		kong.Description("A CLI for encrypted backups"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": string(Version),
		})
	err := ctx.Run(&cli.Globals)
	ctx.FatalIfErrorf(err)
}
