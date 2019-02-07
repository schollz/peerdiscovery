FROM golang:1.11

WORKDIR /peerdiscovery
COPY . .
RUN go build ./examples/main.go

CMD ["/peerdiscovery/main"]
