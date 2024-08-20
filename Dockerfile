FROM golang

WORKDIR /peerdiscovery
COPY . .
RUN go get github.com/schollz/progressbar/v3
RUN go build ./examples/ipv4/main.go

CMD ["/peerdiscovery/main"]
