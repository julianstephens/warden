package common

import "context"

type BackendType int

const (
	LocalStorage BackendType = 1 << iota
	S3
	SFTP
)

//go:generate stringer -type=BackendType

var BackendTypeStringMap = func() map[string]BackendType {
	m := make(map[string]BackendType)
	for i := LocalStorage; i <= SFTP; i = i << 1 {
		m[i.String()] = i
	}
	return m
}()

var BackendTypes = func() []string {
	var m []string
	for i := LocalStorage; i <= SFTP; i = i << 1 {
		m = append(m, i.String())
	}
	return m
}()

type Backend interface {
	Save(ctx context.Context, event Event, reader IReader) error
}

type WardenBackend struct {
	Self    BackendType
	Handler EventHandler
	Name    string
}

var Resources = []string{"masterkey", "config"}
