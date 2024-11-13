package main

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	defer cancel()

	errChan := make(chan error, 1)
	defer close(errChan)

	go show(ctx, &c.Store, &c.StoreFile, c.Resource, errChan)

	// 	<-sigs
	// 	signal.Stop(sigs)
	// 	fmt.Println("got interrupt")
	// 	os.Exit(warden.ExitCodeInterrupt)
	// 	// warden.Printf("Ctrl/Cmd+C again to quit...")
	// }()

	// <-sigCtx.Done()
	// stop()
	// os.Exit(warden.ExitCodeInterrupt)

	// return <-errChan
	return <-errChan
}

func show(ctx context.Context, storeLoc *string, storeFile *string, resource string, errChan chan<- error) {
Loop:
	for {
		var s *store.Store
		var err error

		if storeLoc != nil {
			s, err = store.OpenStore(ctx, *storeLoc)
			time.Sleep(10 * time.Second)
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
			master := s.Key()
			cMaster := *master
			cMaster.Data = master.Decrypt().Data
			warden.PPrint(cMaster)
		case "config":
			warden.PPrint(s.Config())
		default:
			errChan <- errors.New("invalid resource")
			break Loop
		}

		errChan <- nil
	}
}
