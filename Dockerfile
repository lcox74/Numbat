FROM golang:alpine

RUN apk add build-base python3-dev pkgconf py3-numpy openblas-dev

ENV CGO_ENABLED=1
ENV CGO_LDFLAGS="-lopenblas"
WORKDIR /app

CMD ["go", "run", "."]
