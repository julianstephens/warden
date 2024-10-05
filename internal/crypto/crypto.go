package crypto

import (
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"

	pkgerr "github.com/pkg/errors"

	"github.com/julianstephens/warden/internal/warden"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
)

const (
	passwordEntropy = 60
	keySize         = 32
	nonceSize       = 24
	saltSize        = 32
	totalBufPadding = nonceSize + keySize
)

var (
	ErrInvalidPassword   = errors.New("invalid password")
	ErrInvalidRandomSize = errors.New("cannot generate random array of zero length")
)

type Key struct {
	EncryptionKey []byte
	Aead          cipher.AEAD
}

func Hash(data []byte) warden.ID {
	return sha256.Sum256(data)
}

// NewKey generates a new random encryption key and MAC keys
func NewKey(password string) (Key, error) {
	err := passwordvalidator.Validate(password, passwordEntropy)
	if err != nil {
		return Key{}, pkgerr.Wrap(ErrInvalidPassword, err.Error())
	}

	salt, err := NewRandom(saltSize)
	if err != nil {
		return Key{}, fmt.Errorf("unable to create random salt: %+v", err)
	}

	key := argon2.IDKey([]byte(password), salt, uint32(5), uint32(64*1024), uint8(4), uint32(32))

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return Key{}, fmt.Errorf("unable to generate new MAC key: %+v", err)
	}

	return Key{EncryptionKey: key, Aead: aead}, nil
}

// NewRandom generates a cryptographically secure random byte array
func NewRandom(size int) ([]byte, error) {
	if size == 0 {
		return []byte{}, ErrInvalidRandomSize
	}
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}

// NewNonce generates a random nonce with capacity for the ciphertext
func NewNonce(size int, ciphertextLen int, overheadLen int) ([]byte, error) {
	res := make([]byte, size, size+ciphertextLen+overheadLen)
	_, err := rand.Read(res)
	if err != nil {
		return []byte{}, err
	}

	var total byte
	for _, x := range res {
		total |= x
	}

	if total > 0 {
		return res, nil
	}

	return []byte{}, fmt.Errorf("got invalid all-zero nonce")
}

// Encrypt secures data with XChacha20-Poly1305 algo
func Encrypt(key Key, plaintext []byte, additionalData *[]byte) ([]byte, error) {
	nonce, err := NewNonce(key.Aead.NonceSize(), len(plaintext), key.Aead.Overhead())
	if err != nil {
		return []byte{}, err
	}

	var res []byte
	if additionalData == nil {
		res = key.Aead.Seal(nonce, nonce, plaintext, nil)
	} else {
		res = key.Aead.Seal(nonce, nonce, plaintext, *additionalData)
	}

	return res, nil
}

func Decrypt(key Key, encrypted []byte) ([]byte, error) {
	if len(encrypted) < key.Aead.NonceSize() {
		return []byte{}, fmt.Errorf("ciphertext is too short")
	}

	nonce, ciphertext := encrypted[:key.Aead.NonceSize()], encrypted[key.Aead.NonceSize():]

	decrypted, err := key.Aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return []byte{}, err
	}

	return decrypted, nil
}
