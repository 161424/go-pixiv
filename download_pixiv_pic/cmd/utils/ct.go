package utils

import (
	"crypto/aes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func AesEcrypt(key []byte, plaintext []byte) ([]byte, error) {
	// create cipher
	c, err := aes.NewCipher(key)

	// allocate space for ciphered data
	out := make([]byte, len(plaintext))

	// encrypt
	c.Encrypt(out, plaintext)

	return out, err
}

func AesDecrypt(key []byte, ct string) (string, error) {
	ciphertext, _ := hex.DecodeString(ct)

	c, err := aes.NewCipher(key)

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	s := string(pt[:])
	return s, err
}

func MD5(str string) string {
	data := []byte(str) //切片
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}
