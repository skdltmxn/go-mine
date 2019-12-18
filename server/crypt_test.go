package server

import (
	"bytes"
	"crypto/aes"
	"testing"
)

func TestAESCFB8(t *testing.T) {
	key := []byte("0123456789abcdef")
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Errorf("aes.NewCipher failed: %+v", err)
	}

	plain := []byte("Characteristics of a Golang test function")
	cipher := make([]byte, len(plain))

	// key == iv
	cfb := newCFB8Encrypter(block, key)

	// encrypt
	cfb.XORKeyStream(cipher, plain)

	// decrypt
	cfb.XORKeyStream(cipher, cipher)

	if bytes.Equal(plain, cipher) {
		t.Errorf("Plain != Cipher")
	}
}
