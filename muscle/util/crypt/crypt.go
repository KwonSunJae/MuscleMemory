package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// 암호화 함수
func Encrypt(plainText, key string) (string, error) {
	// AES 키는 16, 24, 32 바이트 길이여야 함
	aesKey := []byte(key)
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	// 패딩
	plainTextBytes := []byte(plainText)
	plainTextBytes = pkcs7Pad(plainTextBytes, aes.BlockSize)

	cipherText := make([]byte, aes.BlockSize+len(plainTextBytes))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainTextBytes)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// 복호화 함수
func Decrypt(cipherText, key string) (string, error) {
	aesKey := []byte(key)
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	if len(cipherTextBytes) < aes.BlockSize {
		return "", fmt.Errorf("암호문 길이가 너무 짧습니다.")
	}

	iv := cipherTextBytes[:aes.BlockSize]
	cipherTextBytes = cipherTextBytes[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherTextBytes, cipherTextBytes)

	// 언패딩
	plainTextBytes := pkcs7Unpad(cipherTextBytes)
	return string(plainTextBytes), nil
}

// PKCS7 패딩
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// PKCS7 언패딩
func pkcs7Unpad(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
