package crypto

import (
	"crypto/sha256"

	"github.com/julianstephens/warden/internal/config"
)

func Hash(data []byte) config.ID {
	return sha256.Sum256(data)
}
