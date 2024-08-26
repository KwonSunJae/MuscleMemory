package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"os"
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

func getMACAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range interfaces {
		if iface.HardwareAddr != nil {
			return iface.HardwareAddr.String()[:16], nil
		}
	}
	return "", fmt.Errorf("MAC 주소를 찾을 수 없습니다.")
}

func getUserName() (string, error) {
	userName := os.Getenv("USER") // Unix 계열 (Linux, macOS)
	if userName == "" {
		userName = os.Getenv("USERNAME") // Windows
	}

	if userName == "" {
		return "", fmt.Errorf("사용자 이름을 찾을 수 없습니다.")
	}

	return userName, nil
}

func GetOwner() (string, error) {
	userName, err := getUserName()
	if err != nil {
		return "", err
	}

	macAddr, err := getMACAddress()
	if err != nil {
		return "", err
	}

	return Encrypt(userName, macAddr)
}

func CompareOwner(owner1, owner2 string) bool {
	key, err := getMACAddress()
	if err != nil {
		return false
	}
	decryptedOwner1, err := Decrypt(owner1, key)
	if err != nil {
		return false
	}
	decryptedOwner2, err := Decrypt(owner2, key)
	if err != nil {
		return false
	}
	fmt.Println(decryptedOwner1, decryptedOwner2, decryptedOwner1 == decryptedOwner2)

	return decryptedOwner1 == decryptedOwner2
}
