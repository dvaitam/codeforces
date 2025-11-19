FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o codeforces-web webserver.go

FROM alpine:3.20

WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/codeforces-web /usr/local/bin/codeforces-web
COPY --from=builder /app /app

EXPOSE 8081

ENV DB_DSN=""

ENTRYPOINT ["/usr/local/bin/codeforces-web"]
