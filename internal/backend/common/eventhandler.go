package common

import "context"

type EventHandler interface {
	WriteConfig(data []byte) error
	WriteKey(ctx context.Context, filename string, reader IReader) error
	WritePack(data []byte) error
}
