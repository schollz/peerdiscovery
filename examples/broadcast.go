package main

import "github.com/schollz/peerdiscovery"

func main() {
	p := new(peerdiscovery.PeerDiscovery)
	p.Broadcast()
}
