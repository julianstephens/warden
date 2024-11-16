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

func (c *ShowCmd) Run(ctx context.Context, globals *Globals) error {
	warden.Log.Debug().Msg("ShowCmd.Run")

	ctx, cancel := context.WithCancel(ctx)
	errChan := make(chan error)

	defer func() {
		cancel()
		close(errChan)
	}()

	ctx = warden.Log.WithContext(ctx)
	go show(ctx, &c.Store, &c.StoreFile, c.Resource, errChan)

	return <-errChan
}

func show(ctx context.Context, storeLoc *string, storeFile *string, resource string, errChan chan<- error) {
Loop:
	for {
		var s *store.Store
		var err error

		if storeLoc != nil {
			s, err = store.OpenStore(ctx, *storeLoc)
			if err != nil {
				errChan <- err
				break
			}
		} else if storeFile != nil {
			fmt.Println("got store file")
		} else {
			errChan <- errors.New("no store or store definition provided")
			break
		}

		switch resource {
		case "masterkey":
			warden.Log.Debug().Msg("copying master key...")
			master := s.Key()
			cMaster := *master
			warden.Log.Debug().Msg("master key copied.")
			warden.Log.Debug().Msg("decrypting key data...")
			cMaster.Data = master.Decrypt().Data
			warden.Log.Debug().Msg("key data decrypted.")
			warden.PPrint(cMaster)
		case "config":
			warden.PPrint(s.Config())
		default:
			errChan <- errors.New("invalid resource")
			break Loop
		}

		errChan <- nil
		break
	}
}
