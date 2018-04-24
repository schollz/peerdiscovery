# peerdiscovery

[![travis](https://travis-ci.org/schollz/peerdiscovery.svg?branch=master)](https://travis-ci.org/schollz/peerdiscovery) 
[![go report card](https://goreportcard.com/badge/github.com/schollz/peerdiscovery)](https://goreportcard.com/report/github.com/schollz/peerdiscovery) 
[![coverage](https://img.shields.io/badge/coverage-83%25-brightgreen.svg)](https://gocover.io/github.com/schollz/peerdiscovery)
[![godocs](https://godoc.org/github.com/schollz/peerdiscovery?status.svg)](https://godoc.org/github.com/schollz/peerdiscovery) 

Local peer discovery using UDP broadcast. I needed a peerdiscovery for [croc](https://github.com/schollz/croc) and everything I tried had problems, so I made another one.

![Example of peer discovery]()

## Install

Make sure you have Go 1.5+.

```
go get -u github.com/schollz/peerdiscovery
```

## Usage 

**Basic usage:**

```golang
p, _ := peerdiscovery.New()
discoveries, err := p.Discover()
if err != nil {
    panic(err)
}
for i, d := range discoveries {
    fmt.Printf("discovered '%s' with payload '%s'\n", i, d.Address, d.Payload)
}
```

See the docs for more usage.


## Contributing

Pull requests are welcome. Feel free to...

- Revise documentation
- Add new features
- Fix bugs
- Suggest improvements

## License

MIT