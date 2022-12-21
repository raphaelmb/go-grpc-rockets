FROM golang:1.19 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/server/main.go

FROM alpine AS production

COPY --from=builder /app .

CMD ["./app"]




