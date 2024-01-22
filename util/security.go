package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// 长度必须是16的倍数
var key = []byte("Haibara Ai!!!!!!")

func MD5(original string) string {
	hash := md5.New()
	hash.Write([]byte(original))
	return hex.EncodeToString(hash.Sum(nil))
}

func SHA256(original string) string {
	hash := sha256.Sum256([]byte(original))
	return hex.EncodeToString(hash[:])
}

func EncAES(original string) string {
	encrypted, err := EncryptAES([]byte(original), key)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(encrypted)
}

func DecAES(encrypted string) string {
	enc, _ := base64.StdEncoding.DecodeString(encrypted)
	original, err := DecryptAES(enc, key)
	if err != nil {
		panic(err)
	}
	return string(original)
}

func EncryptAES(original, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	original = pkcs7Padding(original, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encrypted := make([]byte, len(original))
	blockMode.CryptBlocks(encrypted, original)
	return encrypted, nil
}

func DecryptAES(encrypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	original := make([]byte, len(encrypted))
	blockMode.CryptBlocks(original, encrypted)
	original = pkcs7UnPadding(original)
	return original, nil
}

func pkcs7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func pkcs7UnPadding(original []byte) []byte {
	length := len(original)
	unPadding := int(original[length-1])
	return original[:(length - unPadding)]
}
