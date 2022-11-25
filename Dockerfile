FROM golang:alpine as builder

WORKDIR /app

COPY main.go main.go
COPY go.mod go.mod
COPY go.sum go.sum

RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -o sre-bot .

FROM alpine:latest

RUN apk update && apk add ca-certificates

COPY --from=builder /app/sre-bot sre-bot

ENV BOT_TOKEN ""
ENV APP_LEVEL_TOKEN ""

EXPOSE 8081

ENTRYPOINT [ "./sre-bot" ]
