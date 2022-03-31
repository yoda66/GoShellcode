## GoShellCode

This github repository is supporting code for the webcast "Shellcode Execution with Golang".
There are three source file with shellcode execution methods.

* Direct Syscall execution
* Create a thread in the same process
* Inject a thread into a "srvhost.exe" process

To run cgo on a Windows platform you will need a MinGW compiler. 
I strongly advise TDM-GCC-64, as it is cross-compliant for both GOARCH=amd64 and GOARCH=386. 

http://tdm-gcc.tdragon.net/

Compiling a cgo program requires a few environment variables:

set GOOS=windows set GOARCH=amd64 set CGO_ENABLED=1

Unfortunately if you're working with mingw32-make, it doesn't seem to support .ONESHELL and if you have Git for Windows one-line solutions are also very fragile. The best solution I found to set these variables is batch files that wrap mingw32-make.


1. To build as executable binary:

	go build

2. To build as DLL assuming 64-bit MINGW compiler is installed.

	go build -buildmode=c-shared -o SysCall.dll Syscall.go



Enjoy!

Joff Thyer
Copyright (c) 2021
River Gum Security LLC

