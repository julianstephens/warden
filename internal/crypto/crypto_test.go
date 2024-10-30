package crypto_test

import (
	"errors"
	"testing"

	. "github.com/julianstephens/warden/internal/crypto"
)

// func assertEqual[T comparable](t *testing.T, expected T, actual T) {
// 	t.Helper()
// 	if expected == actual {
// 		return
// 	}
// 	t.Errorf("expected (%+v) is not equal to actual (%+v)", expected, actual)
// }

func AssertSliceEqual[T comparable](t *testing.T, expected []T, actual []T) {
	t.Helper()
	if len(expected) != len(actual) {
		t.Errorf("expected (%+v) is not equal to actual (%+v): len(expected)=%d len(actual)=%d",
			expected, actual, len(expected), len(actual))
	}
	for i := range expected {
		if expected[i] != actual[i] {
			t.Errorf("expected[%d] (%+v) is not equal to actual[%d] (%+v)",
				i, expected[i],
				i, actual[i])
		}
	}
}

func TestNewRandom(t *testing.T) {
	cases := []struct {
		name     string
		input    int
		expected int
		gotError error
	}{
		{
			name:     "should return byte slice of size 10",
			input:    10,
			expected: 10,
			gotError: nil,
		},
		{
			name:     "should error on input equals zero",
			input:    0,
			expected: 0,
			gotError: ErrInvalidRandomSize,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, err := NewRandom(c.input)

			if c.gotError != nil {
				if !errors.Is(err, c.gotError) {
					t.Fatalf("expected an error: %+v, but got: %+v", c.gotError, err)
					return
				}
			}

			if len(res) != c.expected {
				t.Fatalf("expected byte array of len: %d, but got: %d", c.expected, len(res))
			}
		})
	}
}

// func TestNewKey(t *testing.T) {
// 	r, err := NewRandom(24)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	securePassword := string(r)
// 	cases := []struct {
// 		name     string
// 		input    string
// 		expected int
// 		gotError error
// 	}{
// 		{
// 			name:     "should return new encryption key and aead cipher",
// 			input:    securePassword,
// 			expected: 32,
// 			gotError: nil,
// 		},
// 		{
// 			name:     "should return error when password is insecure",
// 			input:    "blah",
// 			expected: 0,
// 			gotError: ErrInvalidPassword,
// 		},
// 	}

// 	for _, c := range cases {
// 		t.Run(c.name, func(t *testing.T) {
// 			key, err := NewSessionKey(c.input)

// 			if c.gotError != nil {
// 				if !errors.Is(err, c.gotError) {
// 					t.Fatalf("expected an error: %+v, but got: %+v", c.gotError, err)
// 				}
// 				return
// 			}

// 			if key.Aead == nil {
// 				t.Fatal("expected crypto.AEAD, but got nil")
// 			}
// 			if len(key.EncryptionKey) != c.expected {
// 				t.Fatalf("expected encryption key of len: %d, but got: %d", c.expected, len(key.EncryptionKey))
// 			}
// 		})
// 	}
// }

// func TestEncryptDecrypt(t *testing.T) {
// 	r, err := NewRandom(24)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	key, err := NewSessionKey(string(r))
// 	if err != nil {
// 		t.Errorf("failed to create new encryption key: %+v", err)
// 	}

// 	text := "Hello, world!"
// 	enc, err := Encrypt(key, []byte(text), nil)
// 	if err != nil {
// 		t.Errorf("failed to encrypt text: %+v", err)
// 	}

// 	dec, err := Decrypt(key, enc)
// 	if err != nil {
// 		t.Errorf("failed to decrypt text: %+v", err)
// 	}

// 	fmt.Println(len(string(enc)))
// 	fmt.Println(len(string(dec)))

// 	assertSliceEqual(t, []byte(text), dec)
// }
