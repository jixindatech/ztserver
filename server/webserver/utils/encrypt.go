package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
)

var localKey string = ""
var localIv string = "yehuoshaobujin!!"

func SetEncryptKey(key string) error {
	if len(key) != 32 {
		return fmt.Errorf("%s", "invalid config secret length")
	}

	localKey = key
	return nil
}

// CBCEncryptWithPKCS7 CBC模式PKCS7填充AES加密
func CBCEncryptWithPKCS7(encodeStr string) (cryptedStr string, err error) {
	key := localKey
	iv := localIv
	if len(encodeStr) == 0 {
		return
	}
	// 根据key 生成密文
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return
	}

	// padding
	blockSize := block.BlockSize()

	// blockSize必须等于len(iv)
	if blockSize != len(iv) {
		err = errors.New("IV length must equal block size")
		return
	}
	encodeBytes := []byte(encodeStr)
	encodeBytes = pKCS7Padding(encodeBytes, blockSize)

	// 加密
	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	crypted := make([]byte, len(encodeBytes))
	blockMode.CryptBlocks(crypted, encodeBytes)

	// base64编码
	cryptedStr = base64.StdEncoding.EncodeToString(crypted)
	return
}

func pKCS7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	// 填充
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

// CBCDecryptWithPKCS7 CBC模式PKCS7填充AES解密
func CBCDecryptWithPKCS7(decodeStr string) (origDataStr string, err error) {
	key := localKey
	iv := localIv

	if len(decodeStr) == 0 {
		return
	}
	// 先解密base64
	decodeBytes, err := base64.StdEncoding.DecodeString(decodeStr)
	if err != nil {
		return
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return
	}

	// blockSize必须等于len(iv)
	if block.BlockSize() != len(iv) {
		err = errors.New("IV length must equal block size")
		return
	}
	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	origData := make([]byte, len(decodeBytes))
	blockMode.CryptBlocks(origData, decodeBytes)
	origData = pKCS7UnPadding(origData)

	origDataStr = string(origData)
	return
}

func pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
