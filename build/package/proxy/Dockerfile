FROM golang:1.21-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED 0
ENV GOOS linux

COPY cmd cmd
COPY go.mod go.sum ./
COPY certs/ca.key certs/ca.crt certs/
COPY scripts/gen_cert.sh scripts/gen_key.sh scripts/
COPY internal internal

RUN go mod download && go mod tidy
RUN go build -o bin/proxy-server cmd/proxy/main.go

FROM alpine:3.19.1 as release

RUN apk update && apk add --no-cache openssl

COPY --from=builder /app/certs /certs
COPY --from=builder /app/scripts/ /scripts/
COPY --from=builder /app/bin/proxy-server /proxy-server

EXPOSE 8080

CMD [ "bin/sh", "-c", "/scripts/gen_key.sh && /proxy-server" ]
