package main

import (
	"fmt"
	"strings"

	"github.com/lcox74/numbat/python"
)

func logf(format string, args ...any) {
	fmt.Printf("[go] "+format+"\n", args...)
}

// oneLine collapses numpy's multi-line array str onto a single line.
func oneLine(s string) string {
	rows := strings.Split(s, "\n")
	for i := range rows {
		rows[i] = strings.TrimSpace(rows[i])
	}

	return strings.Join(rows, " ")
}

func main() {
	logf("linked cpython %s", python.Version())

	err := python.Run(func() {
		np := python.Import("numpy")
		array := np.Attr("array")
		f64 := np.Attr("float64")

		// a = np.array([[1, 2], [3, 4]], dtype=np.float64)
		a := array.CallKw(
			[]any{[]any{[]any{1, 2}, []any{3, 4}}},
			map[string]any{"dtype": f64},
		)
		logf("a = %s", oneLine(a.String()))

		// b = np.array([[5, 6], [7, 8]], dtype=np.float64)
		b := array.CallKw(
			[]any{[]any{[]any{5, 6}, []any{7, 8}}},
			map[string]any{"dtype": f64},
		)
		logf("b = %s", oneLine(b.String()))

		// a @ b
		res := a.MatMul(b)
		logf("a @ b = %s", oneLine(res.String()))
	})
	if err != nil {
		logf("warning: %v", err)
		return
	}

	logf("finalised cleanly")
}
