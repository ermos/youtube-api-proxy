FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY main.go ./main.go
COPY internal ./internal

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o http -v .

FROM alpine:latest

ARG PORT=8080

ENV PORT=$PORT

COPY --from=builder /build/http /http

EXPOSE $PORT

ENTRYPOINT ["/http"]

