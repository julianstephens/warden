package backend

import (
	"bytes"
	"io"
)

type IReader interface {
	Length() int64
	Reset() error
}

type ByteReader struct {
	*bytes.Reader
	IReader
	Len int64
}

func (b *ByteReader) Length() int64 {
	return b.Len
}

func (b *ByteReader) Reset() error {
	_, err := b.Reader.Seek(0, io.SeekStart)
	return err
}

func NewByteReader(data []byte) *ByteReader {
	return &ByteReader{
		Reader: bytes.NewReader(data),
		Len:    int64(len(data)),
	}
}
