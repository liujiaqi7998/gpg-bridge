//go:build windows

package main

import "golang.org/x/sys/windows"

const (
	swHide = 0
)

var (
	kernel32Proc       = windows.NewLazySystemDLL("kernel32.dll")
	user32Proc         = windows.NewLazySystemDLL("user32.dll")
	getConsoleWindow   = kernel32Proc.NewProc("GetConsoleWindow")
	showWindowProc     = user32Proc.NewProc("ShowWindow")
)

func hideConsoleWindow() {
	hwnd, _, _ := getConsoleWindow.Call()
	if hwnd == 0 {
		return
	}
	showWindowProc.Call(hwnd, swHide)
}
