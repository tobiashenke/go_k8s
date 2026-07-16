FROM golang:1.26 as builder

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -a -o api main.go

FROM debian:bookworm-slim
WORKDIR /
COPY --from=builder /app/api .

ENTRYPOINT ["/api"]
