FROM golang:1.24-alpine

COPY go.mod go.sum .

RUN go mod download

COPY . .

RUN go build

COPY token.json .

CMD ["./dw-archivist", "-target-id=${TARGET_PLAYLIST_ID}" ,"-discover-weekly-id=${SOURCE_PLAYLIST_ID}"]

