FROM golang:1.21-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED 0
ENV GOOS linux

COPY cmd/webapi cmd/webapi
COPY go.mod go.sum ./
COPY internal internal

RUN go mod download && go mod tidy
RUN go build -o bin/webapi cmd/webapi/main.go

FROM scratch as release

COPY --from=builder /app/bin/webapi /webapi

EXPOSE 8000

CMD [ "/webapi" ]
