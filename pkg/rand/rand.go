package rand

import (
	"crypto/rand"
	"encoding/base64"
)

func bytes(n int) ([]byte, error) {

	b := make([]byte, n)

	bytesRead, err := rand.Read(b)

	if err != nil || bytesRead < n {
		panic("Failed to generate random bytes")
	}

	return b, nil
}

func String(byteCount int) string {
	b, _ := bytes(byteCount)

	return base64.URLEncoding.EncodeToString(b)
}
