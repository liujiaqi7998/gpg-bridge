//go:build windows

package singleinstance

import (
	"errors"
	"fmt"

	"golang.org/x/sys/windows"
)

const mutexName = "Local\\gpg-bridge-single-instance"

var ErrAlreadyActive = errors.New("already active")

type Guard struct {
	handle windows.Handle
}

func Acquire() (*Guard, error) {
	name, err := windows.UTF16PtrFromString(mutexName)
	if err != nil {
		return nil, fmt.Errorf("encode mutex name: %w", err)
	}
	mutex, err := windows.CreateMutex(nil, false, name)
	if err != nil {
		return nil, fmt.Errorf("create mutex: %w", err)
	}
	if windows.GetLastError() == windows.ERROR_ALREADY_EXISTS {
		windows.CloseHandle(mutex)
		return nil, ErrAlreadyActive
	}
	return &Guard{handle: mutex}, nil
}

func (g *Guard) Close() error {
	if g == nil || g.handle == 0 {
		return nil
	}
	return windows.CloseHandle(g.handle)
}
