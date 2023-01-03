package identify

import (
	"context"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	blhost "github.com/libp2p/go-libp2p/p2p/host/blank"
	swarmt "github.com/libp2p/go-libp2p/p2p/net/swarm/testing"

	"github.com/stretchr/testify/require"
)

func TestHandlerClose(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h1 := blhost.NewBlankHost(swarmt.GenSwarm(t))
	defer h1.Close()
	ids1, err := NewIDService(h1)
	require.NoError(t, err)
	ph := newPeerHandler(h1.ID(), ids1)
	closedCh := make(chan struct{}, 2)
	ph.start(ctx, func() {
		closedCh <- struct{}{}
	})

	require.NoError(t, ph.stop())
	select {
	case <-closedCh:
	case <-time.After(time.Second):
		t.Fatal("expected the handler to close")
	}

	require.NoError(t, ph.stop())
	select {
	case <-closedCh:
		t.Fatal("expected only one close event")
	case <-time.After(10 * time.Millisecond):
	}
}

func TestPeerSupportsProto(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h1 := blhost.NewBlankHost(swarmt.GenSwarm(t))
	defer h1.Close()
	ids1, err := NewIDService(h1)
	require.NoError(t, err)

	rp := peer.ID("test")
	ph := newPeerHandler(rp, ids1)
	require.NoError(t, h1.Peerstore().AddProtocols(rp, "test"))
	require.True(t, ph.peerSupportsProtos(ctx, []string{"test"}))
	require.False(t, ph.peerSupportsProtos(ctx, []string{"random"}))

	// remove support for protocol and check
	require.NoError(t, h1.Peerstore().RemoveProtocols(rp, "test"))
	require.False(t, ph.peerSupportsProtos(ctx, []string{"test"}))
}
