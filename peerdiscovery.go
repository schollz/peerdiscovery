package peerdiscovery

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// IPVersion specifies the version of the Internet Protocol to be used.
type IPVersion uint

const (
	IPv4 IPVersion = 4
	IPv6 IPVersion = 6
)

// Discovered is the structure of the discovered peers,
// which holds their local address (port removed) and
// a payload if there is one.
type Discovered struct {
	// Address is the local address of a discovered peer.
	Address string
	// Payload is the associated payload from discovered peer.
	Payload []byte
}

// Settings are the settings that can be specified for
// doing peer discovery.
type Settings struct {
	// Limit is the number of peers to discover, use < 1 for unlimited.
	Limit int
	// Port is the port to broadcast on (the peers must also broadcast using the same port).
	// The default port is 9999.
	Port string
	// MulticastAddress specifies the multicast address.
	// You should be able to use any of 224.0.0.0/4 or ff00::/8.
	// By default it uses the Simple Service Discovery Protocol
	// address (239.255.255.250 for IPv4 or ff02::c for IPv6).
	MulticastAddress string
	// Payload is the bytes that are sent out with each broadcast. Must be short.
	Payload []byte
	// Delay is the amount of time between broadcasts. The default delay is 1 second.
	Delay time.Duration
	// TimeLimit is the amount of time to spend discovering, if the limit is not reached.
	// The default time limit is 10 seconds.
	TimeLimit time.Duration
	// AllowSelf will allow discovery the local machine (default false)
	AllowSelf bool
	// DisableBroadcast will not allow sending out a broadcast
	DisableBroadcast bool
	// IPVersion specifies the version of the Internet Protocol (default IPv4)
	IPVersion IPVersion

	portNum                 int
	multicastAddressNumbers net.IP
}

// peerDiscovery is the object that can do the discovery for finding LAN peers.
type peerDiscovery struct {
	settings Settings

	received map[string][]byte
	sync.RWMutex
}

// initialize returns a new peerDiscovery object which can be used to discover peers.
// The settings are optional. If any setting is not supplied, then defaults are used.
// See the Settings for more information.
func initialize(settings Settings) (p *peerDiscovery, err error) {
	p = new(peerDiscovery)
	p.Lock()
	defer p.Unlock()

	// initialize settings
	p.settings = settings

	// defaults
	if p.settings.Port == "" {
		p.settings.Port = "9999"
	}
	if p.settings.IPVersion == 0 {
		p.settings.IPVersion = IPv4
	}
	if p.settings.MulticastAddress == "" {
		if p.settings.IPVersion == IPv4 {
			p.settings.MulticastAddress = "239.255.255.250"
		} else {
			p.settings.MulticastAddress = "ff02::c"
		}
	}
	if len(p.settings.Payload) == 0 {
		p.settings.Payload = []byte("hi")
	}
	if p.settings.Delay == 0 {
		p.settings.Delay = 1 * time.Second
	}
	if p.settings.TimeLimit == 0 {
		p.settings.TimeLimit = 10 * time.Second
	}
	p.received = make(map[string][]byte)
	p.settings.multicastAddressNumbers = net.ParseIP(p.settings.MulticastAddress)
	if p.settings.multicastAddressNumbers == nil {
		err = fmt.Errorf("Multicast Address %s could not be converted to an IP",
			p.settings.MulticastAddress)
		return
	}
	p.settings.portNum, err = strconv.Atoi(p.settings.Port)
	if err != nil {
		return
	}
	return
}

