package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
)

func Cipher(keyPhrase, value []byte) ([]byte, error) {
	gcm := gcmInstance(keyPhrase)
	nonce := make([]byte, gcm.NonceSize())
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, []byte(value), nil), nil
}
func gcmInstance(keyPhrase []byte) cipher.AEAD {
	aesBlock, err := aes.NewCipher([]byte(keyPhrase))
	if err != nil {
		log.Fatal(err)
	}
	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		log.Fatalln(err)
	}
	return gcm
}
func UnCipher(keyPhrase, ciphered []byte) ([]byte, error) {
	gcm := gcmInstance(keyPhrase)
	nonceSize := gcm.NonceSize()
	nonce, cipheredText := ciphered[:nonceSize], ciphered[nonceSize:]
	return gcm.Open(nil, nonce, cipheredText, nil)
}
