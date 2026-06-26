package bench

import (
	"os"
	"runtime"
	"testing"

	"github.com/lcox74/numbat/python"
)

// sizes covers a spread from small to large so the benchmarks show how matmul
// scales.
var sizes = []int{64, 128, 256, 512, 1024, 2048}

// TestMain brings the CPython interpreter up once for the whole test system
// so it simplifies the test code.
func TestMain(m *testing.M) {
	runtime.LockOSThread()

	if err := python.Initialize(); err != nil {
		panic(err)
	}

	code := m.Run()

	if err := python.Finalize(); err != nil {
		panic(err)
	}
	runtime.UnlockOSThread()

	os.Exit(code)
}
