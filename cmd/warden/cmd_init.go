package main

import (
	"context"
	"fmt"

	"github.com/julianstephens/warden/internal/backend"
	"github.com/julianstephens/warden/internal/crypto"
	"github.com/julianstephens/warden/internal/store"
	"github.com/mitchellh/mapstructure"
)

type InitCmd struct {
	BackendType string         `required:"" short:"t" enum:"${backendTypes}" help:"The backend to create (${backendTypes})"`
	Path        string         `short:"p" type:"path" help:"The location of the encrypted backup store"`
	Params      map[string]int `help:"Argon2id params (t, m, p, T)" default:""`
}

func (i *InitCmd) Run(globals *Globals) error {
	ctx := context.Background()

	t := backend.BackendTypeStringMap[i.BackendType]
	if t == backend.BackendType(0) {
		return fmt.Errorf("received invalid backend type: %+v", t)
	}

	var params crypto.Params
	if i.Params != nil {
		err := mapstructure.Decode(i.Params, &params)
		if err != nil {
			return fmt.Errorf("unable to parse argon2 params: %+v", err)
		}
	}

	password, err := crypto.ReadPassword()
	if err != nil {
		return err
	}

	var be backend.Backend
	switch t {
	case backend.LocalStorage:
		if i.Path == "" {
			return fmt.Errorf("path to store must be provided for local storage backend type")
		}

		be, err = backend.NewBackend(t, backend.LocalStorageParams{Location: i.Path})
		if err != nil {
			return fmt.Errorf("unable to initliaze localstorage backend: %+v", err)
		}
	}

	store, err := store.NewStore(be)
	if err != nil {
		return err
	}
	fmt.Println(store)

	err = store.Init(ctx, &params, password)
	if err != nil {
		return err
	}

	return nil
}
