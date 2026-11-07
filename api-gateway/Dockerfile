FROM golang:1.27.3-alpine3.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:3.24

WORKDIR /app

COPY --from=builder /app/main .

ENV PORT=8080
EXPOSE 8080

CMD ["./main"]
