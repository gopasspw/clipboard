// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows
// +build windows

package clipboard

import (
	"context"
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	cfUnicodetext = 13
	gmemMoveable  = 0x0002
)

var (
	user32                     = windows.NewLazySystemDLL("user32.dll")
	isClipboardFormatAvailable = user32.NewProc("IsClipboardFormatAvailable")
	openClipboard              = user32.NewProc("OpenClipboard")
	closeClipboard             = user32.NewProc("CloseClipboard")
	emptyClipboard             = user32.NewProc("EmptyClipboard")
	getClipboardData           = user32.NewProc("GetClipboardData")
	setClipboardData           = user32.NewProc("SetClipboardData")

	kernel32     = windows.NewLazySystemDLL("kernel32")
	globalAlloc  = kernel32.NewProc("GlobalAlloc")
	globalFree   = kernel32.NewProc("GlobalFree")
	globalLock   = kernel32.NewProc("GlobalLock")
	globalUnlock = kernel32.NewProc("GlobalUnlock")
	lstrcpy      = kernel32.NewProc("lstrcpyW")
)

// waitOpenClipboard opens the clipboard, waiting for up to a second to do so.
func waitOpenClipboard() error {
	started := time.Now()
	limit := started.Add(time.Second)
	var r uintptr
	var err error
	for time.Now().Before(limit) {
		r, _, err = openClipboard.Call(0)
		if r != 0 {
			return nil
		}
		time.Sleep(time.Millisecond)
	}
	return err
}

func readAll(_ context.Context) ([]byte, error) {
	// LockOSThread ensure that the whole method will keep executing on the same thread from begin to end (it actually locks the goroutine thread attribution).
	// Otherwise if the goroutine switch thread during execution (which is a common practice), the OpenClipboard and CloseClipboard will happen on two different threads, and it will result in a clipboard deadlock.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if formatAvailable, _, err := isClipboardFormatAvailable.Call(cfUnicodetext); formatAvailable == 0 {
		return nil, err
	}
	err := waitOpenClipboard()
	if err != nil {
		return nil, err
	}

	h, _, err := getClipboardData.Call(cfUnicodetext)
	if h == 0 {
		_, _, _ = closeClipboard.Call()
		return nil, err
	}

	l, _, err := globalLock.Call(h)
	if l == 0 {
		_, _, _ = closeClipboard.Call()
		return nil, err
	}

	text := syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(l))[:]) //nolint:unsafeptr

	r, _, err := globalUnlock.Call(h)
	if r == 0 {
		_, _, _ = closeClipboard.Call()
		return nil, err
	}

	closed, _, err := closeClipboard.Call()
	if closed == 0 {
		return nil, err
	}
	return []byte(text), nil
}

func writeAll(_ context.Context, text []byte, _ bool) error {
	// LockOSThread ensure that the whole method will keep executing on the same thread from begin to end (it actually locks the goroutine thread attribution).
	// Otherwise if the goroutine switch thread during execution (which is a common practice), the OpenClipboard and CloseClipboard will happen on two different threads, and it will result in a clipboard deadlock.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	err := waitOpenClipboard()
	if err != nil {
		return err
	}

	r, _, err := emptyClipboard.Call(0)
	if r == 0 {
		_, _, _ = closeClipboard.Call()
		return err
	}

	data, err := syscall.UTF16FromString(string(text))
	if err != nil {
		_, _, _ = closeClipboard.Call()
		return err
	}

	// "If the hMem parameter identifies a memory object, the object must have
	// been allocated using the function with the GMEM_MOVEABLE flag."
	h, _, err := globalAlloc.Call(gmemMoveable, uintptr(len(data)*int(unsafe.Sizeof(data[0]))))
	if h == 0 {
		_, _, _ = closeClipboard.Call()
		return err
	}
	defer func() {
		if h != 0 {
			globalFree.Call(h)
		}
	}()

	l, _, err := globalLock.Call(h)
	if l == 0 {
		_, _, _ = closeClipboard.Call()
		return err
	}

	r, _, err = lstrcpy.Call(l, uintptr(unsafe.Pointer(&data[0])))
	if r == 0 {
		_, _, _ = closeClipboard.Call()
		return err
	}

	r, _, err = globalUnlock.Call(h)
	if r == 0 {
		if err.(syscall.Errno) != 0 {
			_, _, _ = closeClipboard.Call()
			return err
		}
	}

	r, _, err = setClipboardData.Call(cfUnicodetext, h)
	if r == 0 {
		_, _, _ = closeClipboard.Call()
		return err
	}
	h = 0 // suppress deferred cleanup
	closed, _, err := closeClipboard.Call()
	if closed == 0 {
		return err
	}
	return nil
}

func unsupported() bool {
	return false
}
