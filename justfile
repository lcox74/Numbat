compose := "docker compose run --rm runner"

# List the available recipes.
default:
    @just --list

# Build the runner image. Only needed for dockerfile changes.
build:
    docker compose build

# Run the demo program
run: build
    {{compose}} go run .

# Run the benchmark with allocation stats.
bench: build
    {{compose}} go test -bench=. -benchmem ./bench

