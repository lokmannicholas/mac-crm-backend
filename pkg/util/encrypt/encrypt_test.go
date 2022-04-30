package encrypt

import (
	"io/ioutil"
	"testing"
)

func TestASE(t *testing.T) {
	var str = "A1234566"
	var ASE_KEY = "19388973-e2a8-456a-9538-c56af8b80157" // 16, 24, 32 charaters

	result := ASEEncrypt(str, ASE_KEY)
	t.Log(result)
	ori := ASEDecrypt(result, ASE_KEY)
	t.Log(ori)
}

func TestMD5Hash(t *testing.T) {
	t.Log(MD5Hash("19388973-e2a8-456a-9538-c56af8b80157"))
}
func TestOpenSSLGenKeys(t *testing.T) {
	pubByte, privByte := GenOPENSSLKeys(2048)
	err := ioutil.WriteFile("test.pub", pubByte, 0644)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = ioutil.WriteFile("test", privByte, 0644)
	if err != nil {
		t.Fatal(err.Error())
	}

}
func TestOpenSSLRSA(t *testing.T) {
	msg := []byte{'A', 'B'}
	pubPath := "test.pub"
	privPath := "test"
	pub, err := GetOPENSSLPublicKey(pubPath)
	if err != nil {
		t.Fatal(err.Error())
	}
	encryptedMsg := EncryptWithPublicKey(msg, pub)
	t.Log(string(encryptedMsg))

	priv, err := GetOPENSSLPrivateKey(privPath)
	if err != nil {
		t.Fatal(err.Error())
	}
	decryptedMsg := DecryptWithPrivateKey(encryptedMsg, priv)
	t.Log(string(decryptedMsg))

}