// Discover will use the created settings to scan for LAN peers. It will return
// an array of the discovered peers and their associate payloads. It will not
// return broadcasts sent to itself.
func Discover(settings ...Settings) (discoveries []Discovered, err error) {
	s := Settings{}
	if len(settings) > 0 {
		s = settings[0]
	}
	p, err := initialize(s)
	if err != nil {
		return
	}

	p.RLock()
	address := net.JoinHostPort(p.settings.MulticastAddress, p.settings.Port)
	portNum := p.settings.portNum
	payload := p.settings.Payload
	tickerDuration := p.settings.Delay
	timeLimit := p.settings.TimeLimit
	p.RUnlock()

	// get interfaces
	ifaces, err := net.Interfaces()
	if err != nil {
		return
	}

	// Open up a connection
	c, err := net.ListenPacket(fmt.Sprintf("udp%d", p.settings.IPVersion), address)
	if err != nil {
		return
	}
	defer c.Close()

	group := p.settings.multicastAddressNumbers

	// ipv{4,6} have an own PacketConn, which does not implement net.PacketConn
	var p2 interface{}
	if p.settings.IPVersion == IPv4 {
		p2 = ipv4.NewPacketConn(c)
	} else {
		p2 = ipv6.NewPacketConn(c)
	}

	for i := range ifaces {
		if p.settings.IPVersion == IPv4 {
			p2.(*ipv4.PacketConn).JoinGroup(&ifaces[i], &net.UDPAddr{IP: group, Port: portNum})
		} else {
			p2.(*ipv6.PacketConn).JoinGroup(&ifaces[i], &net.UDPAddr{IP: group, Port: portNum})
		}
	}

	go p.listen()
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()
	start := time.Now()
	for t := range ticker.C {
		exit := false
		p.RLock()
		if len(p.received) >= p.settings.Limit && p.settings.Limit > 0 {
			exit = true
		}
		p.RUnlock()

		if !s.DisableBroadcast {
			// write to multicast
			dst := &net.UDPAddr{IP: group, Port: portNum}
			for i := range ifaces {
				if p.settings.IPVersion == IPv4 {
					p24 := p2.(*ipv4.PacketConn)
					if errMulticast := p24.SetMulticastInterface(&ifaces[i]); errMulticast != nil {
						// log.Print(errMulticast)
						continue
					}
					p24.SetMulticastTTL(2)
					if _, errMulticast := p24.WriteTo([]byte(payload), nil, dst); errMulticast != nil {
						// log.Print(errMulticast)
						continue
					}
				} else {
					p26 := p2.(*ipv6.PacketConn)
					if errMulticast := p26.SetMulticastInterface(&ifaces[i]); errMulticast != nil {
						// log.Print(errMulticast)
						continue
					}
					p26.SetMulticastHopLimit(2)
					if _, errMulticast := p26.WriteTo([]byte(payload), nil, dst); errMulticast != nil {
						// log.Print(errMulticast)
						continue
					}
				}
			}
		}

		if exit || t.Sub(start) > timeLimit {
			break
		}
	}

	if !s.DisableBroadcast {
		// send out broadcast that is finished
		dst := &net.UDPAddr{IP: group, Port: portNum}
		for i := range ifaces {
			if p.settings.IPVersion == IPv4 {
				p24 := p2.(*ipv4.PacketConn)
				if errMulticast := p24.SetMulticastInterface(&ifaces[i]); errMulticast != nil {
					continue
				}
				p24.SetMulticastTTL(2)
				if _, errMulticast := p24.WriteTo([]byte(payload), nil, dst); errMulticast != nil {
					continue
				}
			} else {
				p26 := p2.(*ipv6.PacketConn)
				if errMulticast := p26.SetMulticastInterface(&ifaces[i]); errMulticast != nil {
					continue
				}
				p26.SetMulticastHopLimit(2)
				if _, errMulticast := p26.WriteTo([]byte(payload), nil, dst); errMulticast != nil {
					continue
				}
			}
		}
	}

	discoveries = make([]Discovered, len(p.received))
	i := 0
	p.RLock()
	for ip, payload := range p.received {
		discoveries[i] = Discovered{
			Address: ip,
			Payload: payload,
		}
		i++
	}
	p.RUnlock()
	return
}

const (
	// https://en.wikipedia.org/wiki/User_Datagram_Protocol#Packet_structure
	maxDatagramSize = 66507
)

// Listen binds to the UDP address and port given and writes packets received
// from that address to a buffer which is passed to a hander
func (p *peerDiscovery) listen() (recievedBytes []byte, err error) {
	p.RLock()
	address := net.JoinHostPort(p.settings.MulticastAddress, p.settings.Port)
	portNum := p.settings.portNum
	allowSelf := p.settings.AllowSelf
	p.RUnlock()
	localIPs := getLocalIPs()

	// get interfaces
	ifaces, err := net.Interfaces()
	if err != nil {
		return
	}
	// log.Println(ifaces)

	// Open up a connection
	c, err := net.ListenPacket(fmt.Sprintf("udp%d", p.settings.IPVersion), address)
	if err != nil {
		return
	}
	defer c.Close()

	group := p.settings.multicastAddressNumbers
	var p2 interface{}
	if p.settings.IPVersion == IPv4 {
		p2 = ipv4.NewPacketConn(c)
	} else {
		p2 = ipv6.NewPacketConn(c)
	}

	for i := range ifaces {
		if p.settings.IPVersion == IPv4 {
			p2.(*ipv4.PacketConn).JoinGroup(&ifaces[i], &net.UDPAddr{IP: group, Port: portNum})
		} else {
			p2.(*ipv6.PacketConn).JoinGroup(&ifaces[i], &net.UDPAddr{IP: group, Port: portNum})
		}
	}

	// Loop forever reading from the socket
	for {
		buffer := make([]byte, maxDatagramSize)
		var (
			n       int
			src     net.Addr
			errRead error
		)
		if p.settings.IPVersion == IPv4 {
			n, _, src, errRead = p2.(*ipv4.PacketConn).ReadFrom(buffer)
		} else {
			n, _, src, errRead = p2.(*ipv6.PacketConn).ReadFrom(buffer)
		}
		// log.Println(n, src.String(), err, buffer[:n])
		if errRead != nil {
			err = errRead
			return
		}

		srcHost, _, _ := net.SplitHostPort(src.String())

		if _, ok := localIPs[srcHost]; ok && !allowSelf {
			continue
		}

		// log.Println(src, hex.Dump(buffer[:n]))

		p.Lock()
		if _, ok := p.received[srcHost]; !ok {
			p.received[srcHost] = buffer[:n]
		}
		p.Unlock()
		p.RLock()
		if len(p.received) >= p.settings.Limit && p.settings.Limit > 0 {
			p.RUnlock()
			break
		}
		p.RUnlock()
	}

	return
}

// getLocalIPs returns the local ip address
func getLocalIPs() (ips map[string]struct{}) {
	ips = make(map[string]struct{})
	ips["localhost"] = struct{}{}
	ips["127.0.0.1"] = struct{}{}
	ips["::1"] = struct{}{}

	ifaces, err := net.Interfaces()
	if err != nil {
		return
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, address := range addrs {
			ip, _, err := net.ParseCIDR(address.String())
			if err != nil {
				// log.Printf("Failed to parse %s: %v", address.String(), err)
				continue
			}

			ips[ip.String()+"%"+iface.Name] = struct{}{}
			ips[ip.String()] = struct{}{}
		}
	}
	return
}
