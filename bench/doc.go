// Package bench holds matrix-multiplication benchmarks that compare matmul
// through smuggled NumPy against native Go implementations (a naive scalar
// loop, gonum's pure-Go BLAS, and gonum backed by OpenBLAS).
//
// Run them through the compose container, e.g. `just bench`.
package bench
