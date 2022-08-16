# Techniques to add from OSEP

1. AV Sandbox Detection by using Golang to call VirtualAllocExNuma
2. Sleep timers (very easy to implement)
3. Multi-stage payloads involving Windows Script Host (cscript.exe) shellcode runners
4. Offering a option for automatic code-signing if the user can provide valid code-signing certificates (this is already a feature in Cobalt Strike and many other C2 Frameworks)
5. Part of the multi-stage payloads require AMSI bypasses
6. PowerPick module from PowerSploit to allow execution of Powershell commands without Powershell
7. Polymorphic function (alphanumeric upper-lower) generation to evade signature detection technique (on scantime of the decrypter stubs), easily implemented since that's the point of Exocet. To generate a cross-compilable go file to serve as a delivery vehicle, any detection resulting after that is a issue with the payload itself
8. Converting all of the techniques from the OSEP course into Golang rather than C#.
9. Inline-hooking and patching of freshly spawned processes with shellcode trampolines
10. (Anti-debugging) ThreadLocalStorage Callbacks
11. (Anti-debugging) Registration of Top-Level Exception Handlers
12. Easy platform targeting by command line (all stages of the shellcode, including Windows Script Host Shellcode Runners, droppers, powershell modules, bypasses, and the final obfuscated RAT) will be auto-generateed in a finished folder just by specifying target platform (Windows, Mac, Linux, Unix) and architecture, as well as the automatic compiling, packing, and stripping of debugging symbols. Just specify the target platform (ex. -windows) and architecture (-amd64) and all instructions and stages will be premade for you.
13. Microsoft Doc and Excel Macro and VBA targeting. Allow specification of specific Microsoft Office Suite Versions before a VBA-stomp. Even though the OSEP performs a proof of concept with a hex editor, a generic word document can be reconstructed as a template payload with write primitives, much like how a malicious zip file can be reconstructed with Python scripts by reconstructing the headers that describe the document. **Actually, part of the docx, xlsx, docm, and xlsm file IS A ZIP FILE, therefore you have to reconstruct that as well.**

**I have seen personally, the deployment of exocet and analysis of payloads by threat actors for some time. From what I have personally seen, they have been used as a stager to deploy Beacons, Geacons (Golang Beacons for non Windows NT architectures) and used as a second-to-last stage payload**

# Packing conceals the key

The strings command against a PE32 executable can reveal the decryption password. The solution is to pack the payload via UPX.

Test it by

`strings testLinuxPayload.elf | grep 'key'`

And then test it again after running `upx-ucl`

`upx-ucl -o testLinuxPayload-packed.elf testLinuxPayload.elf`

`strings testLinuxPayload-packed.elf | grep 'key'`

Will need to add a UPX packer in Go

# Gobfuscate cannot be imported as a custom source code package

Attempts to import gobfuscate and call their individual modules and functions returns errors

# Metasploit breakage
Upgrading Golang or installing a go module appears to break Metasploit. It's really pissing me off.

Do not under any circumstances attempt to run install.sh yet. It will break metasploit if Go is upgraded because msf has some sort of go module called Zeitwerk.




```

        12: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/kernel.rb:30:in `require'
        11: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader/callbacks.rb:41:in `on_dir_autoloaded'
        10: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader/callbacks.rb:41:in `synchronize'
         9: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader/callbacks.rb:55:in `block in on_dir_autoloaded'
         8: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader/callbacks.rb:70:in `on_namespace_loaded'
         7: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader/callbacks.rb:70:in `each'
         6: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader/callbacks.rb:71:in `block in on_namespace_loaded'
         5: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:533:in `set_autoloads_in_dir'
         4: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:725:in `ls'
         3: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:725:in `foreach'
         2: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:733:in `block in ls'
         1: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:534:in `block in set_autoloads_in_dir'
/usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:558:in `rescue in block in set_autoloads_in_dir': wrong constant name Github.com.backup inferred by MsfAutoload::TempInflector from directory (Zeitwerk::NameError)

  /usr/share/metasploit-framework/lib/msf/core/modules/external/go/pkg/mod/cache/download/github.com.backup

