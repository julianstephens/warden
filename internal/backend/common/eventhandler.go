package common

import "context"

type EventHandler interface {
	WriteConfig(ctx context.Context, reader IReader) error
	WriteKey(ctx context.Context, filename string, reader IReader) error
	WritePack(data []byte) error
}
