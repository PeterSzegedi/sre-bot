FROM golang:alpine as builder

WORKDIR /app

COPY cmd cmd
COPY pkg pkg
COPY go.mod go.mod
COPY go.sum go.sum

RUN cd cmd;go get -d -v
RUN cd cmd;CGO_ENABLED=0 GOOS=linux go build -o sre-bot .

FROM alpine:latest

RUN apk update && apk add ca-certificates

COPY --from=builder /app/cmd/sre-bot sre-bot

ENV BOT_TOKEN ""
ENV APP_LEVEL_TOKEN ""
ENV PD_TOKEN ""
ENV L1_SCHEDULE ""
ENV L2_SCHEDULE ""

EXPOSE 8081

ENTRYPOINT [ "./sre-bot" ]
