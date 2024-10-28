package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/julianstephens/warden/internal/backend/common"
	"github.com/julianstephens/warden/internal/crypto"
	"github.com/julianstephens/warden/internal/warden"
)

type Store struct {
	conf    warden.Config
	backend common.Backend
	master  *Key
}

func NewStore(be common.Backend) (*Store, error) {
	s := &Store{backend: be}
	return s, nil
}

func (s *Store) Init(ctx context.Context, params crypto.Params, password string) error {
	conf, err := warden.CreateConfig(params.ToMap())
	if err != nil {
		return err
	}

	return s.init(ctx, password, conf)
}

func (s *Store) init(ctx context.Context, password string, config warden.Config) error {
	params, err := warden.MapToStruct[*crypto.Params](config.Params)
	if err != nil {
		return err
	}

	master, err := AddKey(ctx, s, *params, password)
	if err != nil {
		return err
	}
	s.master = master

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
