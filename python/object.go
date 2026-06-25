package python

/*
#cgo pkg-config: python3-embed
#include <Python.h>
*/
import "C"
import "unsafe"

// Object owns a reference to a CPython object. Call DecRef when done.
type Object struct {
	p *C.PyObject
}

// wrap takes ownership of a returned pointer and registers it with the active
// scope so it is decref'd when the scope ends.
func wrap(p *C.PyObject) *Object {
	if p == nil {
		C.PyErr_Print()
		panic("python call returned nil")
	}

	o := &Object{p}
	track(o)
	return o
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

func (o *Object) Call(args ...any) *Object {
	return o.CallKw(args, nil)
}

// CallKw calls o with positional args and keyword args. Pass nil kwargs for a
// plain call.
func (o *Object) CallKw(args []any, kwargs map[string]any) *Object {
	t := C.PyTuple_New(C.Py_ssize_t(len(args)))
	defer C.Py_DecRef(t)

	for i, a := range args {
		// PyTuple_SetItem steals the reference.
		C.PyTuple_SetItem(t, C.Py_ssize_t(i), toPyObject(a))
	}

	// PyObject_Call wants a dict or NULL; leave kw nil when there are none.
	var kw *C.PyObject
	if len(kwargs) > 0 {
		kw = C.PyDict_New()
		defer C.Py_DecRef(kw)

		for k, v := range kwargs {
			ck := C.CString(k)
			val := toPyObject(v)

			// PyDict_SetItemString does not steal, so drop our reference after.
			C.PyDict_SetItemString(kw, ck, val)
			C.Py_DecRef(val)
			C.free(unsafe.Pointer(ck))
		}
	}

	return wrap(C.PyObject_Call(o.p, t, kw))
}

// MatMul is the Python `@` operator: o @ other.
func (o *Object) MatMul(other *Object) *Object {
	return wrap(C.PyNumber_MatrixMultiply(o.p, other.p))
}

// DecRef releases the reference. It is idempotent, so an explicit DecRef and the
// automatic scope cleanup can both run without double-freeing.
func (o *Object) DecRef() {
	if o != nil && o.p != nil {
		C.Py_DecRef(o.p)
		o.p = nil
	}
}

func (o *Object) Float() float64 {
	return float64(C.PyFloat_AsDouble(o.p))
}

func (o *Object) String() string {
	// PyUnicode_AsUTF8 only works on str objects
	s := C.PyObject_Str(o.p)
	if s == nil {
		C.PyErr_Print()
		return ""
	}
	defer C.Py_DecRef(s)

	return C.GoString(C.PyUnicode_AsUTF8(s))
}

// toPyObject converts a Go value into a new Python object reference.
func toPyObject(v any) *C.PyObject {
	switch x := v.(type) {
	case *Object:
		C.Py_IncRef(x.p)
		return x.p
	case []any:
		list := C.PyList_New(C.Py_ssize_t(len(x)))
		for i, e := range x {
			C.PyList_SetItem(list, C.Py_ssize_t(i), toPyObject(e))
		}

		return list
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
