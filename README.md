## GoShellCode

This github repository is supporting code for the webcast "Shellcode Execution with Golang".
There is a single Golang source file with three shellcode execution methods
contained within it.

* Direct Syscall execution
* Create a thread in the same process
* Inject a thread into a "srvhost.exe" process

1. To build as executable binary:

	go build

2. To build as DLL assuming 64-bit MINGW compiler is installed.

	go build -buildmode=c-shared -o gosc.dll gosc.go


Enjoy!

Joff Thyer
Copyright (c) 2021
River Gum Security LLC

