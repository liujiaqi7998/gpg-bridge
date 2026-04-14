//go:build windows

package main

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	mbOK          = 0x00000000
	mbIconInfo    = 0x00000040
)

var (
	user32MessageBoxProc = user32Proc.NewProc("MessageBoxW")
)

func showInfoMessage(title string, message string) {
	titlePtr, err := windows.UTF16PtrFromString(title)
	if err != nil {
		return
	}
	messagePtr, err := windows.UTF16PtrFromString(message)
	if err != nil {
		return
	}
	user32MessageBoxProc.Call(0, uintptr(unsafe.Pointer(messagePtr)), uintptr(unsafe.Pointer(titlePtr)), mbOK|mbIconInfo)
}
