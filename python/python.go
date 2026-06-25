package python

/*
#cgo pkg-config: python3-embed
#include <Python.h>
*/
import "C"

import (
	"errors"
	"runtime"
)

func Version() string {
	return C.GoString(C.Py_GetVersion())
}

// Run initialises the interpreter, runs f, then finalises. All Python work
// must happen inside f, which is pinned to a single OS thread.
func Run(f func()) (err error) {
	// Keep interpreter calls on one OS thread.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	C.Py_Initialize()
	if C.Py_IsInitialized() == 0 {
		return errors.New("failed to initialise cpython interpreter")
	}

	// Shutdown once f returns, even on an early return inside it.
	defer func() {
		if C.Py_FinalizeEx() != 0 && err == nil {
			err = errors.New("finalise reported an error")
		}
	}()

	// Track every object created during f and release them in LIFO order.
	WithScope(f)
	return nil
}
