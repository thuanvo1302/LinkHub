FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o /bin/linkhub-api ./cmd/api

FROM alpine:3.22
WORKDIR /app
COPY --from=builder /bin/linkhub-api /usr/local/bin/linkhub-api
EXPOSE 8081
CMD ["linkhub-api"]
