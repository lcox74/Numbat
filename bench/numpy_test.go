package bench

import (
	"fmt"
	"testing"

	"github.com/lcox74/numbat/python"
)

// BenchmarkMatMulNumpy measures NumPy's matmul driven through the embedded
// CPython interpreter with the NumPy extension.
func BenchmarkMatMulNumpy(b *testing.B) {
	for _, n := range sizes {
		b.Run(fmt.Sprintf("%dx%d", n, n), func(b *testing.B) {
			python.WithGIL(func() {
				np := python.Import("numpy")
				rand := np.Attr("random").Attr("rand")

				a := rand.Call(n, n)
				c := rand.Call(n, n)

				// Report FLOPs: a matmul of two NxN matrices is ~2*N^3 bytes.
				// Go labels the resulting column "MB/s", but it really is
				// FLOPS/s.
				b.SetBytes(int64(2 * n * n * n))
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					res := a.MatMul(c)
					res.DecRef()
				}
			})
		})
	}
}
