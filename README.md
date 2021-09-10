# EXOCET - AV-evading, undetectable, payload delivery tool
<b>Chang Tan</b>
<br>
<b>AWS Certified Cloud Practitioner and Solutions Architect Associate</b>
<br>
changtan@listerunlimited.com
<br>

# Upcoming update! Direct encrypted shellcode execution!

I need a bit of help, because I successfully implemented CGO to execute encrypted shellcode but it is throwing memory access violations exit status 0xc0000005. It shouldn't be anything related to DEP (Data Execution Prevention) because the file CGOTest/working-template-shellcode-executor.go did run.

![](https://raw.githubusercontent.com/tanc7/EXOCET-AV-Evasion/master/shellcode-exec-works.png)

Once I figure this out, CGO was a pain in the ass to implement, we can now create crypters that execute INLINE-ASSEMBLY. Which was considered a impossibility until now.

# Come play around with the shellcode crypter and executor while you help me fix these memory access violation bugs

Note this requires Golang and the MinGW toolchain to be installed on Windows with you running and generating the shellcode on Windows. The reason why, is because CGO cannot be cross-compiled like our other EXOCET modules. To install the toolchain you need to go to [https://www.msys2.org/](https://www.msys2.org/) and follow the guide. Then you must add gcc to your environment variables in Windows

![](https://raw.githubusercontent.com/tanc7/EXOCET-AV-Evasion/master/sysenv.png)


**Step 1: Generate shellcode, this could be from msfvenom Meterpreter payloads, Cobalt Strike Beacons, or your own custom shellcode in C compatible format**
![](https://raw.githubusercontent.com/tanc7/EXOCET-AV-Evasion/master/generate-shellcode.png)

**Step 2: Copy only the bytes of the shellcode, excluding the quotes into a text file like sc.txt**
![](https://raw.githubusercontent.com/tanc7/EXOCET-AV-Evasion/master/copy-shellcode.png)

**Step 3: Your shellcode file should look like this. Raw shellcode**
![](https://raw.githubusercontent.com/tanc7/EXOCET-AV-Evasion/master/Shellcode-File.png)

**Step 4: Now run the command `go run exocet-shellcode-exec.go sc.txt shellcodetest.go KEY`**

**Step 5: You can attempt to run it but you'll run into memory access violation errors for some reason, which I am still working on**

# Note on Memory Access Violation Problem

Apparently, aside from the major limitations of CGO that prohibit or dramatically frustrates cross-compilation, the issue is that the shellcode we want to execute is landing in a section of memory (analyzed in WinDBG x64) that is not RWX. In other words, unless we write C code that explicitly allows execution in memory of the shellcode, it will always throw access violation errors.

The other method, that I observed other developers of rudimentary Go modules [https://gist.github.com/mgeeky/bb0fd5652b234fbd1c7630d7e5c8542d](https://gist.github.com/mgeeky/bb0fd5652b234fbd1c7630d7e5c8542d), is that they use Go's Windows API to interact with ntdll.dll and kernel32.dll to call VirtualAlloc and specify areas of RWX memory pages. This method works better, but it seems that the shellcode must be in num-transformed format only for it to work.

I am still working on this you guys. I may combine multiple programming languages together to write a proper shellcode execution module


![](https://upload.wikimedia.org/wikipedia/en/4/46/Exocet_impact.jpg)
<br>
![](https://raw.githubusercontent.com/tanc7/EXOCET-AV-Evasion/master/nodetections.png)

EXOCET is superior to Metasploit's "Evasive Payloads" modules as EXOCET uses AES-256 in GCM Mode (Galois/Counter Mode). Metasploit's Evasion Payloads uses a easy to detect RC4 encryption. While RC4 can decrypt faster, AES-256 is much more difficult to ascertain the intent of the malware.

However, it is possible to use Metasploit to build a Evasive Payload, and then chain that with EXOCET. So EXOCET will decrypt via AES-256, and then the Metasploit Evasive Payload then decrypts itself from RC4.

Much like my previous project, DarkLordObama, this toolkit is designed to be a delivery/launch vehicle, much like Veil-Evasion does. 

<a href="https://github.com/tanc7/dark-lord-obama">Dark Lord Obama Project</a>

However, EXOCET is not limited to a single codebase or platforms that are running Python. EXOCET works on ALL supported platforms and architectures that Go supports.

# Exocet Overview

EXOCET, is effectively a crypter-type malware dropper that can recycle easily detectable payloads like WannaCry, encrypt them using AES-GCM (Galois/Counter Mode), which is more secure than AES-CBC, and then create a dropper file for a majority of architectures and platforms out there. 

Basically...

1. It ingests dangerous malware that are now detectable by antivirus engines
2. It then encrypts them and produces it's own Go file
3. Then that Go file can be cross-compiled to 99% of known architectures
4. Upon execution, the encrypted payload is written to the disk and immediately executed on the command line

That means 32-bit, and 64-bit architectures, and it works on Linux, Windows, Macs, Unix, Android, iPhone, etc. You take, anything, and I mean ANYTHING, like the 1988 Morris Worm that nearly brought down the internet (which exploited a flaw in the fingerd listener daemon on UNIX), and make it a viable cyberweapon again.

EXOCET is designed to be used with the DSX Program, or the "Cyber Metal Gear" as I envisioned it. Being able to launch and proliferate dangerous malware without a traceable launch trail.

EXOCET is written entirely in Go.

# How to use

EXOCET, regardless of which binary you use to run it, requires Golang to work. By default, it generates a crypter .go file.

1. Windows users: <a href="https://golang.org/doc/install">Install Go Here</a>
2. Linux users: run `sudo apt-get update && sudo apt-get install -y golang`

## To run it

`go run EXOCET.go detectablemalware.exe outputmalware.go password123`, where password123 is your key. The longer and more random the key, the better it is able to evade detection from static binary analysis.

For 64-bit Windows Targets...

`env GOOS=windows GOARCH=amd64 go build outputmalware.go`

And out comes a `outputmalware.exe` file

For 64-bit MacOS Targets

`env GOOS=darwin GOARCH=amd64 go build outputmalware.go`

For 64-bit Linux Targets

`env GOOS=linux GOARCH=amd64 go build outputmalware.go`

See this reference on github for your parameters for other operating systems like Android <a href="https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63">Reference for Go Cross Compilation</a>

Furthermore, there are prebuilt binaries that I have made, meaning you just have to run `./EXOCET` or `EXOCET-Windows.exe`

If you are targeting Windows systems, it's highly recommended that you use `EXOCET-Windows-Proc-Inject` or `EXOCET-Windows-Process-injection.exe` instead. As it will immediately search and query any running commonly running Windows processes to inject Meterpreter payloads into.


# EXOCET live demo

<iframe width="560" height="315" src="https://github.com/tanc7/EXOCET-AV-Evasion/raw/master/exocetdemo.mp4" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

# Reason for the name

On May 4th, 1982, during the Falklands War, a squadron of Argentinan Super Eterdards launched a French made Exocet missile at the HMS Sheffield. Despite the Royal Navy's attempts to stop the missile, one struck, sinking the Sheffield. That incident literally put Argentina on the map as a show of force against a global colonial power.

<a href="https://www.theguardian.com/uk-news/2017/oct/15/exocet-missile-how-sinking-hms-sheffield-made-famous">News Article of the sinking of the HMS Sheffield</a>

Very much like how Onel de Guzman's actions with the ILOVEYOU virus put the Philippines on the map as a cyber threat. 

<a href="https://en.wikipedia.org/wiki/ILOVEYOU">ILOVEYOU Virus on Wikipedia</a>
