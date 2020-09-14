package conngater

import (
	"net"
	"sync"

	"github.com/libp2p/go-libp2p-core/connmgr"
	"github.com/libp2p/go-libp2p-core/control"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"

	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"

	logging "github.com/ipfs/go-log"
)

type BasicConnectionGater struct {
	sync.Mutex

	blockedPeers   map[peer.ID]struct{}
	blockedAddrs   map[string]struct{}
	blockedSubnets map[string]*net.IPNet
}

var log = logging.Logger("net/conngater")

func NewBasicConnectionGater() *BasicConnectionGater {
	// XXX
	return nil
}

// BlockPeer adds a peer to the set of blocked peers
func (cg *BasicConnectionGater) BlockPeer(p peer.ID) {
	cg.Lock()
	defer cg.Unlock()

	cg.blockedPeers[p] = struct{}{}
}

// UnblockPeer removes a peer from the set of blocked peers
func (cg *BasicConnectionGater) UnblockPeer(p peer.ID) {
	cg.Lock()
	defer cg.Unlock()

	delete(cg.blockedPeers, p)
}

// BlockAddr adds an IP address to the set of blocked addresses
func (cg *BasicConnectionGater) BlockAddr(ip net.IP) {
	cg.Lock()
	defer cg.Unlock()

	cg.blockedAddrs[ip.String()] = struct{}{}
}

// UnblockAddr removes an IP address from the set of blocked addresses
func (cg *BasicConnectionGater) UnblockAddr(ip net.IP) {
	cg.Lock()
	defer cg.Unlock()

	delete(cg.blockedAddrs, ip.String())
}

// BlockSubnet adds an IP subnet to the set of blocked addresses
func (cg *BasicConnectionGater) BlockSubnet(ipnet *net.IPNet) {
	cg.Lock()
	defer cg.Unlock()

	cg.blockedSubnets[ipnet.String()] = ipnet
}

// UnblockSubnet removes an IP address from the set of blocked addresses
func (cg *BasicConnectionGater) UnblockSubnet(ipnet *net.IPNet) {
	cg.Lock()
	defer cg.Unlock()

	delete(cg.blockedSubnets, ipnet.String())
}

// ConnectionGater interface
var _ connmgr.ConnectionGater = (*BasicConnectionGater)(nil)

func (cg *BasicConnectionGater) InterceptPeerDial(p peer.ID) (allow bool) {
	cg.Lock()
	defer cg.Unlock()

	_, block := cg.blockedPeers[p]
	return !block
}

func (cg *BasicConnectionGater) InterceptAddrDial(p peer.ID, a ma.Multiaddr) (allow bool) {
	cg.Lock()
	defer cg.Unlock()

	ip, err := manet.ToIP(a)
	if err != nil {
		log.Warnf("error converting multiaddr to IP addr: %s", err)
		return true
	}

	_, block := cg.blockedAddrs[ip.String()]
	if block {
		return false
	}

	for _, ipnet := range cg.blockedSubnets {
		if ipnet.Contains(ip) {
			return false
		}
	}

	return true
}

func (cg *BasicConnectionGater) InterceptAccept(cma network.ConnMultiaddrs) (allow bool) {
	cg.Lock()
	defer cg.Unlock()

	a := cma.RemoteMultiaddr()

	ip, err := manet.ToIP(a)
	if err != nil {
		log.Warnf("error converting multiaddr to IP addr: %s", err)
		return true
	}

	_, block := cg.blockedAddrs[ip.String()]
	if block {
		return false
	}

	for _, ipnet := range cg.blockedSubnets {
		if ipnet.Contains(ip) {
			return false
		}
	}

	return true
}

func (cg *BasicConnectionGater) InterceptSecured(dir network.Direction, p peer.ID, cma network.ConnMultiaddrs) (allow bool) {
	if dir == network.DirOutbound {
		// we have already filtered those in InterceptPeerDial/InterceptAddrDial
		return true
	}

	// we have already filtered addrs in InterceptAccept, so we just check the peer ID
	cg.Lock()
	defer cg.Unlock()

	_, block := cg.blockedPeers[p]
	return !block
}

func (cg *BasicConnectionGater) InterceptUpgraded(network.Conn) (allow bool, reason control.DisconnectReason) {
	return true, 0
}
