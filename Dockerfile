FROM golang:1.22.0-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/app

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

COPY .env .

COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs

CMD ./main
