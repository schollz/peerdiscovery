package peerdiscovery

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDiscovery(t *testing.T) {
	for _, version := range []IPVersion{IPv4, IPv6} {
		// should not be able to "discover" itself
		discoveries, err := Discover()
		assert.Nil(t, err)
		assert.Zero(t, len(discoveries))

		// should be able to "discover" itself
		discoveries, err = Discover(Settings{
			Limit:     -1,
			AllowSelf: true,
			Payload:   []byte("payload"),
			Delay:     500 * time.Millisecond,
			TimeLimit: 1 * time.Second,
			IPVersion: version,
		})
		assert.Nil(t, err)
		assert.NotZero(t, len(discoveries))
	}
}

func TestDiscoverySelf(t *testing.T) {
	for _, version := range []IPVersion{IPv4, IPv6} {
		// broadcast self to self
		go func() {
			_, err := Discover(Settings{
				Limit:     -1,
				Payload:   []byte("payload"),
				Delay:     10 * time.Millisecond,
				TimeLimit: 1 * time.Second,
				IPVersion: version,
			})
			assert.Nil(t, err)
		}()
		discoveries, err := Discover(Settings{
			Limit:            1,
			Payload:          []byte("payload"),
			Delay:            500 * time.Millisecond,
			TimeLimit:        1 * time.Second,
			DisableBroadcast: true,
			AllowSelf:        true,
		})
		assert.Nil(t, err)
		assert.NotZero(t, len(discoveries))
	}
}
