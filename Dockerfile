FROM golang:1.20

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o server main.go

CMD ["./server"]