package local

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/julianstephens/warden/internal/backend/common"
	"github.com/julianstephens/warden/internal/warden"
	pkgerr "github.com/pkg/errors"
)

type Local struct {
	common.WardenBackend
	location string
}

const (
	name = "LocalStorage"
)

var (
	ErrStoreInitialization = errors.New("unable to initialize new store")
)

func NewLocalStorage(params common.LocalStorageParams) (*Local, error) {
	local := &Local{
		WardenBackend: common.WardenBackend{Self: common.LocalStorage, Name: name, Handler: &LocalHandler{}},
		location:      params.Location,
	}

	if err := scaffold(local.location); err != nil {
		return nil, err
	}

	return local, nil
}

func makeReadonly(filename string) error {
	err := os.Chmod(filename, 0444)
	if err != nil {
		return fmt.Errorf("unable to make file read-only: %+v", err)
	}
	return nil
}

func scaffold(storeLoc string) error {
	if storeLoc == "" {
		return &warden.InvalidStoreError{Msg: "expected valid path to store, got empty string"}
	}

	// create directory
	if _, err := os.Stat(storeLoc); err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(storeLoc, os.ModePerm); err != nil {
				return pkgerr.Wrap(ErrStoreInitialization, err.Error())
			}
		} else {
			return pkgerr.Wrap(ErrStoreInitialization, err.Error())
		}
	}

	return nil
}

func (l *Local) Put(ctx context.Context, event common.Event, reader common.IReader) error {
	warden.Printf("got %s event (%s)", event.Type.String(), event.Name)

	switch event.Type {
	case common.Key:
		ctx := context.WithValue(ctx, "location", l.location)
		return l.WardenBackend.Handler.WriteKey(ctx, fmt.Sprintf("%s.json", event.Name), reader)
	case common.Config:
		// TODO: handle config save
	case common.Pack:
		// TODO: handle pack save
	default:
		return fmt.Errorf("got invalid event type: %s", event.Type.String())
	}

	return nil
}
