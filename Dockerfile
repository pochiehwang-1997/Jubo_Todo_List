FROM golang:alpine

COPY . /app
WORKDIR /app

RUN go mod download
RUN go build

CMD go run main.go