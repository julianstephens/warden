package common

type FileType int

const (
	Config FileType = 1 << iota
	Key
	Pack
)

type Event struct {
	Type FileType
	Name *string
}
