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

	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
	"github.com/ipfs/go-datastore/query"
	logging "github.com/ipfs/go-log"
)

type BasicConnectionGater struct {
	sync.RWMutex

	blockedPeers   map[peer.ID]struct{}
	blockedAddrs   map[string]struct{}
	blockedSubnets map[string]*net.IPNet

	ds datastore.Datastore
}

var log = logging.Logger("net/conngater")

const ns = "/libp2p/net/conngater"
const keyPeer = "/peer/"
const keyAddr = "/addr/"
const keySubnet = "/subnet/"

// NewBasicConnectionGater creates a new connection gater.
// The ds argument is an (optional, can be nil) datastore to persist the connection gater
// filters
func NewBasicConnectionGater(ds datastore.Datastore) *BasicConnectionGater {
	cg := &BasicConnectionGater{
		blockedPeers:   make(map[peer.ID]struct{}),
		blockedAddrs:   make(map[string]struct{}),
		blockedSubnets: make(map[string]*net.IPNet),
	}

	if ds != nil {
		cg.ds = namespace.Wrap(ds, datastore.NewKey(ns))
		cg.loadRules()
	}

	return cg
}

func (cg *BasicConnectionGater) loadRules() {
	// load blocked peers
	res, err := cg.ds.Query(query.Query{Prefix: keyPeer})
	if err != nil {
		log.Errorf("error querying datastore for blocked peers: %s", err)
		return
	}

	for r := range res.Next() {
		if r.Error != nil {
			log.Errorf("query result error: %s", r.Error)
			return
		}

		p := peer.ID(r.Entry.Value)
		cg.blockedPeers[p] = struct{}{}
	}

	// load blocked addrs
	res, err = cg.ds.Query(query.Query{Prefix: keyAddr})
	if err != nil {
		log.Errorf("error querying datastore for blocked addrs: %s", err)
		return
	}

	for r := range res.Next() {
		if r.Error != nil {
			log.Errorf("query result error: %s", r.Error)
			return
		}

		ip := net.IP(r.Entry.Value)
		cg.blockedAddrs[ip.String()] = struct{}{}
	}

	// load blocked subnets
	res, err = cg.ds.Query(query.Query{Prefix: keySubnet})
	if err != nil {
		log.Errorf("error querying datastore for blocked subnets: %s", err)
		return
	}

	for r := range res.Next() {
		if r.Error != nil {
			log.Errorf("query result error: %s", r.Error)
			return
		}

		ipnetStr := string(r.Entry.Value)
		_, ipnet, err := net.ParseCIDR(ipnetStr)
		if err != nil {
			log.Errorf("error parsing CIDR subnet: %s", err)
			return
		}
		cg.blockedSubnets[ipnetStr] = ipnet
	}
}

// BlockPeer adds a peer to the set of blocked peers
func (cg *BasicConnectionGater) BlockPeer(p peer.ID) {
	cg.Lock()
	defer cg.Unlock()

	cg.blockedPeers[p] = struct{}{}

	if cg.ds != nil {
		err := cg.ds.Put(datastore.NewKey(keyPeer+p.String()), []byte(p))
		if err != nil {
			log.Errorf("error writing blocked peer to datastore: %s", err)
		}
	}
}

// UnblockPeer removes a peer from the set of blocked peers
func (cg *BasicConnectionGater) UnblockPeer(p peer.ID) {
	cg.Lock()
	defer cg.Unlock()

	delete(cg.blockedPeers, p)

	if cg.ds != nil {
		err := cg.ds.Delete(datastore.NewKey(keyPeer + p.String()))
		if err != nil {
			log.Errorf("error deleting blocked peer from datastore: %s", err)
		}
	}
}

// ListBlockedPeers return a list of blocked peers
func (cg *BasicConnectionGater) ListBlockedPeers() []peer.ID {
	cg.RLock()
	defer cg.RUnlock()

	result := make([]peer.ID, 0, len(cg.blockedPeers))
	for p := range cg.blockedPeers {
		result = append(result, p)
	}

	return result
}

