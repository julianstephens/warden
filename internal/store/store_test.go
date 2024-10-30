package store_test

import (
	"context"
	"os"
	"path"
	"testing"

	mp "github.com/agiledragon/gomonkey/v2"

	"github.com/julianstephens/warden/internal/backend"
	"github.com/julianstephens/warden/internal/backend/common"
	"github.com/julianstephens/warden/internal/crypto"
	"github.com/julianstephens/warden/internal/store"
	"github.com/julianstephens/warden/internal/warden"
)

const (
	testDir = "./tmp/test"
	testPwd = "testsecurepassword123"
)

func resetStore(t *testing.T) {
	os.RemoveAll(testDir)
	err := warden.EnsureDir(testDir)
	if err != nil {
		t.Fatal(err)
	}
}

func createAndInitStore(ctx context.Context, t *testing.T) *store.Store {
	be, err := backend.NewBackend(common.LocalStorage, common.LocalStorageParams{Location: testDir})
	if err != nil {
		t.Fatal(err)
	}
	store := store.NewStore(be)

	err = store.Init(ctx, crypto.DefaultParams, testPwd)
	if err != nil {
		t.Fatal(err)
	}

	return store
}

func assertSliceEqual[T comparable](t *testing.T, expected []T, actual []T) {
	t.Helper()
	if len(expected) != len(actual) {
		t.Errorf("expected (%+v) is not equal to actual (%+v): len(expected)=%d len(actual)=%d",
			expected, actual, len(expected), len(actual))
	}
	for i := range expected {
		if expected[i] != actual[i] {
			t.Errorf("expected[%d] (%+v) is not equal to actual[%d] (%+v)",
				i, expected[i],
				i, actual[i])
		}
	}
}

func TestInit(t *testing.T) {
	resetStore(t)

	ctx := context.Background()
	createAndInitStore(ctx, t)

	conf, err := warden.LoadJSON[warden.Config](path.Join(testDir, "config.json"))
	if err != nil {
		t.Fatalf("expected config at %s, got err: %+v", path.Join(testDir, "config.json"), err)
	}

	if conf.ID == "" {
		t.Fatal("expected config id, got empty string")
	}

	if conf.Params == nil {
		t.Fatal("expected config params, got nil")
	}

	if _, err := os.Stat(path.Join(testDir, "keys")); os.IsNotExist(err) {
		t.Fatal("key directory not found")
	}

	keys, _ := os.ReadDir(path.Join(testDir, "keys"))
	if len(keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(keys))
	}
}

func TestOpen(t *testing.T) {
	resetStore(t)

	ctx := context.Background()
	original := createAndInitStore(ctx, t)

	patches := mp.ApplyFuncReturn(crypto.ReadPassword, testPwd, nil)
	defer patches.Reset()

	opened, err := store.OpenStore(ctx, testDir)
	if err != nil {
		t.Fatal(err)
	}

	originalKey := original.Key()
	openedKey := opened.Key()

	if originalKey.Hostname != openedKey.Hostname {
		t.Fatalf("expected hostname %s, got %s", originalKey.Hostname, openedKey.Hostname)
	}

	if originalKey.Username != openedKey.Username {
		t.Fatalf("expected user %s, got %s", originalKey.Username, openedKey.Username)
	}

	if !openedKey.CreatedAt.After(originalKey.CreatedAt) {
		t.Fatalf("expected new key creation stamp after original: %s, %s", openedKey.CreatedAt, originalKey.CreatedAt)
	}

	if string(original.Key().Decrypt().Data) == string(opened.Key().Decrypt().Data) {
		t.Fatalf("expected decrypted key %s, got %s", string(original.Key().Decrypt().Data), string(opened.Key().Decrypt().Data))
	}
}
