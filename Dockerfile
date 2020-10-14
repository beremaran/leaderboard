FROM golang:1.15.2-buster AS builder

WORKDIR /go/src/app
COPY . .

WORKDIR /go/src/app/cmd/leaderboard
# TODO: make use of Makefile here
RUN go mod download && \
    go get -u github.com/swaggo/swag/cmd/swag && \
    swag init --parseInternal -g main.go && \
    CGO_ENABLED=0 go build -tags netgo -a -v

FROM alpine:latest

WORKDIR /app
COPY --from=builder /go/src/app/cmd/leaderboard/leaderboard leaderboard

EXPOSE 1323
CMD "/app/leaderboard"
