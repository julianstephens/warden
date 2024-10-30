package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/julianstephens/warden/internal/store"
	"github.com/julianstephens/warden/internal/warden"
)

type ShowCmd struct {
	CommonFlags
	Resource string `arg:"" enum:"${resources}" help:"the resource to show (${resources})"`
}

func (c *ShowCmd) Run(globals *Globals) error {
	ctx := context.Background()

	var s *store.Store
	var err error

	if c.Store != "" {
		s, err = store.OpenStore(ctx, c.Store)
		if err != nil {
			return err
		}
	} else if c.StoreFile != "" {
		fmt.Println("got store file")
	} else {
		return errors.New("no store or store definition provided")
	}

	switch c.Resource {
	case "masterkey":
		master := s.Key()
		cMaster := *master
		cMaster.Data = master.Decrypt().Data
		warden.PPrint(cMaster)
	case "config":
		warden.PPrint(s.Config())
	default:
		return errors.New("invalid resource")
	}
	return nil
}
