package bench

import (
	"fmt"
	"math/rand"
	"testing"
)

// matmul is a naive way to do matrix multiplication, I don't really know any
// better ways but I am sure there are.
func matmul(a, b, c []float64, n int) {
	for i := range n {
		for k := range n {
			aik := a[i*n+k]
			row := c[i*n : i*n+n]
			brow := b[k*n : k*n+n]
			for j := range n {
				row[j] += aik * brow[j]
			}
		}
	}
}

// BenchmarkMatMulGoNaive is the native-Go baseline
func BenchmarkMatMulGoNaive(b *testing.B) {
	for _, n := range sizes {
		b.Run(fmt.Sprintf("%dx%d", n, n), func(b *testing.B) {
			r := rand.New(rand.NewSource(1))
			a := make([]float64, n*n)
			bb := make([]float64, n*n)
			c := make([]float64, n*n)
			for i := range a {
				a[i] = r.Float64()
				bb[i] = r.Float64()
			}

			// Report FLOPs: a matmul of two NxN matrices is ~2*N^3 bytes.
			// Go labels the resulting column "MB/s", but it really is
			// FLOPS/s.
			b.SetBytes(int64(2 * n * n * n))
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				for j := range c {
					c[j] = 0
				}

				matmul(a, bb, c, n)
			}
		})
	}
}
