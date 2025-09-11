package platform

import (
	"errors"
	"syscall"
	"unsafe"
)

const (
	FO_DELETE          = 0x0003
	FOF_ALLOWUNDO      = 0x0040
	FOF_NOCONFIRMATION = 0x0010
)

type SHFILEOPSTRUCTW struct {
	Hwnd                  uintptr
	WFunc                 uint32
	PFrom                 *uint16
	PTo                   *uint16
	FFlags                uint16
	FAnyOperationsAborted int32
	HNameMappings         uintptr
	LpszProgressTitle     *uint16
}

var (
	shell32              = syscall.NewLazyDLL("shell32.dll")
	procSHFileOperationW = shell32.NewProc("SHFileOperationW")
)

func MoveToTrash(path string) error {
	// tupla-nollaterminoitu lista
	u16, err := syscall.UTF16PtrFromString(path + "\x00")
	if err != nil {
		return err
	}
	sh := SHFILEOPSTRUCTW{
		WFunc:  FO_DELETE,
		PFrom:  u16,
		FFlags: FOF_ALLOWUNDO | FOF_NOCONFIRMATION,
	}
	r1, _, callErr := procSHFileOperationW.Call(uintptr(unsafe.Pointer(&sh)))
	if r1 != 0 {
		if callErr != syscall.Errno(0) {
			return callErr
		}
		return errors.New("SHFileOperationW failed")
	}
	if sh.FAnyOperationsAborted != 0 {
		return errors.New("operation aborted")
	}
	return nil
}
