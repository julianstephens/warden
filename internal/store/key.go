package store

import (
	"fmt"
	"time"

	"github.com/julianstephens/warden/internal/crypto"
	"github.com/julianstephens/warden/internal/warden"
)

type Key struct {
	id warden.ID

	self   *crypto.Key
	master *crypto.Key

	Username  string
	Hostname  string
	CreatedAt time.Time

	Salt []byte
	Data []byte
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

func AddKey(store *Store, id warden.ID, password string, master *crypto.Key) error {
	// k := &Key{}
	return nil
}

func RemoveKey(store *Store, id warden.ID) error { return nil }