// BlockAddr adds an IP address to the set of blocked addresses
func (cg *BasicConnectionGater) BlockAddr(ip net.IP) {
	cg.Lock()
	defer cg.Unlock()

	cg.blockedAddrs[ip.String()] = struct{}{}

	if cg.ds != nil {
		err := cg.ds.Put(datastore.NewKey(keyAddr+ip.String()), []byte(ip))
		if err != nil {
			log.Errorf("error writing blocked addr to datastore: %s", err)
		}
	}
}

// UnblockAddr removes an IP address from the set of blocked addresses
func (cg *BasicConnectionGater) UnblockAddr(ip net.IP) {
	cg.Lock()
	defer cg.Unlock()

	delete(cg.blockedAddrs, ip.String())

	if cg.ds != nil {
		err := cg.ds.Delete(datastore.NewKey(keyAddr + ip.String()))
		if err != nil {
			log.Errorf("error deleting blocked addr from datastore: %s", err)
		}
	}
}

// ListBlockedAddrs return a list of blocked IP addresses
func (cg *BasicConnectionGater) ListBlockedAddrs() []net.IP {
	cg.RLock()
	defer cg.RUnlock()

	result := make([]net.IP, 0, len(cg.blockedAddrs))
	for ipStr := range cg.blockedAddrs {
		ip := net.ParseIP(ipStr)
		result = append(result, ip)
	}

	return result
}

// BlockSubnet adds an IP subnet to the set of blocked addresses
func (cg *BasicConnectionGater) BlockSubnet(ipnet *net.IPNet) {
	cg.Lock()
	defer cg.Unlock()

	cg.blockedSubnets[ipnet.String()] = ipnet

	if cg.ds != nil {
		err := cg.ds.Put(datastore.NewKey(keySubnet+ipnet.String()), []byte(ipnet.String()))
		if err != nil {
			log.Errorf("error writing blocked addr to datastore: %s", err)
		}
	}
}

// UnblockSubnet removes an IP address from the set of blocked addresses
func (cg *BasicConnectionGater) UnblockSubnet(ipnet *net.IPNet) {
	cg.Lock()
	defer cg.Unlock()

	delete(cg.blockedSubnets, ipnet.String())

	if cg.ds != nil {
		err := cg.ds.Delete(datastore.NewKey(keySubnet + ipnet.String()))
		if err != nil {
			log.Errorf("error deleting blocked subnet from datastore: %s", err)
		}
	}
}

// ListBlockedSubnets return a list of blocked IP subnets
func (cg *BasicConnectionGater) ListBlockedSubnets() []*net.IPNet {
	cg.RLock()
	defer cg.RUnlock()

	result := make([]*net.IPNet, 0, len(cg.blockedSubnets))
	for _, ipnet := range cg.blockedSubnets {
		result = append(result, ipnet)
	}

	return result
}

// ConnectionGater interface
var _ connmgr.ConnectionGater = (*BasicConnectionGater)(nil)

func (cg *BasicConnectionGater) InterceptPeerDial(p peer.ID) (allow bool) {
	cg.RLock()
	defer cg.RUnlock()

	_, block := cg.blockedPeers[p]
	return !block
}

func (cg *BasicConnectionGater) InterceptAddrDial(p peer.ID, a ma.Multiaddr) (allow bool) {
	// we have already filrted blocked peersin InterceptPeerDial, so we just check the IP
	cg.RLock()
	defer cg.RUnlock()

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
	cg.RLock()
	defer cg.RUnlock()

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
	cg.RLock()
	defer cg.RUnlock()

	_, block := cg.blockedPeers[p]
	return !block
}

func (cg *BasicConnectionGater) InterceptUpgraded(network.Conn) (allow bool, reason control.DisconnectReason) {
	return true, 0
}
