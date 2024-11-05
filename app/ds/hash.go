package ds

import (
	"bytes"
	"crypto/sha256"
	"io"
	"os"
)

var (
	left  = []byte{0x01, 0x10}
	right = []byte{0xFE, 0xEF}
)

func ExtractHash(msg []byte) []byte {
	if !bytes.Equal(msg[:len(left)], left) {
		return nil
	}

	if !bytes.Equal(msg[len(msg)-len(right):], right) {
		return nil
	}

	return msg[len(left) : len(msg)-len(right)]
}

func Hash(file *os.File) ([]byte, error) {
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
