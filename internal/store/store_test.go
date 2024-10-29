package store_test

import (
	"context"
	"testing"

	"github.com/julianstephens/warden/internal/backend"
	"github.com/julianstephens/warden/internal/backend/common"
	"github.com/julianstephens/warden/internal/crypto"
	"github.com/julianstephens/warden/internal/store"
	"github.com/julianstephens/warden/internal/warden"
)

func TestInit(t *testing.T) {
	testDir := "./tmp/test"
	err := warden.EnsureDir(testDir)
	if err != nil {
		t.Fatal(err)
	}

	be, err := backend.NewBackend(common.LocalStorage, common.LocalStorageParams{Location: testDir})
	if err != nil {
		t.Fatal(err)
	}

	store, err := store.NewStore(be)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	password, err := crypto.NewRandom(12)
	if err != nil {
		t.Fatal(err)
	}

	err = store.Init(ctx, crypto.DefaultParams, string(password))
	if err != nil {
		t.Fatal(err)
	}

}
