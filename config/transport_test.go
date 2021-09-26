package config

import (
	"testing"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
)

func TestTransportVariadicOptions(t *testing.T) {
	_, err := TransportConstructor(func(_ peer.ID, _ ...int) transport.Transport { return nil })
	if err != nil {
		t.Fatal(err)
	}
}
