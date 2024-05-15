FROM golang:1.22 AS builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/models ./cmd/models/main.go

FROM alpine:latest

WORKDIR /build

COPY configs /build/configs
COPY swagger /build/swagger
COPY models_storage /build/models_storage
COPY --from=builder /build/bin/models /build/bin/models

CMD ["./bin/models"]
