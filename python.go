package main

/*
#cgo pkg-config: python3-embed
#include <Python.h>
*/
import "C"
import "unsafe"

type Object struct {
	p *C.PyObject
}

func wrap(p *C.PyObject) *Object {
	if p == nil {
		C.PyErr_Print()
		panic("python call returned nil")
	}

	return &Object{p}
}

func Import(name string) *Object {
	c := C.CString(name)
	defer C.free(unsafe.Pointer(c))

	return wrap(C.PyImport_ImportModule(c))
}

func (o *Object) Attr(name string) *Object {
	c := C.CString(name)
	defer C.free(unsafe.Pointer(c))

	return wrap(C.PyObject_GetAttrString(o.p, c))
}

func (o *Object) DecRef() {
	if o != nil && o.p != nil {
		C.Py_DecRef(o.p)
	}
}

func (o *Object) Float() float64 {
	return float64(C.PyFloat_AsDouble(o.p))
}

func (o *Object) String() string {
	// PyUnicode_AsUTF8 only works on str objects; stringify first so any
	// object (e.g. an ndarray) renders and we don't leave a dangling
	// exception that surfaces at interpreter shutdown.
	s := C.PyObject_Str(o.p)
	if s == nil {
		C.PyErr_Print()
		return ""
	}
	defer C.Py_DecRef(s)

	return C.GoString(C.PyUnicode_AsUTF8(s))
}

func AsPyObj(v any) *C.PyObject {
	switch x := v.(type) {
	case *Object:
		C.Py_IncRef(x.p)
		return x.p
	case float64:
		return C.PyFloat_FromDouble(C.double(x))
	case int:
		return C.PyLong_FromLongLong(C.longlong(x))
	case bool:
		b := C.long(0)
		if x {
			b = 1
		}

		return C.PyBool_FromLong(b)
	}

	return nil
}
