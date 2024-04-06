package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"
)

func fileIsExists(fn string) (bool, error) {
	if _, err := os.Stat(fn); err != nil {
		if !os.IsNotExist(err) {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func makeSessionToken() string {
	s := fmt.Sprintf("%d", time.Now().Unix())
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}
