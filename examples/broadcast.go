package main

import (
	"log"
	"time"

	"github.com/schollz/peerdiscovery"
)

func main() {
	p := peerdiscovery.New(peerdiscovery.Settings{
		Limit:   1,
		Payload: []byte(peerdiscovery.RandStringBytesMaskImprSrc(10)),
		Delay:   100 * time.Millisecond,
	})
	discoveries, err := p.Discover()
	if err != nil {
		log.Fatal(err)
	} else {
		for _, d := range discoveries {
			log.Printf("%s: %s", d.Address, d.Payload)
		}
	}
}
