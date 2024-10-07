package store

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/julianstephens/warden/internal/crypto"
	"github.com/julianstephens/warden/internal/warden"
)

type Store struct {
	conf warden.Config
	key  Key
	loc  string
}

var (
	dirs = []string{"keys", "data"}
)

func NewStore(loc string, password string) (*Store, error) {
	conf, err := warden.CreateConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to create store config: %+v", err)
	}

	err = scaffold(loc)
	if err != nil {
		return nil, err
	}

	return &Store{conf: conf, loc: loc}, nil
}

func (s *Store) Sync() error {
	confData, err := json.Marshal(s.conf)
	if err != nil {
		return err
	}

	confId := crypto.Hash(confData)
	fmt.Print(confId)

	// if _, err := os.Stat(s.loc + "/" + "config"); err != nil {
	// 	// does not exist
	// } else {
	// 	// exists
	// }

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
