package store

import (
	"fmt"

	"github.com/julianstephens/warden/internal/config"
)

type Store struct {
	conf config.Config
}

func NewStore(loc string) (Store, error) {
	conf, err := config.CreateConfig()
	if err != nil {
		return Store{}, fmt.Errorf("unable to create store config: %+v", err)
	}

	store := Store{conf: conf}

	return store, nil
}
