FROM golang:latest AS builder

RUN apt-get update

WORKDIR /furrybot
COPY . .
ENV CGO_ENABLED=0
RUN make build

FROM alpine:latest

COPY --from=builder /furrybot/build/furrybot /furrybot/bot

CMD ["/furrybot/bot"]