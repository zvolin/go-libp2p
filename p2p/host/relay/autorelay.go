package relay

import (
	"context"
	"fmt"
	"math/rand"
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
	relays map[peer.ID]struct{}
	status autonat.NATStatus
}

func NewAutoRelay(ctx context.Context, bhost *basic.BasicHost, discover discovery.Discoverer, router routing.PeerRouting) *AutoRelay {
	ar := &AutoRelay{
		host:       bhost,
		discover:   discover,
		router:     router,
		addrsF:     bhost.AddrsFactory,
		relays:     make(map[peer.ID]struct{}),
		disconnect: make(chan struct{}, 1),
		status:     autonat.NATStatusUnknown,
	}
	ar.autonat = autonat.NewAutoNAT(ctx, bhost, ar.baseAddrs)
	bhost.AddrsFactory = ar.hostAddrs
	bhost.Network().Notify(ar)
	go ar.background(ctx)
	return ar
}

func (ar *AutoRelay) baseAddrs() []ma.Multiaddr {
	return ar.addrsF(ar.host.AllAddrs())
}

func (ar *AutoRelay) hostAddrs(addrs []ma.Multiaddr) []ma.Multiaddr {
	return ar.relayAddrs(ar.addrsF(addrs))
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
			ar.mx.Lock()
			ar.status = autonat.NATStatusUnknown
			ar.mx.Unlock()
			wait = autonat.AutoNATRetryInterval

		case autonat.NATStatusPublic:
			ar.mx.Lock()
			if ar.status != autonat.NATStatusPublic {
				push = true
			}
			ar.status = autonat.NATStatusPublic
			ar.mx.Unlock()

		case autonat.NATStatusPrivate:
			update := ar.findRelays(ctx)
			ar.mx.Lock()
			if update || ar.status != autonat.NATStatusPrivate {
				push = true
			}
			ar.status = autonat.NATStatusPrivate
			ar.mx.Unlock()
		}

		if push {
			push = false
			ar.host.PushIdentify()
		}

		select {
		case <-ar.disconnect:
			push = true
		case <-time.After(wait):
		case <-ctx.Done():
			return
		}
	}
}

func (ar *AutoRelay) findRelays(ctx context.Context) bool {
	retry := 0

again:
	ar.mx.Lock()
	haveRelays := len(ar.relays)
	ar.mx.Unlock()
	if haveRelays >= DesiredRelays {
		return false
	}
	need := DesiredRelays - haveRelays

	limit := 1000

	dctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	pis, err := discovery.FindPeers(dctx, ar.discover, RelayRendezvous, limit)
	cancel()
	if err != nil {
		log.Debugf("error discovering relays: %s", err.Error())

		if haveRelays == 0 {
			retry++
			if retry > 5 {
				log.Debug("no relays connected; giving up")
				return false
			}

			log.Debug("no relays connected; retrying in 30s")
			select {
			case <-time.After(30 * time.Second):
				goto again
			case <-ctx.Done():
				return false
			}
		}
	}

	log.Debugf("discovered %d relays", len(pis))

	pis = ar.selectRelays(ctx, pis)
	update := 0

	for _, pi := range pis {
		ar.mx.Lock()
		if _, ok := ar.relays[pi.ID]; ok {
			ar.mx.Unlock()
			continue
		}
		ar.mx.Unlock()

		cctx, cancel := context.WithTimeout(ctx, 60*time.Second)

		if len(pi.Addrs) == 0 {
			pi, err = ar.router.FindPeer(cctx, pi.ID)
			if err != nil {
				log.Debugf("error finding relay peer %s: %s", pi.ID, err.Error())
				cancel()
				continue
			}
		}

		err = ar.host.Connect(cctx, pi)
		cancel()
		if err != nil {
			log.Debugf("error connecting to relay %s: %s", pi.ID, err.Error())
			continue
		}

		log.Debugf("connected to relay %s", pi.ID)
		ar.mx.Lock()
		ar.relays[pi.ID] = struct{}{}
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
		retry++
		if retry > 5 {
			log.Debug("no relays connected; giving up")
			return false
		}

		log.Debug("no relays connected; retrying in 30s")
		select {
		case <-time.After(30 * time.Second):
			goto again
		case <-ctx.Done():
			return false
		}
	}

	return update > 0
}

func (ar *AutoRelay) selectRelays(ctx context.Context, pis []pstore.PeerInfo) []pstore.PeerInfo {
	// TODO better relay selection strategy; this just selects random relays
	//      but we should probably use ping latency as the selection metric

	shuffleRelays(pis)
	return pis
}

// This function is computes the NATed relay addrs when our status is private:
// - The public addrs are removed from the address set.
// - The non-public addrs are included verbatim so that peers behind the same NAT/firewall
//   can still dial us directly.
// - On top of those, we add the relay-specific addrs for the relays to which we are
//   connected. For each non-private relay addr, we encapsulate the p2p-circuit addr
//   through which we can be dialed.
func (ar *AutoRelay) relayAddrs(addrs []ma.Multiaddr) []ma.Multiaddr {
	ar.mx.Lock()
	if ar.status != autonat.NATStatusPrivate {
		ar.mx.Unlock()
		return addrs
	}

	relays := make([]peer.ID, 0, len(ar.relays))
	for p := range ar.relays {
		relays = append(relays, p)
	}
	ar.mx.Unlock()

	raddrs := make([]ma.Multiaddr, 0, 4*len(relays)+2)

	// only keep private addrs from the original addr set
	for _, addr := range addrs {
		if manet.IsPrivateAddr(addr) {
			raddrs = append(raddrs, addr)
		}
	}

	// add relay specific addrs to the list
	for _, p := range relays {
		addrs := cleanupAddressSet(ar.host.Peerstore().Addrs(p))

		circuit, err := ma.NewMultiaddr(fmt.Sprintf("/p2p/%s/p2p-circuit", p.Pretty()))
		if err != nil {
			panic(err)
		}

		for _, addr := range addrs {
			pub := addr.Encapsulate(circuit)
			raddrs = append(raddrs, pub)
		}
	}

	return raddrs
}

func shuffleRelays(pis []pstore.PeerInfo) {
	for i := range pis {
		j := rand.Intn(i + 1)
		pis[i], pis[j] = pis[j], pis[i]
	}
}

// Notifee
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
