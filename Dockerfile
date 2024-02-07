FROM golang:1.22-alpine AS build

WORKDIR /usr/src/icecon
COPY go.mod go.sum ./
RUN go mod download
COPY * .
RUN go generate -v ./...
RUN go build -v -ldflags "-s -w" .

###

FROM alpine:3.18

COPY --from=build /usr/src/icecon/icecon /usr/local/bin

STOPSIGNAL SIGTERM
ENTRYPOINT ["/usr/local/bin/icecon"]
