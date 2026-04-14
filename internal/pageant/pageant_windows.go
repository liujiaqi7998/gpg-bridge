//go:build windows

package pageant

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	puttyIPCMagic     = 0x804e50ba
	puttyIPCMaxLen    = 16384
	fileMapName       = "gpg_bridge"
	pageantWindow     = "Pageant"
	wmCopyData        = 0x004A
	fileMapAllAccess  = 0x000f001f
)

type copyDataStruct struct {
	dwData uintptr
	cbData uint32
	lpData uintptr
}

var (
	user32            = windows.NewLazySystemDLL("user32.dll")
	kernel32          = windows.NewLazySystemDLL("kernel32.dll")
	procFindWindowW   = user32.NewProc("FindWindowW")
	procSendMessageW  = user32.NewProc("SendMessageW")
	procMapViewOfFile = kernel32.NewProc("MapViewOfFile")
)

func Query(req []byte) ([]byte, error) {
	if len(req) > puttyIPCMaxLen {
		return nil, fmt.Errorf("request too large: %d", len(req))
	}
	name := fmt.Sprintf("%s-%d", fileMapName, windows.GetCurrentProcessId())
	namePtr, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return nil, fmt.Errorf("encode mapping name: %w", err)
	}

	hMap, err := windows.CreateFileMapping(windows.InvalidHandle, nil, windows.PAGE_READWRITE, 0, puttyIPCMaxLen, namePtr)
	if err != nil {
		return nil, fmt.Errorf("create file mapping: %w", err)
	}
	defer windows.CloseHandle(hMap)

	view, _, callErr := procMapViewOfFile.Call(uintptr(hMap), fileMapAllAccess, 0, 0, uintptr(puttyIPCMaxLen))
	if view == 0 {
		return nil, fmt.Errorf("map view of file: %v", callErr)
	}
	defer windows.UnmapViewOfFile(view)

	buf := unsafe.Slice((*byte)(unsafe.Pointer(view)), puttyIPCMaxLen)
	copy(buf, req)

	className, _ := windows.UTF16PtrFromString(pageantWindow)
	hwnd, _, _ := procFindWindowW.Call(uintptr(unsafe.Pointer(className)), uintptr(unsafe.Pointer(className)))
	if hwnd == 0 {
		return nil, fmt.Errorf("cannot find Pageant window")
	}

	mappingNameUTF16, err := windows.UTF16FromString(name)
	if err != nil {
		return nil, fmt.Errorf("build mapping payload: %w", err)
	}
	cds := copyDataStruct{
		dwData: puttyIPCMagic,
		cbData: uint32(len(mappingNameUTF16) * 2),
		lpData: uintptr(unsafe.Pointer(&mappingNameUTF16[0])),
	}
	res, _, sendErr := procSendMessageW.Call(hwnd, wmCopyData, 0, uintptr(unsafe.Pointer(&cds)))
	if res == 0 {
		return nil, fmt.Errorf("send WM_COPYDATA: %v", sendErr)
	}

	respLen := int(uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 | uint32(buf[3]))
	respLen += 4
	if respLen <= 4 || respLen > puttyIPCMaxLen {
		return nil, fmt.Errorf("invalid response length: %d", respLen)
	}
	resp := make([]byte, respLen)
	copy(resp, buf[:respLen])
	return resp, nil
}
