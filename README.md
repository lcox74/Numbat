# Numbat

This isn't a useful project, its just an experiment to see what I could do with
CPython in Golang. Turns out that I can drive Python libraries like NumPy
through this, GIL and all. There's a demo in `main.go` that builds two matrices
and multiplies them with NumPy, and a benchmark suite that pits smuggled NumPy
against native Go. I've been pretty much living out of the
[Python C API Reference Manual].

This is companion code for the blog post [Smuggling NumPy Into Go].

## Build & Run

`cgo` links against CPython through `pkg-config: python3-embed`. The repository
ships a container with every dependency already in place as I didn't want to
deal with additional pain. The recipes in the `justfile` wrap the compose
container:

```sh
just run     # build the image and run the matmul demo
just bench   # run the matmul benchmarks with allocation stats
```

If you'd rather skip `just`, the demo is plain `docker compose`:

```sh
docker compose up --build
```

I wouldn't run this locally.

## Benchmarks

The suite in `bench/` times a matrix multiply across a spread of sizes (64 up
to 2048) and compares four backends:

- NumPy driven through the embedded CPython interpreter
- a naive scalar Go loop
- gonum backed by OpenBLAS through the netlib cgo bindings

## License

Released under the [MIT License](LICENSE).

[Python C API Reference Manual]: https://docs.python.org/3.14/c-api/
[Smuggling NumPy Into Go]: https://lachlancox.dev/blog/2026-06-numpy-in-go/

