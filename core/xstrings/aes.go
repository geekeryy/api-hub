package xstrings

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"runtime"
)

var (
	ErrKeyLengthSixteen = errors.New("a sixteen or twenty-four or thirty-two length secret key is required")
	ErrPaddingSize      = errors.New("padding size error please check the secret key or iv")
	ErrIvAes            = errors.New("a sixteen-length ivaes is required")
)

const (
	Ivaes = "abcdefgh12345678"
)

// encrypt
func AesCbcEncryptBase64(plainText, secretKey string, ivAes []byte) (string, error) {
	if len(secretKey) != 16 && len(secretKey) != 24 && len(secretKey) != 32 {
		return "", ErrKeyLengthSixteen
	}
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}
	paddingText := PKCS5Padding([]byte(plainText), block.BlockSize())

	var iv []byte
	if len(ivAes) != 0 {
		if len(ivAes) != block.BlockSize() {
			return "", ErrIvAes
		} else {
			iv = ivAes
		}
	} else {
		iv = []byte(Ivaes)
	} // To initialize the vector, it needs to be the same length as block.blocksize
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(paddingText))
	blockMode.CryptBlocks(cipherText, paddingText)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// decrypt
func AesCbcDecryptBase64(cipherText, secretKey string, ivAes []byte) (string, error) {
	if len(secretKey) != 16 && len(secretKey) != 24 && len(secretKey) != 32 {
		return "", ErrKeyLengthSixteen
	}
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}
	decodedCipherText, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case runtime.Error:
				err = fmt.Errorf("runtime err=%v,Check that the key or text is correct", r)
			default:
				err = fmt.Errorf("error=%v,check the cipherText ", r)
			}
		}
	}()
	var iv []byte
	if len(ivAes) != 0 {
		if len(ivAes) != block.BlockSize() {
			return "", ErrIvAes
		} else {
			iv = ivAes
		}
	} else {
		iv = []byte(Ivaes)
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	paddingText := make([]byte, len(decodedCipherText))
	blockMode.CryptBlocks(paddingText, decodedCipherText)

	plainText, err := PKCS5UnPadding(paddingText, block.BlockSize())
	if err != nil {
		return "", err
	}
	return string(plainText), nil
}

func PKCS5Padding(plainText []byte, blockSize int) []byte {
	padding := blockSize - (len(plainText) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	newText := append(plainText, padText...)
	return newText
}

func PKCS5UnPadding(plainText []byte, blockSize int) ([]byte, error) {
	length := len(plainText)
	number := int(plainText[length-1])
	if number >= length || number > blockSize {
		return nil, ErrPaddingSize
	}
	return plainText[:length-number], nil
}
