FROM golang

WORKDIR /peerdiscovery
COPY . .
RUN go build -v

CMD ["/peerdiscovery/ipv6"]
