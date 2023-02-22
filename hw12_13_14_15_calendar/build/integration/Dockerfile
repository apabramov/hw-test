FROM golang:1.17-alpine

WORKDIR /app
COPY . .

RUN go mod download
RUN go mod verify

ENV CGO_ENABLED=0

CMD go test -v -p 1 -tags integration ./test/integration