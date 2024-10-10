package main

import (
	"fmt"
	"os"

	"github.com/julianstephens/warden/internal/backend"
	"github.com/julianstephens/warden/internal/store"
	"golang.org/x/term"
)

type InitCmd struct {
	BackendType string `required:"" short:"t" enum:"${backendTypes}" help:"The backend to create"`
}

func (i *InitCmd) Run(ctx *Globals) error {
	pwd, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	t := backend.BackendTypeStringMap[i.BackendType]
	if t == backend.BackendType(0) {
		return fmt.Errorf("received invalid backend type: %+v", t)
	}

	store, err := store.NewStore(t)
	if err != nil {
		return err
	}
	fmt.Print(store)

	return nil
}
