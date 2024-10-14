package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/julianstephens/warden/internal/backend"
	"github.com/julianstephens/warden/internal/crypto"
	"github.com/julianstephens/warden/internal/warden"
)

type Store struct {
	conf    warden.Config
	backend backend.Backend
	master  *crypto.Key
}

func NewStore(be backend.Backend) (*Store, error) {
	s := &Store{backend: be}
	return s, nil
}

func (s *Store) Init(ctx context.Context, params *crypto.Params, password string) error {
	p := warden.DefaultIfNil[crypto.Params](params, crypto.DefaultParams)

	conf, err := warden.CreateConfig(p.ToMap())
	if err != nil {
		return err
	}

	err = s.init(ctx, password, conf)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) init(ctx context.Context, password string, config warden.Config) error {
	salt := crypto.NewSalt()
	master, err := crypto.NewSessionKey(salt)
	if err != nil {
		return err
	}
	s.master = master

	masterJson, err := json.Marshal(master)
	if err != nil {
		return fmt.Errorf("unable to convert master key to json: %+v", err)
	}

	params, err := warden.MapToStruct[crypto.Params](config.Params)
	if err != nil {
		return err
	}

	derivedUser, err := crypto.NewIDKey(warden.DefaultIfNil[crypto.Params](params, crypto.DefaultParams), password, salt)
	if err != nil {
		return err
	}

	encMaster, err := crypto.Encrypt(*derivedUser, masterJson, nil)
	if err != nil {
		return err
	}

	err = s.backend.Handle(ctx, backend.Key, encMaster)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Sync() error {
	confData, err := json.Marshal(s.conf)
	if err != nil {
		return err
	}

	confId := crypto.Hash(confData)
	fmt.Print(confId)

	return nil
}
