FROM golang:1.27-alpine

RUN apk add --no-cache \
    chromium \
    udev \
    ttf-freefont \
    ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

ENV PORT=8081
ENV CHROMIUM_PATH=/usr/bin/chromium-browser
EXPOSE 8081

CMD ["./main"]
