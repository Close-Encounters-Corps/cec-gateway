FROM golang:1.17
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG COMMITSHA
RUN go build -ldflags "-X main.COMMITSHA=$COMMITSHA" -o cec-gw



FROM ubuntu:focal

WORKDIR /app

RUN apt update && apt install ca-certificates -y

COPY --from=0 /build/cec-gw /app/cec-gw
CMD ["/app/cec-gw"]
