package relay

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	basic "github.com/libp2p/go-libp2p/p2p/host/basic"

	autonat "github.com/libp2p/go-libp2p-autonat"
	_ "github.com/libp2p/go-libp2p-circuit"
	discovery "github.com/libp2p/go-libp2p-discovery"
	inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	routing "github.com/libp2p/go-libp2p-routing"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

const (
	RelayRendezvous = "/libp2p/relay"
)

var (
	DesiredRelays = 3

	BootDelay = 20 * time.Second
)

// AutoRelay is a Host that uses relays for connectivity when a NAT is detected.
type AutoRelay struct {
	host     *basic.BasicHost
	discover discovery.Discoverer
	router   routing.PeerRouting
	autonat  autonat.AutoNAT
	addrsF   basic.AddrsFactory

	disconnect chan struct{}

	mx     sync.Mutex
	relays map[peer.ID]pstore.PeerInfo
	addrs  []ma.Multiaddr
}

func NewAutoRelay(ctx context.Context, bhost *basic.BasicHost, discover discovery.Discoverer, router routing.PeerRouting) *AutoRelay {
	ar := &AutoRelay{
		host:       bhost,
		discover:   discover,
		router:     router,
		addrsF:     bhost.AddrsFactory,
		relays:     make(map[peer.ID]pstore.PeerInfo),
		disconnect: make(chan struct{}, 1),
	}
	ar.autonat = autonat.NewAutoNAT(ctx, bhost, ar.baseAddrs)
	bhost.AddrsFactory = ar.hostAddrs
	bhost.Network().Notify(ar)
	go ar.background(ctx)
	return ar
}

func (ar *AutoRelay) hostAddrs(addrs []ma.Multiaddr) []ma.Multiaddr {
	ar.mx.Lock()
	defer ar.mx.Unlock()
	if ar.addrs != nil && ar.autonat.Status() == autonat.NATStatusPrivate {
		return ar.addrs
	} else {
		return ar.addrsF(addrs)
	}
}

func (ar *AutoRelay) baseAddrs() []ma.Multiaddr {
	return ar.addrsF(ar.host.AllAddrs())
}

func (ar *AutoRelay) background(ctx context.Context) {
	select {
	case <-time.After(autonat.AutoNATBootDelay + BootDelay):
	case <-ctx.Done():
		return
	}

	// when true, we need to identify push
	push := false

	for {
		wait := autonat.AutoNATRefreshInterval
		switch ar.autonat.Status() {
		case autonat.NATStatusUnknown:
			wait = autonat.AutoNATRetryInterval

		case autonat.NATStatusPublic:
			// invalidate addrs
			ar.mx.Lock()
			if ar.addrs != nil {
				ar.addrs = nil
				push = true
			}
			ar.mx.Unlock()

			// if we had previously announced relay addrs, push our public addrs
			if push {
				push = false
				ar.host.PushIdentify()
			}

		case autonat.NATStatusPrivate:
			push = false // clear, findRelays pushes as needed
			ar.findRelays(ctx)
		}

		select {
		case <-ar.disconnect:
			// invalidate addrs
			ar.mx.Lock()
			if ar.addrs != nil {
				ar.addrs = nil
				push = true
			}
			ar.mx.Unlock()
		case <-time.After(wait):
		case <-ctx.Done():
			return
		}
	}
}

func (ar *AutoRelay) findRelays(ctx context.Context) {
again:
	ar.mx.Lock()
	haveRelays := len(ar.relays)
	if haveRelays >= DesiredRelays {
		ar.mx.Unlock()
		// this dance is necessary to cover the Private->Public->Private transition
		// where we were already connected to enough relays while private and dropped
		// the addrs while public
		if ar.addrs == nil {
			ar.updateAddrs()
		}
		return
	}
	need := DesiredRelays - len(ar.relays)
	ar.mx.Unlock()

	limit := 1000

	dctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	pis, err := discovery.FindPeers(dctx, ar.discover, RelayRendezvous, limit)
	cancel()
	if err != nil {
		log.Debugf("error discovering relays: %s", err.Error())
		return
	}

	pis = ar.selectRelays(ctx, pis, 20, 50)
	update := 0

	for _, pi := range pis {
		ar.mx.Lock()
		if _, ok := ar.relays[pi.ID]; ok {
			ar.mx.Unlock()
			continue
		}
		ar.mx.Unlock()

		cctx, cancel := context.WithTimeout(ctx, 15*time.Second)
		err = ar.host.Connect(cctx, pi)
		cancel()
		if err != nil {
			log.Debugf("error connecting to relay %s: %s", pi.ID, err.Error())
			continue
		}

		log.Debugf("connected to relay %s", pi.ID)
		ar.mx.Lock()
		ar.relays[pi.ID] = pi
		haveRelays++
		ar.mx.Unlock()

		// tag the connection as very important
		ar.host.ConnManager().TagPeer(pi.ID, "relay", 42)

		update++
		need--
		if need == 0 {
			break
		}
	}

	if haveRelays == 0 {
		// we failed to find any relays and we are not connected to any!
		// wait a little and try again, the discovery query might have returned only dead peers
		select {
		case <-time.After(30 * time.Second):
			goto again
		case <-ctx.Done():
			return
		}
	}

	if update > 0 || ar.addrs == nil {
		ar.updateAddrs()
	}
}

