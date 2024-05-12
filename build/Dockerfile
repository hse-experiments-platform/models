FROM golang:1.22 AS builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/models ./cmd/models/main.go

FROM alpine:latest

WORKDIR /build

RUN git clone https://github.com/hse-experiments-platform/models-storage

COPY configs /build/configs
COPY swagger /build/swagger
COPY --from=builder /build/bin/models /build/bin/models

CMD ["./bin/models"]
