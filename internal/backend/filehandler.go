package backend

type FileType int

const (
	Config FileType = 1 << iota
	Key
	Pack
)

type FileHandler interface {
	putConfig(data []byte) error
	putKey(data []byte) error
	putPack(data []byte) error
}
