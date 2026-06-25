# Numbat

This isn't a useful project, its just an experiment to see what I could do with
CPython in Golang. Turns out that I can drive Python libraries like NumPy
through this. There's a demo in `main.go` and I am working on doing a more
complicated demo and maybe some benchmarking against native Go code. I've been
pretty much living out of the [Python C API Reference Manual].

This is companion code for the blog post *Smuggling NumPy into Go* which I am
currently writing up.

## Build & Run

`cgo` links against CPython through `pkg-config: python3-embed`. The repository
ships a container with every dependency already in place as I didn't want to
deal with additional pain:

```sh
docker compose up --build
```

I wouldn't run this locally.

## License

Released under the [MIT License](LICENSE).

[Python C API Reference Manual]: https://docs.python.org/3.14/c-api/
