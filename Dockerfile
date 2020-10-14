FROM golang:1.15.2-buster AS builder

WORKDIR /go/src/app
COPY . .

WORKDIR /go/src/app/cmd/leaderboard
RUN go mod download && \
    swag init --parseInternal -g cmd/leaderboard/main.go && \
    CGO_ENABLED=0 go build -tags netgo -a -v

FROM alpine:latest

WORKDIR /app
COPY --from=builder /go/src/app/cmd/leaderboard/leaderboard leaderboard

EXPOSE 1323
CMD "/app/leaderboard"
