FROM golang:1.22.6-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY .env .
EXPOSE 3000
EXPOSE 3001
ENTRYPOINT ["/app/main", "--env=production"]