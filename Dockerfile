FROM golang:1.15-alpine as builder

WORKDIR /app

COPY . .

RUN go build -o traefik-hosts-generator

FROM debian:10-slim

COPY --from=builder /app/traefik-hosts-generator /traefik-hosts-generator

ENTRYPOINT [ "/traefik-hosts-generator" ]