Possible ways to address this:

  * Tell Zeitwerk to ignore this particular directory.
  * Tell Zeitwerk to ignore one of its parent directories.
  * Rename the directory to comply with the naming conventions.
  * Modify the inflector to handle this case.

┌──(root💀kali)-[/home/kali]
└─# mv /usr/share/metasploit-framework/lib/msf/core/modules/external/go/pkg/mod/cache/download/github.com.backup .

┌──(root💀kali)-[/home/kali]
└─# msfdb init && msfdb start && msfconsole
[i] Database already started
[i] The database appears to be already configured, skipping initialization
[i] Database already started
Traceback (most recent call last):
        46: from /usr/bin/msfconsole:18:in `<main>'
        45: from /usr/bin/msfconsole:18:in `require'
        44: from /usr/share/metasploit-framework/lib/msfenv.rb:17:in `<top (required)>'
        43: from /usr/share/metasploit-framework/lib/msfenv.rb:17:in `require'
        42: from /usr/share/metasploit-framework/config/environment.rb:4:in `<top (required)>'
        41: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/railties-6.1.4.1/lib/rails/railtie.rb:207:in `method_missing'
        40: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/railties-6.1.4.1/lib/rails/railtie.rb:207:in `public_send'
        39: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/railties-6.1.4.1/lib/rails/application.rb:391:in `initialize!'
        38: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/railties-6.1.4.1/lib/rails/initializable.rb:60:in `run_initializers'
        37: from /usr/lib/ruby/2.7.0/tsort.rb:205:in `tsort_each'
        36: from /usr/lib/ruby/2.7.0/tsort.rb:226:in `tsort_each'
        35: from /usr/lib/ruby/2.7.0/tsort.rb:347:in `each_strongly_connected_component'
        34: from /usr/lib/ruby/2.7.0/tsort.rb:347:in `call'
        33: from /usr/lib/ruby/2.7.0/tsort.rb:347:in `each'
        32: from /usr/lib/ruby/2.7.0/tsort.rb:349:in `block in each_strongly_connected_component'
        31: from /usr/lib/ruby/2.7.0/tsort.rb:431:in `each_strongly_connected_component_from'
        30: from /usr/lib/ruby/2.7.0/tsort.rb:350:in `block (2 levels) in each_strongly_connected_component'
        29: from /usr/lib/ruby/2.7.0/tsort.rb:228:in `block in tsort_each'
        28: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/railties-6.1.4.1/lib/rails/initializable.rb:61:in `block in run_initializers'
        27: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/railties-6.1.4.1/lib/rails/initializable.rb:32:in `run'
        26: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/railties-6.1.4.1/lib/rails/initializable.rb:32:in `instance_exec'
        25: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/railties-6.1.4.1/lib/rails/application/finisher.rb:133:in `block in <module:Finisher>'
        24: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:508:in `eager_load_all'
        23: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:508:in `each'
        22: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:393:in `eager_load'
        21: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:393:in `synchronize'
        20: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:404:in `block in eager_load'
        19: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:725:in `ls'
        18: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:725:in `foreach'
        17: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:733:in `block in ls'
        16: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:416:in `block (2 levels) in eager_load'
        15: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:416:in `const_get'
        14: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/kernel.rb:30:in `require'
        13: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader/callbacks.rb:41:in `on_dir_autoloaded'
        12: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader/callbacks.rb:41:in `synchronize'
        11: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader/callbacks.rb:55:in `block in on_dir_autoloaded'
        10: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader/callbacks.rb:70:in `on_namespace_loaded'
         9: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader/callbacks.rb:70:in `each'
         8: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader/callbacks.rb:71:in `block in on_namespace_loaded'
         7: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:533:in `set_autoloads_in_dir'
         6: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:725:in `ls'
         5: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:725:in `foreach'
         4: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:733:in `block in ls'
         3: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:551:in `block in set_autoloads_in_dir'
         2: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:583:in `autoload_subdir'
         1: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:777:in `cdef?'
/usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:777:in `const_defined?': wrong constant name Golang.org.backup (NameError)
        44: from /usr/bin/msfconsole:18:in `<main>'
        43: from /usr/bin/msfconsole:18:in `require'
        42: from /usr/share/metasploit-framework/lib/msfenv.rb:17:in `<top (required)>'
        41: from /usr/share/metasploit-framework/lib/msfenv.rb:17:in `require'
        40: from /usr/share/metasploit-framework/config/environment.rb:4:in `<top (required)>'
        39: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/railties-6.1.4.1/lib/rails/railtie.rb:207:in `method_missing'
        38: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/railties-6.1.4.1/lib/rails/railtie.rb:207:in `public_send'
        37: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/railties-6.1.4.1/lib/rails/application.rb:391:in `initialize!'
        36: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/railties-6.1.4.1/lib/rails/initializable.rb:60:in `run_initializers'
        35: from /usr/lib/ruby/2.7.0/tsort.rb:205:in `tsort_each'
        34: from /usr/lib/ruby/2.7.0/tsort.rb:226:in `tsort_each'
        33: from /usr/lib/ruby/2.7.0/tsort.rb:347:in `each_strongly_connected_component'
        32: from /usr/lib/ruby/2.7.0/tsort.rb:347:in `call'
        31: from /usr/lib/ruby/2.7.0/tsort.rb:347:in `each'
        30: from /usr/lib/ruby/2.7.0/tsort.rb:349:in `block in each_strongly_connected_component'
        29: from /usr/lib/ruby/2.7.0/tsort.rb:431:in `each_strongly_connected_component_from'
        28: from /usr/lib/ruby/2.7.0/tsort.rb:350:in `block (2 levels) in each_strongly_connected_component'
        27: from /usr/lib/ruby/2.7.0/tsort.rb:228:in `block in tsort_each'
        26: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/railties-6.1.4.1/lib/rails/initializable.rb:61:in `block in run_initializers'
        25: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/railties-6.1.4.1/lib/rails/initializable.rb:32:in `run'
        24: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/railties-6.1.4.1/lib/rails/initializable.rb:32:in `instance_exec'
        23: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/railties-6.1.4.1/lib/rails/application/finisher.rb:133:in `block in <module:Finisher>'
        22: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:508:in `eager_load_all'
        21: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:508:in `each'
        20: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:393:in `eager_load'
        19: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:393:in `synchronize'
        18: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:404:in `block in eager_load'
        17: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:725:in `ls'
        16: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:725:in `foreach'
        15: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:733:in `block in ls'
        14: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:416:in `block (2 levels) in eager_load'
        13: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:416:in `const_get'
        12: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/kernel.rb:30:in `require'
        11: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader/callbacks.rb:41:in `on_dir_autoloaded'
        10: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader/callbacks.rb:41:in `synchronize'
         9: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader/callbacks.rb:55:in `block in on_dir_autoloaded'
         8: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader/callbacks.rb:70:in `on_namespace_loaded'
         7: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader/callbacks.rb:70:in `each'
         6: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader/callbacks.rb:71:in `block in on_namespace_loaded'
         5: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:533:in `set_autoloads_in_dir'
         4: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:725:in `ls'
         3: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:725:in `foreach'
         2: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:733:in `block in ls'
         1: from /usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:534:in `block in set_autoloads_in_dir'
/usr/share/metasploit-framework/vendor/bundle/ruby/2.7.0/gems/zeitwerk-2.4.2/lib/zeitwerk/loader.rb:558:in `rescue in block in set_autoloads_in_dir': wrong constant name Golang.org.backup inferred by MsfAutoload::TempInflector from directory (Zeitwerk::NameError)

  /usr/share/metasploit-framework/lib/msf/core/modules/external/go/pkg/mod/cache/download/golang.org.backup

Possible ways to address this:

  * Tell Zeitwerk to ignore this particular directory.
  * Tell Zeitwerk to ignore one of its parent directories.
  * Rename the directory to comply with the naming conventions.
  * Modify the inflector to handle this case.

```
