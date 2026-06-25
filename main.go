package main

/*
#cgo pkg-config: python3-embed
#include <Python.h>
*/
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

const source = "print('Hello, World!')"

func logf(format string, args ...any) {
	fmt.Printf("[go] "+format+"\n", args...)
}

func main() {
	// Keep interpreter calls on one OS thread.
	runtime.LockOSThread()

	// Check if we can get the version from libpython
	logf("linked cpython %s", C.GoString(C.Py_GetVersion()))

	// Initialise the Interpreter
	C.Py_Initialize()
	if C.Py_IsInitialized() == 0 {
		logf("warning: failed to initialised cpython interpreter")
		return
	}
	logf("cpython interpreter initialised")

	// Run some python code
	code := C.CString(source)
	defer C.free(unsafe.Pointer(code))
	if C.PyRun_SimpleString(code) != 0 {
		logf("warning: python code failed")
		return
	}

	// Shutdown
	if C.Py_FinalizeEx() != 0 {
		logf("warning: finalise reported an error")
		return
	}
	logf("finalised cleanly")
}
