FROM golang:1.23.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go get ./cmd/api/
RUN go build -o main ./cmd/api/main.go

FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/main .
COPY .env .

EXPOSE 8080

CMD [ "/app/main" ]
