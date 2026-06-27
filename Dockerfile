FROM golang:1.25-alpine AS builder

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

# Run as a non-root user (Trivy DS-0002).
RUN adduser -D -u 10001 app
USER app

EXPOSE $PORT

ENTRYPOINT ["/http"]

