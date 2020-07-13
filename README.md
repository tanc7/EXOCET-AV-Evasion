# Hello

For this endowed with the opportunity to receive the source code of this payload, this is the EXOCET module written in Golang.

# Primary purpose

It's main purpose is to take common malware such as

1. Shamoon wiper malware
2. Metasploit payloads
3. And other executables

And then use a combination of base64 encoding, obfuscation by generating arrays and array-maps to aid in it's reconstitution and using AES-256 encryption to take what is normally, a easily detectable payload, into one that is completely undetectable.

This is a subset of the DSX project. And is what gives the DSX it's infamous lethality. This is the "Exocet Missile" armament of the DSX Deep Standoff Attack Ship Server.

A lot of the malicious binaries are available at theZoo repo and will be converted into test payloads.

File dropping is avoided, as we do not want to touch the disk. Rather we decrypt, decode, and de-obfuscate the payload in our new executable and execute it in memory.