func (ar *AutoRelay) selectRelays(ctx context.Context, pis []pstore.PeerInfo, count, maxq int) []pstore.PeerInfo {
	// TODO better relay selection strategy; this just selects random relays
	//      but we should probably use ping latency as the selection metric

	if len(pis) == 0 {
		log.Debugf("no relays discovered")
		return pis
	}

	if len(pis) < count {
		count = len(pis)
	}

	if len(pis) < maxq {
		maxq = len(pis)
	}

	// only select relays that can be found by routing
	type queryResult struct {
		pi  pstore.PeerInfo
		err error
	}
	result := make([]pstore.PeerInfo, 0, count)
	resultCh := make(chan queryResult, maxq)

	qctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// shuffle to randomize the order of queries
	shuffleRelays(pis)
	for _, pi := range pis[:maxq] {
		go func(p peer.ID) {
			pi, err := ar.router.FindPeer(qctx, p)
			if err != nil {
				log.Debugf("error finding relay peer %s: %s", p, err.Error())
			}
			resultCh <- queryResult{pi: pi, err: err}
		}(pi.ID)
	}

	rcount := 0
	for len(result) < count && rcount < maxq {
		select {
		case qr := <-resultCh:
			rcount++
			if qr.err == nil {
				pi := cleanupAddressSet(qr.pi)
				if len(pi.Addrs) > 0 {
					result = append(result, pi)
				} else {
					log.Debugf("ignoring relay peer %s: cleaned up address set is empty", pi.ID)
				}
			}

		case <-qctx.Done():
			break
		}
	}

	shuffleRelays(result)
	return result
}

func (ar *AutoRelay) updateAddrs() {
	ar.doUpdateAddrs()
	ar.host.PushIdentify()
}

// This function updates our NATed advertised addrs (ar.addrs)
// The public addrs are rewritten so that they only retain the public IP part; they
// become undialable but are useful as a hint to the dialer to determine whether or not
// to dial private addrs.
// The non-public addrs are included verbatim so that peers behind the same NAT/firewall
// can still dial us directly.
// On top of those, we add the relay-specific addrs for the relays to which we are
// connected. For each non-private relay addr, we encapsulate the p2p-circuit addr
// through which we can be dialed.
func (ar *AutoRelay) doUpdateAddrs() {
	ar.mx.Lock()
	defer ar.mx.Unlock()

	addrs := ar.baseAddrs()
	raddrs := make([]ma.Multiaddr, 0, len(addrs)+len(ar.relays))

	// remove our public addresses from the list
	for _, addr := range addrs {
		if manet.IsPublicAddr(addr) {
			continue
		}
		raddrs = append(raddrs, addr)
	}

	// add relay specific addrs to the list
	for _, pi := range ar.relays {
		circuit, err := ma.NewMultiaddr(fmt.Sprintf("/p2p/%s/p2p-circuit", pi.ID.Pretty()))
		if err != nil {
			panic(err)
		}

		for _, addr := range pi.Addrs {
			if !manet.IsPrivateAddr(addr) {
				pub := addr.Encapsulate(circuit)
				raddrs = append(raddrs, pub)
			}
		}
	}

	ar.addrs = raddrs
}

// This function cleans up a relay's address set to remove private addresses and curtail
// addrsplosion. For the latter, we use the following heuristic:
// - if the address set includes a (tcp) address with the default port 4001,
//   we remove all tcp addrs with a different port
func cleanupAddressSet(pi pstore.PeerInfo) pstore.PeerInfo {
	// pass-1: find default port
	has4001 := false
	for _, addr := range pi.Addrs {
		port, err := tcpPort(addr)
		if err != nil {
			continue
		}
		if port == 4001 {
			has4001 = true
			break
		}
	}

	// pass-2: cleanup
	var newAddrs []ma.Multiaddr

	for _, addr := range pi.Addrs {
		if manet.IsPrivateAddr(addr) {
			continue
		}

		if has4001 {
			port, err := tcpPort(addr)
			if err == nil && port != 4001 {
				continue
			}
		}

		newAddrs = append(newAddrs, addr)
	}

	return pstore.PeerInfo{ID: pi.ID, Addrs: newAddrs}
}

func tcpPort(addr ma.Multiaddr) (int, error) {
	val, err := addr.ValueForProtocol(ma.P_TCP)
	if err != nil {
		// not tcp
		return 0, err
	}
	return strconv.Atoi(val)
}

func shuffleRelays(pis []pstore.PeerInfo) {
	for i := range pis {
		j := rand.Intn(i + 1)
		pis[i], pis[j] = pis[j], pis[i]
	}
}

func (ar *AutoRelay) Listen(inet.Network, ma.Multiaddr)      {}
func (ar *AutoRelay) ListenClose(inet.Network, ma.Multiaddr) {}
func (ar *AutoRelay) Connected(inet.Network, inet.Conn)      {}

func (ar *AutoRelay) Disconnected(net inet.Network, c inet.Conn) {
	p := c.RemotePeer()

	ar.mx.Lock()
	defer ar.mx.Unlock()

	if ar.host.Network().Connectedness(p) == inet.Connected {
		// We have a second connection.
		return
	}

	if _, ok := ar.relays[p]; ok {
		delete(ar.relays, p)
		select {
		case ar.disconnect <- struct{}{}:
		default:
		}
	}
}

func (ar *AutoRelay) OpenedStream(inet.Network, inet.Stream) {}
func (ar *AutoRelay) ClosedStream(inet.Network, inet.Stream) {}
