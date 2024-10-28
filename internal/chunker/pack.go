package chunker

type BlobType int

const (
	Data int = 1 << iota
	CompressedData
)

type Blob struct {
	Data []byte
	Type BlobType
}

type Header struct {
	Key  []byte
	Data []byte
	Len  int
}

type Pack struct {
	Data   []Blob
	Header Header
}
