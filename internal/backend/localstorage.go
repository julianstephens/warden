package backend

import (
	"context"
	"errors"
	"fmt"
	"os"

	pkgerr "github.com/pkg/errors"
)

type Local struct {
	backend
	location string
}

type LocalStorageParams struct {
	Params
	Location string
}

type LocalHandler struct {
}

// func (h *LocalHandler) Handle(t FileType, data []byte, password string) error { return nil }
func (h *LocalHandler) putConfig(data []byte) error { return nil }
func (h *LocalHandler) putKey(data []byte) error    { return nil }
func (h *LocalHandler) putPack(data []byte) error   { return nil }

const (
	self = LocalStorage
	name = "LocalStorage"
)

var (
	ErrInvalidStore        = errors.New("invalid local store")
	ErrStoreInitialization = errors.New("unable to intialize new store")
)

func newLocalStorage(params LocalStorageParams) (*Local, error) {
	local := &Local{
		backend:  backend{self: LocalStorage, name: name, handler: &LocalHandler{}},
		location: params.Location,
	}

	if err := scaffold(local.location); err != nil {
		return nil, err
	}

	return local, nil
}

func scaffold(storeLoc string) error {
	if storeLoc == "" {
		return pkgerr.Wrap(ErrInvalidStore, "expected valid path to store, got empty string")
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

func (l *Local) GetName() string {
	return l.GetType().String()
}

func (l *Local) GetType() BackendType {
	return self
}

func (l *Local) Sync(ctx context.Context, data []byte) error {
	return nil
}

func (l *Local) Handle(ctx context.Context, t FileType, data any) error {
	d := make([]byte, 1)

	switch t {
	case Config:
		l.handler.putConfig(d)
	case Key:
		l.handler.putKey(d)
	case Pack:
		l.handler.putPack(d)
	default:
		return fmt.Errorf("invalid file type: %s", t.String())
	}

	return nil
}
