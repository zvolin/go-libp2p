package relay

import (
	"context"

	discovery "github.com/libp2p/go-libp2p-discovery"
	host "github.com/libp2p/go-libp2p-host"
)

// RelayHost is a Host that provides Relay services.
type RelayHost struct {
	host.Host
	advertise discovery.Advertiser
}

// New constructs a new RelayHost
func NewRelayHost(ctx context.Context, host host.Host, advertise discovery.Advertiser) *RelayHost {
	h := &RelayHost{Host: host, advertise: advertise}
	discovery.Advertise(ctx, advertise, "/libp2p/relay")
	return h
}

var _ host.Host = (*RelayHost)(nil)
