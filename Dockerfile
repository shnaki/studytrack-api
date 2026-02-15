FROM golang:1.25.4-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /api ./cmd/api

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /api .
COPY db/migrations ./db/migrations

EXPOSE 8080

CMD ["./api"]
