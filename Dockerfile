ARG GO_VERSION=1.17

FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /app

COPY . .
RUN go get

COPY *.go .
RUN go build -o  ./stocky-boi-api

FROM alpine:latest

RUN apk update && apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/stocky-boi-api .

CMD [ "./stocky-boi-api" ]