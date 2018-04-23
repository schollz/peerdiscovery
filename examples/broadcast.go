package main

import "github.com/schollz/peerdiscovery"

func main() {
	p := peerdiscovery.New(peerdiscovery.Settings{
		Limit: 1,
	})
	p.Discover()
}
