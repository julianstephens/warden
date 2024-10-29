package main

import (
	"context"
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
		fmt.Println("got store")
		s, err = store.OpenStore(ctx, c.Store)
		if err != nil {
			return err
		}
	} else if c.StoreFile != "" {
		fmt.Println("got store file")
	} else {
		return fmt.Errorf("no store or store definition provided")
	}

	switch c.Resource {
	case "masterkey":
		fmt.Println("show master")
		master := s.Key()
		cMaster := *master
		cMaster.Data = master.Decrypt().Data
		warden.PPrint(cMaster)
	case "config":
		fmt.Println("show config")
	default:
		fmt.Println("nothing to show")
	}
	return nil
}
