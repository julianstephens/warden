package store

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/julianstephens/warden/internal/backend/common"
	"github.com/julianstephens/warden/internal/crypto"
	"github.com/julianstephens/warden/internal/warden"
)

type Key struct {
	id warden.ID

	master *crypto.Key

	Username  string        `json:"username"`
	Hostname  string        `json:"hostname"`
	CreatedAt time.Time     `json:"createdAt"`
	Params    crypto.Params `json:"params"`
	Salt      []byte        `json:"salt"`
	Data      []byte        `json:"data"`
}

func (k *Key) ID() warden.ID {
	return k.id
}

func (k *Key) Valid() bool {
	return true
}

func (k *Key) String() string {
	template := func(data string) string { return fmt.Sprintf("<Key | %s>", data) }

	if k == nil {
		return template("nil")
	}

	return template(fmt.Sprintf("user: %s, host: %s, created: %s>", k.Username, k.Hostname, k.CreatedAt))
}

func (k *Key) Decrypt() *crypto.Key {
	return k.master
}

func deriveKey(ctx context.Context, params crypto.Params, password string) (key *Key, err error) {
	salt := crypto.NewSalt()

	derivedUser, err := crypto.NewIDKey(params, password, salt)
	if err != nil {
		return
	}

	master, err := crypto.NewSessionKey(salt)
	if err != nil {
		return
	}

	masterJson, err := json.Marshal(master)
	if err != nil {
		return
	}

	encMaster, err := crypto.Encrypt(*derivedUser, masterJson, nil)
	if err != nil {
		return
	}

	username, err := user.Current()
	if err != nil {
		err = fmt.Errorf("unable to get system user: %+v", err)
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		err = fmt.Errorf("unable to get system hostname: %+v", err)
		return
	}

	key = &Key{
		master:    master,
		Username:  username.Username,
		Hostname:  hostname,
		CreatedAt: time.Now(),
		Params:    params,
		Salt:      salt,
		Data:      encMaster,
	}

	return
}

func LoadKey(ctx context.Context, store *Store, params crypto.Params, password string) (*Key, error) {
	k, err := deriveKey(ctx, params, password)
	if err != nil {
		return nil, err
	}

	return k, nil
}

func AddKey(ctx context.Context, store *Store, params crypto.Params, password string) (*Key, error) {
	k, err := deriveKey(ctx, params, password)
	if err != nil {
		return nil, err
	}

	keyJson, err := json.Marshal(k)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal key to json: %+v", err)
	}

	id := crypto.Hash(keyJson)
	k.id = id

	event := common.Event{
		Type: common.Key,
		Name: id.String(),
	}

	err = store.backend.Put(ctx, event, common.NewByteReader(keyJson))
	if err != nil {
		return nil, err
	}

	return k, nil
}

func RemoveKey(id warden.ID) error { return nil }
