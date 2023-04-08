package basichost

import (
	"context"
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	inat "github.com/libp2p/go-libp2p/p2p/net/nat"

	ma "github.com/multiformats/go-multiaddr"
)

// NATManager is a simple interface to manage NAT devices.
// It listens Listen and ListenClose notifications from the network.Network,
// and tries to obtain port mappings for those.
type NATManager interface {
	// NAT gets the NAT device managed by the NAT manager.
	NAT() *inat.NAT

	io.Closer
}

// NewNATManager creates a NAT manager.
func NewNATManager(net network.Network) NATManager {
	return newNatManager(net)
}

type entry struct {
	protocol string
	port     int
}

// natManager takes care of adding + removing port mappings to the nat.
// Initialized with the host if it has a NATPortMap option enabled.
// natManager receives signals from the network, and check on nat mappings:
//   - natManager listens to the network and adds or closes port mappings
//     as the network signals Listen() or ListenClose().
//   - closing the natManager closes the nat and its mappings.
type natManager struct {
	net   network.Network
	natMx sync.RWMutex
	nat   *inat.NAT

	syncFlag chan struct{} // cap: 1

	tracked map[entry]bool // the bool is only used in doSync and has no meaning outside of that function

	refCount  sync.WaitGroup
	ctxCancel context.CancelFunc
}

func newNatManager(net network.Network) *natManager {
	ctx, cancel := context.WithCancel(context.Background())
	nmgr := &natManager{
		net:       net,
		syncFlag:  make(chan struct{}, 1),
		ctxCancel: cancel,
		tracked:   make(map[entry]bool),
	}
	nmgr.refCount.Add(1)
	go nmgr.background(ctx)
	return nmgr
}

// Close closes the natManager, closing the underlying nat
// and unregistering from network events.
func (nmgr *natManager) Close() error {
	nmgr.ctxCancel()
	nmgr.refCount.Wait()
	return nil
}

func (nmgr *natManager) background(ctx context.Context) {
	defer nmgr.refCount.Done()

	defer func() {
		nmgr.natMx.Lock()
		defer nmgr.natMx.Unlock()

		if nmgr.nat != nil {
			nmgr.nat.Close()
		}
	}()

	discoverCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	natInstance, err := inat.DiscoverNAT(discoverCtx)
	if err != nil {
		log.Info("DiscoverNAT error:", err)
		return
	}

	nmgr.natMx.Lock()
	nmgr.nat = natInstance
	nmgr.natMx.Unlock()

	// sign natManager up for network notifications
	// we need to sign up here to avoid missing some notifs
	// before the NAT has been found.
	nmgr.net.Notify((*nmgrNetNotifiee)(nmgr))
	defer nmgr.net.StopNotify((*nmgrNetNotifiee)(nmgr))

	nmgr.doSync() // sync one first.
	for {
		select {
		case <-nmgr.syncFlag:
			nmgr.doSync() // sync when our listen addresses chnage.
		case <-ctx.Done():
			return
		}
	}
}

func (nmgr *natManager) sync() {
	select {
	case nmgr.syncFlag <- struct{}{}:
	default:
	}
}

// doSync syncs the current NAT mappings, removing any outdated mappings and adding any
// new mappings.
func (nmgr *natManager) doSync() {
	for e := range nmgr.tracked {
		nmgr.tracked[e] = false
	}
	var newAddresses []entry
	for _, maddr := range nmgr.net.ListenAddresses() {
		// Strip the IP
		maIP, rest := ma.SplitFirst(maddr)
		if maIP == nil || rest == nil {
			continue
		}

		switch maIP.Protocol().Code {
		case ma.P_IP6, ma.P_IP4:
		default:
			continue
		}

		// Only bother if we're listening on an unicast / unspecified IP.
		ip := net.IP(maIP.RawValue())
		if !ip.IsGlobalUnicast() && !ip.IsUnspecified() {
			continue
		}

		// Extract the port/protocol
		proto, _ := ma.SplitFirst(rest)
		if proto == nil {
			continue
		}

		var protocol string
		switch proto.Protocol().Code {
		case ma.P_TCP:
			protocol = "tcp"
		case ma.P_UDP:
			protocol = "udp"
		default:
			continue
		}
		port, err := strconv.ParseUint(proto.Value(), 10, 16)
		if err != nil {
			// bug in multiaddr
			panic(err)
		}
		e := entry{protocol: protocol, port: int(port)}
		if _, ok := nmgr.tracked[e]; ok {
			nmgr.tracked[e] = true
		} else {
			newAddresses = append(newAddresses, e)
		}
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	// Close old mappings
	for e, v := range nmgr.tracked {
		if !v {
			nmgr.nat.RemoveMapping(e.protocol, e.port)
			delete(nmgr.tracked, e)
		}
	}

	// Create new mappings.
	for _, e := range newAddresses {
		if err := nmgr.nat.AddMapping(e.protocol, e.port); err != nil {
			log.Errorf("failed to port-map %s port %d: %s", e.protocol, e.port, err)
		}
		nmgr.tracked[e] = false
	}
}

// NAT returns the natManager's nat object. this may be nil, if
// (a) the search process is still ongoing, or (b) the search process
// found no nat. Clients must check whether the return value is nil.
func (nmgr *natManager) NAT() *inat.NAT {
	nmgr.natMx.Lock()
	defer nmgr.natMx.Unlock()
	return nmgr.nat
}

type nmgrNetNotifiee natManager

func (nn *nmgrNetNotifiee) natManager() *natManager {
	return (*natManager)(nn)
}

func (nn *nmgrNetNotifiee) Listen(network.Network, ma.Multiaddr) {
	nn.natManager().sync()
}

func (nn *nmgrNetNotifiee) ListenClose(network.Network, ma.Multiaddr) {
	nn.natManager().sync()
}

func (nn *nmgrNetNotifiee) Connected(network.Network, network.Conn)    {}
func (nn *nmgrNetNotifiee) Disconnected(network.Network, network.Conn) {}
