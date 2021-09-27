package config

import (
	"testing"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	tptu "github.com/libp2p/go-libp2p-transport-upgrader"
	"github.com/libp2p/go-tcp-transport"

	"github.com/stretchr/testify/require"
)

func TestTransportVariadicOptions(t *testing.T) {
	_, err := TransportConstructor(func(_ peer.ID, _ ...int) transport.Transport { return nil })
	require.NoError(t, err)
}

func TestConstructorWithOpts(t *testing.T) {
	var options []int
	c, err := TransportConstructor(func(_ *tptu.Upgrader, opts ...int) transport.Transport {
		options = opts
		return tcp.NewTCPTransport(nil)
	}, 42, 1337)
	require.NoError(t, err)
	_, err = c(nil, nil, nil)
	require.NoError(t, err)
	require.Equal(t, options, []int{42, 1337})
}
