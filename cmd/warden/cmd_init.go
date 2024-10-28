package main

import (
	"context"
	"fmt"

	"github.com/julianstephens/warden/internal/backend"
	"github.com/julianstephens/warden/internal/backend/common"
	"github.com/julianstephens/warden/internal/crypto"
	"github.com/julianstephens/warden/internal/store"
)

type InitCmd struct {
	BackendType string         `required:"" short:"t" enum:"${backendTypes}" help:"The backend to create (${backendTypes})" default:"${defaultBackend}"`
	Path        string         `short:"p" type:"path" help:"The location of the encrypted backup store"`
	Params      map[string]int `help:"Argon2id params (t, m, p, T)" default:"${defaultParams}"`
}

func (i *InitCmd) Run(globals *Globals) error {
	ctx := context.Background()

	t := common.BackendTypeStringMap[i.BackendType]
	if t == common.BackendType(0) {
		return fmt.Errorf("received invalid backend type: %+v", t)
	}

	var params crypto.Params
	if i.Params != nil {
		params.P = i.Params["p"]
		params.M = i.Params["m"]
		params.T = i.Params["t"]
		params.L = i.Params["T"]
	}

	password, err := crypto.ReadPassword()
	if err != nil {
		return err
	}

	var be common.Backend
	switch t {
	case common.LocalStorage:
		if i.Path == "" {
			return fmt.Errorf("path to store must be provided for local storage backend type")
		}

		be, err = backend.NewBackend(t, common.LocalStorageParams{Location: i.Path})
		if err != nil {
			return fmt.Errorf("unable to initialize localstorage backend: %+v", err)
		}
	}

	store, err := store.NewStore(be)
	if err != nil {
		return err
	}

	err = store.Init(ctx, params, password)
	if err != nil {
		return err
	}

	return nil
}
