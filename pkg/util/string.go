package util

import (
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"strconv"
)

func Base64Encode(o interface{}) (string, error) {
	v, err := json.Marshal(o)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(v), nil

}
func Base64Decode(v string, o interface{}) error {
	de, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return err
	}

	return json.Unmarshal(de, o)
}

func StrMustToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

const Letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GetLetter(index int) string {
	if index < 0 {
		return ""
	}
	if index < len(Letters) {
		return string(Letters[index-1])
	}
	if index == len(Letters) {
		return string(Letters[index-1])
	}

	l := index / 26

	if index%26 == 0 {
		l -= 1
	}
	if l <= 26 {
		return string(Letters[(l)-1]) + GetLetter(index-(26*l))
	}
	return string(Letters[(l)-1]) + GetLetter(index%26)
}

const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
