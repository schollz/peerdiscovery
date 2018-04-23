package peerdiscovery

import (
	"encoding/hex"
	"log"
	"net"
	"time"
)

var address = "239.0.0.0:9999"

func Broadcast() (err error) {
	conn, err := NewBroadcaster(address)
	if err != nil {
		return
	}
	go Listen()
	for {
		conn.Write([]byte("hello, world"))
		time.Sleep(1 * time.Second)
	}
	return
}

// NewBroadcaster creates a new UDP multicast connection on which to broadcast
func NewBroadcaster(address string) (*net.UDPConn, error) {
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

func Listen() (err error) {
	NewListener(address, msgHandler)
	return
}

const (
	maxDatagramSize = 8192
)

func msgHandler(src *net.UDPAddr, n int, b []byte) {
	log.Println(n, "bytes read from", src)
	log.Println(hex.Dump(b[:n]))
	log.Println(string(b))
}

// Listen binds to the UDP address and port given and writes packets received
// from that address to a buffer which is passed to a hander
func NewListener(address string, handler func(*net.UDPAddr, int, []byte)) {
	// Parse the string address
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
	}

	// Open up a connection
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	conn.SetReadBuffer(maxDatagramSize)

	// Loop forever reading from the socket
	for {
		buffer := make([]byte, maxDatagramSize)
		numBytes, src, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}

		handler(src, numBytes, buffer)
	}
}
