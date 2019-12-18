package server

import (
	"crypto/cipher"
)

// same struct used in crypto/cipher
type cfb8 struct {
	b       cipher.Block
	next    []byte
	out     []byte
	outUsed int

	decrypt bool
}

func (x *cfb8) XORKeyStream(dst, src []byte) {
	if len(dst) < len(src) {
		panic("cfb8: output smaller than input")
	}

	for i := 0; i < len(src); i++ {
		x.b.Encrypt(x.out, x.next)
		val := src[i] ^ x.out[0]

		copy(x.next, x.next[1:])

		if x.decrypt {
			x.next[15] = src[i]
		} else {
			x.next[15] = val
		}

		dst[i] = val
	}
}

func newCFB8Encrypter(block cipher.Block, iv []byte) cipher.Stream {
	return newCFB8(block, iv, false)
}

func newCFB8Decrypter(block cipher.Block, iv []byte) cipher.Stream {
	return newCFB8(block, iv, true)
}

func newCFB8(block cipher.Block, iv []byte, decrypt bool) cipher.Stream {
	blockSize := block.BlockSize()
	if len(iv) != blockSize {
		// stack trace will indicate whether it was de or encryption
		panic("newCFB8: IV length must equal block size")
	}
	x := &cfb8{
		b:       block,
		out:     make([]byte, blockSize),
		next:    make([]byte, blockSize),
		outUsed: blockSize,
		decrypt: decrypt,
	}
	copy(x.next, iv)

	return x
}
