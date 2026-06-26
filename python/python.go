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

// mainThreadState holds the interpreter state
var mainThreadState *C.PyThreadState

// Initialize starts the interpreter once.
func Initialize() error {
	C.Py_Initialize()
	if C.Py_IsInitialized() == 0 {
		return errors.New("failed to initialise cpython interpreter")
	}

	mainThreadState = C.PyEval_SaveThread()
	return nil
}

// Finalize reacquires the GIL and shuts the interpreter down. It must run on
// the same OS thread that called Initialize.
func Finalize() error {
	C.PyEval_RestoreThread(mainThreadState)
	mainThreadState = nil

	if C.Py_FinalizeEx() != 0 {
		return errors.New("finalise reported an error")
	}
	return nil
}

// WithGIL pins the current goroutine to its OS thread, acquires the GIL for
// the duration of f, and runs f inside a fresh object scope.
func WithGIL(f func()) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	gil := C.PyGILState_Ensure()
	defer C.PyGILState_Release(gil)

	// Track every object created during f and release them in LIFO order.
	WithScope(f)
}

// Run initialises the interpreter, runs f under the GIL, then finalises.
func Run(f func()) (err error) {
	// Keep Initialize and Finalize on the same OS thread.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err = Initialize(); err != nil {
		return err
	}

	// Shutdown once f returns, even on an early return inside it.
	defer func() {
		if ferr := Finalize(); ferr != nil && err == nil {
			err = ferr
		}
	}()

	WithGIL(f)
	return nil
}
