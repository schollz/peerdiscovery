FROM golang

WORKDIR /peerdiscovery
COPY . .
RUN go build ./examples/main.go

CMD ["/peerdiscovery/main"]
