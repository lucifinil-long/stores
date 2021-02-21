package utils

import (
	"crypto/rc4"
	"encoding/base64"
)

const (
	// DefaultEncryptKey is default encrypt key
	DefaultEncryptKey = "FlyingWing0"
)

// RC4Base64Encrypt encrpt the raw string using RC4 and formated result as base64 string
func RC4Base64Encrypt(raw string, key string) string {
	k := []byte(key)
	srctmp := []byte(raw)
	cl, _ := rc4.NewCipher(k)
	dst := make([]byte, len(srctmp))
	cl.XORKeyStream(dst, srctmp)
	str := base64.StdEncoding.EncodeToString(dst)

	return str
}

// RC4Base64Descrypt descrpt the base64 format encrpted string encrpted by RC4
func RC4Base64Descrypt(encrpted, key string) string {
	keyBytes := []byte(key)
	str, err := base64.StdEncoding.DecodeString(encrpted)
	if err != nil {
		return ""
	}
	data := []byte(str)
	ct, err := rc4.NewCipher(keyBytes)
	if err != nil {
		return ""
	}
	dst := make([]byte, len(data))
	ct.XORKeyStream(dst, data)
	return string(dst)
}
