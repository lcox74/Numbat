FROM golang:alpine

RUN apk add build-base python3-dev pkgconf py3-numpy

ENV CGO_ENABLED=1
WORKDIR /app

CMD ["go", "run", "main.go"]
