package python

// scope collects every Object created while it is active so they can be
// released together.
type scope struct {
	objs []*Object
}

var active *scope

// track registers o with the active scope, if any.
func track(o *Object) {
	if active != nil {
		active.objs = append(active.objs, o)
	}
}

// release decref's every tracked object in reverse (LIFO) order.
func (s *scope) release() {
	for i := len(s.objs) - 1; i >= 0; i-- {
		s.objs[i].DecRef()
	}
	s.objs = nil
}

// WithScope runs f inside a nested scope. Just helps with lazy cleaning for
// loops.
func WithScope(f func()) {
	s := &scope{}
	prev := active
	active = s
	defer func() {
		active = prev
		s.release()
	}()

	f()
}
