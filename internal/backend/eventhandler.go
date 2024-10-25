package backend

type EventHandler interface {
	putConfig(data []byte) error
	putKey(data []byte) error
	putPack(data []byte) error
}
