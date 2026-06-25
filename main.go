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

func logf(format string, args ...any) {
	fmt.Printf("[go] "+format+"\n", args...)
}

// withCString hands a C string to fn and frees it afterwards.
func withCString(s string, fn func(*C.char)) {
	c := C.CString(s)
	defer C.free(unsafe.Pointer(c))
	fn(c)
}

// importModule imports a Python module by name.
func importModule(name string) *C.PyObject {
	var mod *C.PyObject
	withCString(name, func(c *C.char) {
		mod = C.PyImport_ImportModule(c)
	})

	// Print any error that might happen
	if mod == nil {
		C.PyErr_Print()
	}

	return mod
}

// getAttrString reads an attribute off a Python object.
func getAttrString(o *C.PyObject, attr string) *C.PyObject {
	var v *C.PyObject
	withCString(attr, func(c *C.char) {
		v = C.PyObject_GetAttrString(o, c)
	})

	return v
}

// goString decodes a Python str object into a Go string.
func goString(o *C.PyObject) string {
	return C.GoString(C.PyUnicode_AsUTF8(o))
}

func runPython(f func()) {
	// Keep interpreter calls on one OS thread.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Check if we can get the version from libpython
	logf("linked cpython %s", C.GoString(C.Py_GetVersion()))

	// Initialise the Interpreter
	C.Py_Initialize()
	if C.Py_IsInitialized() == 0 {
		logf("warning: failed to initialised cpython interpreter")
		return
	}
	logf("cpython interpreter initialised")

	// Shutdown once f returns, even on an early return inside it.
	defer func() {
		if C.Py_FinalizeEx() != 0 {
			logf("warning: finalise reported an error")
			return
		}
		logf("finalised cleanly")
	}()

	// Run the Python dependent code
	f()
}

func main() {
	runPython(func() {
		// Import NumPy
		np := importModule("numpy")
		if np == nil {
			logf("warning: could not import numpy")
			return
		}
		defer C.Py_DecRef(np)

		// Get NumPy Version
		ver := getAttrString(np, "__version__")
		if ver == nil {
			logf("warning: numpy has no __version__")
			return
		}
		defer C.Py_DecRef(ver)

		// Print the NumPy Version
		logf("numpy version %s", goString(ver))
	})
}
