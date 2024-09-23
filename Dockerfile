FROM golang:1.23.1

WORKDIR /app

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o ./bin/indexer ./main.go

CMD ["./bin/indexer"]
