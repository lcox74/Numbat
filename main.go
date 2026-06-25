package main

import (
	"fmt"

	"github.com/lcox74/numbat/python"
)

func logf(format string, args ...any) {
	fmt.Printf("[go] "+format+"\n", args...)
}

func main() {
	logf("linked cpython %s", python.Version())

	err := python.Run(func() {
		np := python.Import("numpy")
		defer np.DecRef()

		ones := np.Attr("ones")
		defer ones.DecRef()

		res := ones.Call(5)
		defer res.DecRef()

		logf("np.ones(5) = %s", res.String())
	})
	if err != nil {
		logf("warning: %v", err)
		return
	}

	logf("finalised cleanly")
}
