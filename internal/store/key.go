package store

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/julianstephens/warden/internal/crypto"
	"github.com/julianstephens/warden/internal/warden"
)

type Key struct {
	id warden.ID

	derivedUser *crypto.Key
	master      *crypto.Key

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

func AddKey(store *Store, password string) error {
	salt := crypto.NewSalt()

	derivedUser, err := crypto.NewIDKey(crypto.DefaultParams, password, salt)
	if err != nil {
		return err
	}

	master, err := crypto.NewSessionKey(salt)
	if err != nil {
		return err
	}

	masterJson, err := json.Marshal(master)
	if err != nil {
		return err
	}

	encMaster, err := crypto.Encrypt(*derivedUser, masterJson, nil)
	if err != nil {
		return err
	}

	id := crypto.Hash(masterJson)

	username, err := user.Current()
	if err != nil {
		err = fmt.Errorf("unable to get system user: %+v", err)
		return err
	}

	hostname, err := os.Hostname()
	if err != nil {
		err = fmt.Errorf("unable to get system hostname: %+v", err)
		return err
	}

	k := Key{
		derivedUser: derivedUser,
		master:      master,
		Username:    username.Username,
		Hostname:    hostname,
		CreatedAt:   time.Now(),
		Params:      crypto.DefaultParams,
		Salt:        salt,
		Data:        encMaster,
		id:          id,
	}

	return nil
}

func RemoveKey(store *Store, id warden.ID) error { return nil }
