package main

/*
#cgo pkg-config: python3-embed
#include <Python.h>
*/
import "C"

import (
	"fmt"
	"runtime"
)

func logf(format string, args ...any) {
	fmt.Printf("[go] "+format+"\n", args...)
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
		np := Import("numpy")
		defer np.DecRef()

		onesFn := np.Attr("ones")
		defer onesFn.DecRef()

		t := C.PyTuple_New(1)
		C.PyTuple_SetItem(t, 0, AsPyObj(5))
		res := wrap(C.PyObject_CallObject(onesFn.p, t))
		defer res.DecRef()
		C.Py_DecRef(t)

		logf("np.ones(5) = %s", res.String())
	})
}
