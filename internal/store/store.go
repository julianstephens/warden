package store

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/julianstephens/warden/internal/backend"
	"github.com/julianstephens/warden/internal/crypto"
	"github.com/julianstephens/warden/internal/warden"
	pkgerr "github.com/pkg/errors"
	"golang.org/x/term"
)

type Store struct {
	conf        warden.Config
	backendType backend.BackendType
	Backend     *backend.Backend
}

var (
	dirs = []string{"keys", "data"}
)

func NewStore(t backend.BackendType) (*Store, error) {
	conf, err := warden.CreateConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to create store config: %+v", err)
	}

	s := &Store{conf: conf}

	password, err := readPassword()
	if err != nil {
		return nil, err
	}

	AddKey(s, nil, password)

	params := backend.LocalStorageParams{
		Location: "test",
	}
	be, err := backend.NewBackend(t, params)
	if err != nil {
		return nil, err
	}

	s.backendType = t
	s.Backend = &be

	return s, nil
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

func scaffold(loc string) error {
	err := os.MkdirAll(loc, os.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to create dir at %s: %+v", loc, err)
	}

	for _, d := range dirs {
		err = os.Mkdir(loc+"/"+d, os.ModePerm)
		if err != nil {
			return fmt.Errorf("unable to scaffold repo: %+v", err)
		}
	}

	return nil
}

func readPassword() (string, error) {
	fmt.Print("Enter password for new repo: ")
	pwd, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", pkgerr.Wrap(crypto.ErrInvalidPassword, err.Error())
	}

	fmt.Print("Confirm password: ")
	confPwd, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", pkgerr.Wrap(crypto.ErrInvalidPassword, err.Error())
	}

	if len(pwd) != len(confPwd) {
		return "", pkgerr.Wrap(crypto.ErrInvalidPassword, "passwords do not match")
	}

	for i := range confPwd {
		if pwd[i] != confPwd[i] {
			return "", pkgerr.Wrap(crypto.ErrInvalidPassword, "passwords do not match")
		}
	}

	return string(pwd), nil
}
