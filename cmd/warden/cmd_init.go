package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/julianstephens/warden/internal/backend"
	"github.com/julianstephens/warden/internal/backend/common"
	"github.com/julianstephens/warden/internal/crypto"
	"github.com/julianstephens/warden/internal/store"
	"github.com/julianstephens/warden/internal/warden"
)

type InitCmd struct {
	BackendType string         `required:"" short:"t" enum:"${backendTypes}" help:"The backend to create (${backendTypes})" default:"${defaultBackend}"`
	Store       string         `short:"s" type:"path" help:"The location of the encrypted backup store"`
	Params      map[string]int `help:"Argon2id params (t, m, p, T)" default:"${defaultParams}"`
}

func (c *InitCmd) Run(ctx context.Context, globals *Globals) error {
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

	t := common.BackendTypeStringMap[c.BackendType]
	if t == common.BackendType(0) {
		return fmt.Errorf("received invalid backend type: %+v", t)
	}

	var params crypto.Params
	if c.Params != nil {
		params.P = c.Params["p"]
		params.M = c.Params["m"]
		params.T = c.Params["t"]
		params.L = c.Params["T"]
	}

	password, err := crypto.ReadPassword()
	if err != nil {
		return err
	}

	var be common.Backend
	switch t {
	case common.LocalStorage:
		if c.Store == "" {
			return fmt.Errorf("path to store must be provided for local storage backend type")
		}

		be, err = backend.NewBackend(t, common.LocalStorageParams{Location: c.Store})
		if err != nil {
			return fmt.Errorf("unable to initialize localstorage backend: %+v", err)
		}
	}

	store := store.NewStore(be, c.Store)

	err = store.Init(ctx, params, password)
	if err != nil {
		return err
	}

	return nil
}
