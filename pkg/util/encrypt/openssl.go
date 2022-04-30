package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"

	"dmglab.com/mac-crm/pkg/service"
	"golang.org/x/crypto/ssh"
)

func GetRSAPrivateKey(privateKeyFilePath string) (*rsa.PrivateKey, error) {

	privFile, err := ioutil.ReadFile(privateKeyFilePath)
	if err != nil {
		service.SysLog.Infoln("No RSA private key found, generating temp one", nil)
		return nil, err
	}

	privPem, _ := pem.Decode(privFile)
	if privPem.Type != "RSA PRIVATE KEY" {
		msg := fmt.Sprintf("RSA private key is of the wrong type %s", privPem.Type)
		return nil, errors.New(msg)
	}
	pk := BytesToRSAPrivateKey(privPem.Bytes)
	if pk == nil {
		return nil, errors.New("convert bytes to private key failed")
	}
	return pk, nil
}
func GetRSAPublicKey(publicKeyFilePath string) (*rsa.PublicKey, error) {

	pubFile, err := ioutil.ReadFile(publicKeyFilePath)
	if err != nil {
		service.SysLog.Infoln("No RSA private key found, generating temp one", nil)
		return nil, err
	}

	pub, rest := pem.Decode(pubFile)
	if len(rest) > 0 {
		return nil, errors.New("PEM block contains more than just public key")
	}

	if pub.Type != "RSA PUBLIC KEY" {
		msg := fmt.Sprintf("RSA public key is of the wrong type %s", pub.Type)
		return nil, errors.New(msg)
	}

	pk := BytesToRSAPublicKey(pub.Bytes)
	if pk == nil {
		return nil, errors.New("convert bytes to public key failed")
	}
	return pk, nil

}

func GetOPENSSLPublicKey(publicKeyFilePath string) (*rsa.PublicKey, error) {

	pub, err := ioutil.ReadFile(publicKeyFilePath)
	if err != nil {
		service.SysLog.Infoln("No RSA private key found, generating temp one", nil)
		return nil, err
	}
	parsed, _, _, _, err := ssh.ParseAuthorizedKey([]byte(pub))
	if err != nil {
		return nil, err
	}

	parsedCryptoKey := parsed.(ssh.CryptoPublicKey)

	// Then, we can call CryptoPublicKey() to get the actual crypto.PublicKey
	pubCrypto := parsedCryptoKey.CryptoPublicKey()

	// Finally, we can convert back to an *rsa.PublicKey
	if pubKey := pubCrypto.(*rsa.PublicKey); pubKey != nil {
		return pubKey, nil
	}
	return nil, errors.New("read public key failed")

}

func GetOPENSSLPrivateKey(privateKeyFilePath string) (*rsa.PrivateKey, error) {

	privFile, err := ioutil.ReadFile(privateKeyFilePath)
	if err != nil {
		service.SysLog.Infoln("No RSA private key found, generating temp one", nil)
		return nil, err
	}

	block, _ := pem.Decode([]byte(privFile))
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// PrivateKeyToBytes private key to bytes
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		service.SysLog.Panicln(err.Error())
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}

// BytesToOPENSSLPrivateKey bytes to private key
func BytesToRSAPrivateKey(block []byte) *rsa.PrivateKey {
	key, err := x509.ParsePKCS1PrivateKey(block)
	if err != nil {
		service.SysLog.Panicln(err.Error())
	}
	return key
}
func BytesToRSAPublicKey(block []byte) *rsa.PublicKey {
	rsaPubKey, err := x509.ParsePKIXPublicKey(block)
	if err != nil {
		service.SysLog.Panicln(err.Error())
	}
	return rsaPubKey.(*rsa.PublicKey)
}
func BytesToOPENSSLPrivateKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		service.SysLog.Infoln("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			service.SysLog.Panicln(err.Error())
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		service.SysLog.Panicln(err.Error())
	}
	return key
}

// BytesToOPENSSLPublicKey bytes to public key
func BytesToOPENSSLPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		service.SysLog.Infoln("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			service.SysLog.Panicln(err.Error())
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		service.SysLog.Panicln(err.Error())
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		service.SysLog.Panicln("not ok")
	}
	return key
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		service.SysLog.Panicln(err.Error())
	}
	return ciphertext
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		service.SysLog.Panicln(err.Error())
	}
	return plaintext
}

func GenRSAKeysToByte(bitSize int) ([]byte, []byte) {
	privkey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		service.SysLog.Panicln(err.Error())
	}
	pub := &privkey.PublicKey
	return PublicKeyToBytes(pub), PrivateKeyToBytes(privkey)
}

func GenOPENSSLKeys(bitSize int) ([]byte, []byte) {
	privkey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		service.SysLog.Panicln(err.Error())
	}

	pub, err := ssh.NewPublicKey(privkey.Public())
	if err != nil {
		return nil, nil
	}
	pubKeyStr := string(ssh.MarshalAuthorizedKey(pub))
	privKeyStr := marshalRSAPrivate(privkey)

	return []byte(pubKeyStr), []byte(privKeyStr)
}
func marshalRSAPrivate(priv *rsa.PrivateKey) string {
	return string(pem.EncodeToMemory(&pem.Block{
		Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv),
	}))
}
