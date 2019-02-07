# peerdiscovery

[![travis](https://travis-ci.org/schollz/peerdiscovery.svg?branch=master)](https://travis-ci.org/schollz/peerdiscovery) 
[![go report card](https://goreportcard.com/badge/github.com/schollz/peerdiscovery)](https://goreportcard.com/report/github.com/schollz/peerdiscovery) 
[![coverage](https://cover.run/go/github.com/schollz/peerdiscovery.svg)](https://gocover.io/github.com/schollz/peerdiscovery)
[![godocs](https://godoc.org/github.com/schollz/peerdiscovery?status.svg)](https://godoc.org/github.com/schollz/peerdiscovery) 

Pure-go library for cross-platform thread-safe local peer discovery using UDP multicast. I needed to use peer discovery for [croc](https://github.com/schollz/croc) and everything I tried had problems, so I made another one.


## Install

Make sure you have Go 1.5+.

```
go get -u github.com/schollz/peerdiscovery
```

## Usage 

The following is a code to find the first peer on the local network and print it out.

```golang
discoveries, _ := peerdiscovery.Discover(peerdiscovery.Settings{Limit: 1})
for _, d := range discoveries {
    fmt.Printf("discovered '%s'\n", d.Address)
}
```

Here's the output when running on two computers. (*Run these gifs in sync by hitting Ctl + F5*).

**Computer 1:**

![computer 1](https://user-images.githubusercontent.com/6550035/39165714-ba7167d8-473a-11e8-82b5-fb7401ce2138.gif)

**Computer 2:**

![computer 1](https://user-images.githubusercontent.com/6550035/39165716-ba8db9ec-473a-11e8-96f7-e8c64faac676.gif)

For more examples, see [the scanning example](https://github.com/schollz/peerdiscovery/blob/master/examples/main.go) or [the docs](https://godoc.org/github.com/schollz/peerdiscovery).


## Testing

To test the peer discovery with just one host, one can launch multiple containers. The provided `Dockerfile` will run the example code.
Please make sure to enable [Docker's IPv6 support](https://docs.docker.com/v17.09/engine/userguide/networking/default_network/ipv6/) if you are using IPv6 for peer discovery.

```console
# Build the container, named peertest
$ docker build -t peertest .

# Execute the following command in multiple terminals
$ docker run -t --rm peertest
Scanning for 10 seconds to find LAN peers
 100% |████████████████████████████████████████|  [9s:0s]Found 1 other computers
0) '172.17.0.2' with payload 'zqrecHipCO'
```


## Contributing

Pull requests are welcome. Feel free to...

- Revise documentation
- Add new features
- Fix bugs
- Suggest improvements

## License

MIT
