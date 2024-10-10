package backend

import (
	"context"
)

type Backend interface {
	GetName() string
	GetType() BackendType
	Sync(ctx context.Context, data []byte) error
	Handle(ctx context.Context, t FileType, data any) error
}

type Params interface{}

type backend struct {
	self    BackendType
	handler FileHandler
	name    string
}

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

func NewBackend(t BackendType, params Params) (Backend, error) {
	switch t {
	case LocalStorage:
		return newLocalStorage(params.(LocalStorageParams))
	default:
		return newLocalStorage(params.(LocalStorageParams))
	}
}
