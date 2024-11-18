package local_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/julianstephens/warden/internal/backend/common"
	"github.com/julianstephens/warden/internal/backend/local"
	"github.com/julianstephens/warden/internal/warden"
)

var testParams = common.LocalStorageParams{
	Location: "/tmp/local",
}

func TestNewLocalStorage(t *testing.T) {
	_, err := local.NewLocalStorage(common.LocalStorageParams{Location: ""})
	if err == nil {
		t.Fatal("should error on empty localstorage location")
	}

	local, err := local.NewLocalStorage(testParams)
	if err != nil {
		t.Fatal(err)
	}

	if local.Self != common.LocalStorage {
		t.Fatalf("expected backend type %s, got %s", common.LocalStorage.String(), local.Self.String())
	}
	if local.Name != "LocalStorage" {
		t.Fatalf("expected backend named LocalStorage, got %s", local.Name)
	}

	if _, err = os.Stat(testParams.Location); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("should create store %s", testParams.Location)
	}

	os.RemoveAll(testParams.Location)
}

func TestWriteConfig(t *testing.T) {
	l, err := local.NewLocalStorage(testParams)
	if err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatalf("unable to generate json data: %+v", err)
	}

	var testConf warden.Config
	err = gofakeit.Struct(&testConf)
	if err != nil {
		t.Fatalf("unable to generate test config: %+v", err)
	}

	data, err := json.Marshal(testConf)
	reader := common.NewByteReader(data)

	ctx := context.Background()
	err = l.Handler.WriteConfig(ctx, reader)
	if err == nil || !errors.Is(err, local.ErrNoStoreLocation) {
		t.Fatal("should error on context with no store location")
	}

	ctx = context.WithValue(ctx, "location", testParams.Location)

	err = l.Handler.WriteConfig(ctx, reader)
	if err != nil {
		t.Fatal(err)
	}

	data, err = os.ReadFile(path.Join(testParams.Location, "config.json"))
	if err != nil {
		t.Fatal(err)
	}

	var conf warden.Config
	err = json.Unmarshal(data, &conf)
	if err != nil {
		t.Fatal("unable to unmarshal stored data to warden config")
	}

	if testConf.ID != conf.ID {
		t.Fatalf("expected config id %s, got %s", testConf.ID, conf.ID)
	}

	if eq := reflect.DeepEqual(testConf.Params, conf.Params); !eq {
		t.Fatal("expected stored config params to equal test config params")
	}

	os.RemoveAll(testParams.Location)
}
