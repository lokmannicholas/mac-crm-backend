package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"dmglab.com/mac-crm/pkg/service"
)

func ASEEncrypt(stringToEncrypt string, keyString string) (encryptedString string) {

	keyString = hex.EncodeToString([]byte(MD5Hash(keyString)))
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		service.SysLog.Panicln(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		service.SysLog.Panicln(err.Error())
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		service.SysLog.Panicln(err.Error())
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

func ASEDecrypt(encryptedString string, keyString string) (decryptedString string) {
	keyString = hex.EncodeToString([]byte(MD5Hash(keyString)))
	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	block, err := aes.NewCipher(key)
	if err != nil {
		service.SysLog.Error(err.Error())
		return ""
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		service.SysLog.Error(err.Error())
		return ""
	}

	nonceSize := aesGCM.NonceSize()

	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		service.SysLog.Error(err.Error())
		return ""
	}
	if plaintext == nil {
		return ""
	}
	return string(plaintext[:])
}

func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func GenToken(id string) string {

	b := []byte(id)
	b2, err := time.Now().GobEncode()
	if err != nil {
		return ""
	}
	b = append(b, b2...)
	return base64.StdEncoding.EncodeToString(b)
}

func OpenSSLEncrypt() {

}
