package tools

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"log"
)

func DecryptAES(content []byte, key string, iv string) (result string, err error) {
	keyBytes := []byte(key)
	ivBytes := []byte(iv)
	blockSize := 32
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return result, err
	}

	if len(content) < blockSize {
		return result, errors.New("cipherText too short")
	}

	if len(content)%aes.BlockSize != 0 {
		return result, errors.New("cipherText is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(content, content)

	content = PKCS5UnPadding(content)
	contentStringB64 := string(content)

	decContent, err2 := base64.StdEncoding.DecodeString(contentStringB64)
	if err2 != nil {
		return result, err2
	}

	log.Println(decContent)

	result = string(decContent)
	return result, nil
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	log.Println(length)
	log.Println(unpadding)

	return src[:(length - unpadding)]
}
