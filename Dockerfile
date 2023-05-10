FROM golang:1.20-alpine AS builder

WORKDIR /build

ENV GOPATH=/go/src

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build service.go



FROM alpine:3.18

COPY --from=builder /build/service /usr/local/bin/

RUN apk add --no-cache curl net-tools iputils-ping bind-tools

EXPOSE 26656 26657 1317 9090

USER 0

ENTRYPOINT ["sh", "-c"]
CMD ["service", "-port=8080"]
