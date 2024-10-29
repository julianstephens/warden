package store

import (
	"context"
	"fmt"
	"path"

	"github.com/julianstephens/warden/internal/backend"
	"github.com/julianstephens/warden/internal/backend/common"
	"github.com/julianstephens/warden/internal/crypto"
	"github.com/julianstephens/warden/internal/warden"
)

type Store struct {
	conf    warden.Config
	backend common.Backend
	master  *Key
}

func NewStore(be common.Backend) *Store {
	return &Store{backend: be}
}

func OpenStore(ctx context.Context, storeLoc string) (*Store, error) {
	be, err := backend.NewBackend(common.LocalStorage, common.LocalStorageParams{Location: storeLoc})
	if err != nil {
		return nil, fmt.Errorf("unable to initialize localstorage backend: %+v", err)
	}

	s := NewStore(be)

	err = s.open(ctx, storeLoc)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) open(ctx context.Context, storeLoc string) (err error) {
	password, err := crypto.ReadPassword()
	if err != nil {
		return
	}

	config, err := warden.LoadJSON[warden.Config](path.Join(storeLoc, "config.json"))
	if err != nil {
		return
	}
	s.conf = config

	params, err := warden.MapToStruct[crypto.Params](config.Params)
	if err != nil {
		return
	}

	master, err := LoadKey(ctx, s, params, password)
	if err != nil {
		return
	}
	s.master = master

	return
}

func (s *Store) Init(ctx context.Context, params crypto.Params, password string) error {
	conf, err := warden.CreateConfig(params.ToMap())
	if err != nil {
		return err
	}

	return s.init(ctx, password, conf)
}

func (s *Store) init(ctx context.Context, password string, config warden.Config) (err error) {
	params, err := warden.MapToStruct[*crypto.Params](config.Params)
	if err != nil {
		return
	}

	master, err := AddKey(ctx, s, *params, password)
	if err != nil {
		return
	}
	s.master = master

	return
}

func (s *Store) Key() *Key {
	return s.master
}
