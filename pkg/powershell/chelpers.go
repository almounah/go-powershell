package powershell

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func makeUint64FromPtr(v uintptr) uint64 {
	return *((*uint64)(unsafe.Pointer(&v)))
}
func makeUintptrFromUint64(v uint64) uintptr {
	return *((*uintptr)(unsafe.Pointer(&v)))
}

func allocWrapper(size uint64) (uintptr, error) {
	return nativePowerShell_DefaultAlloc(size)
}

func Ensure(lastError syscall.Errno) error {
	if lastError != 0 {
		return lastError
	} else {
		return syscall.EINVAL
	}
}

func NotNill(ro uintptr, lastError syscall.Errno) error {
	if ro == 0 {
		return Ensure(lastError)
	}
	return nil

}

func LocalAlloc(length uint64) (ptr uintptr, err error) {
	var LocalAlloc_LPTR uint32 = 0x40
	modkernel32    := windows.NewLazySystemDLL("kernel32.dll")
	procLocalAlloc := modkernel32.NewProc("LocalAlloc")
	ptr, _, lastError := syscall.SyscallN(procLocalAlloc.Addr(), uintptr(LocalAlloc_LPTR), uintptr(length))
	err = NotNill(ptr, lastError)
	return
}
func localAllocWrapper(size uint64) (uintptr, error) {
	return LocalAlloc(size)
}

func freeWrapper(v uintptr) {
	nativePowerShell_DefaultFree(v)
}

func localMallocCopyLogStringHolder(input nativePowerShell_LogString_Holder) uintptr {

	size := uint64(unsafe.Sizeof(input))

	data, err := localAllocWrapper(size)
	if err != nil {
		panic("Couldn't allocate memory")
	}

	_ = memcpyLogStringHolder(data, input)

	return data
}

func mallocCopyArrayGenericPowerShellObject(input []nativePowerShell_GenericPowerShellObject) uintptr {

	inputCount := uint64(len(input))

	var size uint64 = 0
	if inputCount != 0 {
		size = inputCount * uint64(unsafe.Sizeof(input[0]))
	}

	data, err := allocWrapper(size)
	if err != nil {
		panic("Couldn't allocate memory")
	}

	if inputCount != 0 {
		_ = memcpyGenericPowerShellObject(data, &input[0], size)
	}

	return data
}

func mallocCopyStr(str string) uintptr {

	size := 2 * uint64((len(str) + 1))
	data, err := allocWrapper(size)
	if err != nil {
		panic("Couldn't allocate memory")
	}
	// safe usage due to data being c pointer
	_ = memcpyStr(data, str)

	return data
}

func wsclen(str uintptr) uint64 {
	var charCode uint16 = 1
	var i uint64 = 0
	for ; charCode != 0; i++ {
		charCode = castToUint16(str + (makeUintptrFromUint64(i) * unsafe.Sizeof(charCode)))
	}
	return i
}

func uintptrMakeString(ptr uintptr) string {
	return cstrToStr(ptr)
}

func makeCStringUintptr(str string) uintptr {
	return mallocCopyStr(str)
}
