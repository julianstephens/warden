package main

import (
	"fmt"
	"os"

	"github.com/julianstephens/warden/internal/store"
	"golang.org/x/term"
)

type InitCmd struct {
	Path string `arg:"" name:"path" help:"Path to new store." type:"path"`
}

func (i *InitCmd) Run(ctx *Globals) error {
	pwd, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	store, err := store.NewStore(i.Path, string(pwd))
	if err != nil {
		return err
	}
	fmt.Print(store)

	return nil
}
