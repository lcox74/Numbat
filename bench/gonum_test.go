package bench

import (
	"fmt"
	"math/rand"
	"testing"

	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/netlib/blas/netlib"
)

// init routes gonum's BLAS calls through OpenBLAS via the netlib cgo bindings,
// instead of gonum's default pure-Go implementation which is alot slower.
func init() {
	blas64.Use(netlib.Implementation{})
}

// BenchmarkMatMulGonumOpenBLAS measures gonum's matmul backed by OpenBLAS.
func BenchmarkMatMulGonumOpenBLAS(b *testing.B) {
	for _, n := range sizes {
		b.Run(fmt.Sprintf("%dx%d", n, n), func(b *testing.B) {
			r := rand.New(rand.NewSource(1))
			a := make([]float64, n*n)
			bb := make([]float64, n*n)
			for i := range a {
				a[i] = r.Float64()
				bb[i] = r.Float64()
			}

			am := mat.NewDense(n, n, a)
			bm := mat.NewDense(n, n, bb)
			var cm mat.Dense

			// Report FLOPs: a matmul of two NxN matrices is ~2*N^3 bytes.
			// Go labels the resulting column "MB/s", but it really is
			// FLOPS/s.
			b.SetBytes(int64(2 * n * n * n))
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				cm.Mul(am, bm)
			}
		})
	}
}
