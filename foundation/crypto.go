package foundation

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"io/ioutil"
)

// Encrypt to text by AES-256, after that encode this by base64
func Encrypt(key []byte, data []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(data))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], data)
	encoded := base64.StdEncoding.EncodeToString(cipherText)
	return encoded, nil
}

// Decrypt by AES-256
func Decrypt(key []byte, encrypted string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := data[:aes.BlockSize]
	src := data[aes.BlockSize:]
	dst := make([]byte, len(src))
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(dst, src)

	return dst, nil
}

func GetKey(path string) []byte {
	p, err := ioutil.ReadFile(path)
	if err != nil {
		PrintError("Failed to read AES key file")
	}

	return GenKey(p)
}

func GenKey(src []byte) []byte {
	hash := sha256.Sum256(src)
	return hash[:]
}
