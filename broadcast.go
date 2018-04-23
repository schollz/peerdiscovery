package peerdiscovery

import (
	"encoding/hex"
	"log"
	"net"
	"time"

	"github.com/pkg/errors"
)

var address = "239.0.0.0:9999"

type PeerDiscovery struct {
	listenerRecieved      []byte
	numClientsToDiscovery int
}

func (p *PeerDiscovery) Broadcast() {
	conn, err := newBroadcast(address)
	if err != nil {
		return
	}
	go p.Listen()
	for {
		if len(p.listenerRecieved) > 0 {
			break
		}
		conn.Write([]byte("hello, world"))
		time.Sleep(1 * time.Second)
	}
	return
}

// newBroadcast creates a new UDP multicast connection on which to broadcast
func newBroadcast(address string) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	return conn, nil

}

func (p *PeerDiscovery) Listen() (err error) {
	p.listenerRecieved, err = listen(address)
	return
}

const (
	maxDatagramSize = 8192
)

// Listen binds to the UDP address and port given and writes packets received
// from that address to a buffer which is passed to a hander
func listen(address string) (recievedBytes []byte, err error) {
	// Parse the string address
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return
	}

	// Open up a connection
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		return
	}
	defer conn.Close()

	conn.SetReadBuffer(maxDatagramSize)

	// Loop forever reading from the socket
	for {
		buffer := make([]byte, maxDatagramSize)
		numBytes, src, err2 := conn.ReadFromUDP(buffer)
		if err2 != nil {
			err = errors.Wrap(err2, "could not read from udp")
			return
		}

		if src.IP.String() == GetLocalIP() {
			continue
		}

		log.Println(numBytes, "bytes read from", src)
		log.Println(hex.Dump(buffer[:numBytes]))
		log.Println(string(buffer))
		recievedBytes = buffer[:numBytes]
		break
	}

	conn2, err := newBroadcast(address)
	if err != nil {
		return
	}
	defer conn2.Close()
	conn2.Write([]byte("ok"))
	return
}
