FROM golang:1.24-alpine

COPY go.mod go.sum .

RUN go mod download

COPY . .

RUN go build

CMD ["./entrypoint.sh"]

