package basichost

import (
	"context"
	"sync"

	inat "github.com/libp2p/go-libp2p-nat"

	goprocess "github.com/jbenet/goprocess"
	lgbl "github.com/libp2p/go-libp2p-loggables"
	inet "github.com/libp2p/go-libp2p-net"
	ma "github.com/multiformats/go-multiaddr"
)

// NATManager takes care of adding + removing port mappings to the nat.
// Initialized with the host if it has a NATPortMap option enabled.
// NATManager receives signals from the network, and check on nat mappings:
//  * NATManager listens to the network and adds or closes port mappings
//    as the network signals Listen() or ListenClose().
//  * closing the NATManager closes the nat and its mappings.
type NATManager struct {
	net   inet.Network
	natmu sync.RWMutex // guards nat (ready could obviate this mutex, but safety first.)
	nat   *inat.NAT

	ready chan struct{}     // closed once the nat is ready to process port mappings
	proc  goprocess.Process // NATManager has a process + children. can be closed.
}

func NewNATManager(net inet.Network) *NATManager {
	nmgr := &NATManager{
		net:   net,
		ready: make(chan struct{}),
	}

	nmgr.proc = goprocess.WithTeardown(func() error {
		// on closing, unregister from network notifications.
		net.StopNotify((*nmgrNetNotifiee)(nmgr))
		return nil
	})

	// discover the nat.
	nmgr.discoverNAT()
	return nmgr
}

// Close closes the NATManager, closing the underlying nat
// and unregistering from network events.
func (nmgr *NATManager) Close() error {
	return nmgr.proc.Close()
}

// Ready returns a channel which will be closed when the NAT has been found
// and is ready to be used, or the search process is done.
func (nmgr *NATManager) Ready() <-chan struct{} {
	return nmgr.ready
}

func (nmgr *NATManager) discoverNAT() {

	nmgr.proc.Go(func(worker goprocess.Process) {
		// inat.DiscoverNAT blocks until the nat is found or a timeout
		// is reached. we unfortunately cannot specify timeouts-- the
		// library we're using just blocks.
		//
		// Note: on early shutdown, there may be a case where we're trying
		// to close before DiscoverNAT() returns. Since we cant cancel it
		// (library) we can choose to (1) drop the result and return early,
		// or (2) wait until it times out to exit. For now we choose (2),
		// to avoid leaking resources in a non-obvious way. the only case
		// this affects is when the daemon is being started up and _immediately_
		// asked to close. other services are also starting up, so ok to wait.
		discoverdone := make(chan struct{})
		var nat *inat.NAT
		go func() {
			defer close(discoverdone)
			nat = inat.DiscoverNAT()
		}()

		// by this point -- after finding the NAT -- we may have already
		// be closing. if so, just exit.
		select {
		case <-worker.Closing():
			return
		case <-discoverdone:
			if nat == nil { // no nat, or failed to get it.
				return
			}
		}

		// wire up the nat to close when nmgr closes.
		// nmgr.proc is our parent, and waiting for us.
		nmgr.proc.AddChild(nat.Process())

		// set the nat.
		nmgr.natmu.Lock()
		nmgr.nat = nat
		nmgr.natmu.Unlock()

		// signal that we're ready to process nat mappings:
		close(nmgr.ready)

		// sign NATManager up for network notifications
		// we need to sign up here to avoid missing some notifs
		// before the NAT has been found.
		nmgr.net.Notify((*nmgrNetNotifiee)(nmgr))

		// if any interfaces were brought up while we were setting up
		// the nat, now is the time to setup port mappings for them.
		// we release ready, then grab them to avoid losing any. adding
		// a port mapping is idempotent, so its ok to add the same twice.
		addrs := nmgr.net.ListenAddresses()
		for _, addr := range addrs {
			// we do it async because it's slow and we may want to close beforehand
			go addPortMapping(nmgr, addr)
		}
	})
}

// NAT returns the NATManager's nat object. this may be nil, if
// (a) the search process is still ongoing, or (b) the search process
// found no nat. Clients must check whether the return value is nil.
func (nmgr *NATManager) NAT() *inat.NAT {
	nmgr.natmu.Lock()
	defer nmgr.natmu.Unlock()
	return nmgr.nat
}

func addPortMapping(nmgr *NATManager, intaddr ma.Multiaddr) {
	nat := nmgr.NAT()
	if nat == nil {
		panic("NATManager addPortMapping called without a nat.")
	}

	// first, check if the port mapping already exists.
	for _, mapping := range nat.Mappings() {
		if mapping.InternalAddr().Equal(intaddr) {
			return // it exists! return.
		}
	}

	ctx := context.TODO()
	lm := make(lgbl.DeferredMap)
	lm["internalAddr"] = func() interface{} { return intaddr.String() }

	defer log.EventBegin(ctx, "natMgrAddPortMappingWait", lm).Done()

	select {
	case <-nmgr.proc.Closing():
		lm["outcome"] = "cancelled"
		return // no use.
	case <-nmgr.ready: // wait until it's ready.
	}

	// actually start the port map (sub-event because waiting may take a while)
	defer log.EventBegin(ctx, "natMgrAddPortMapping", lm).Done()

	// get the nat
	m, err := nat.NewMapping(intaddr)
	if err != nil {
		lm["outcome"] = "failure"
		lm["error"] = err
		return
	}

	extaddr, err := m.ExternalAddr()
	if err != nil {
		lm["outcome"] = "failure"
		lm["error"] = err
		return
	}

	lm["outcome"] = "success"
	lm["externalAddr"] = func() interface{} { return extaddr.String() }
	log.Infof("established nat port mapping: %s <--> %s", intaddr, extaddr)
}

func rmPortMapping(nmgr *NATManager, intaddr ma.Multiaddr) {
	nat := nmgr.NAT()
	if nat == nil {
		panic("NATManager rmPortMapping called without a nat.")
	}

	// list the port mappings (it may be gone on it's own, so we need to
	// check this list, and not store it ourselves behind the scenes)

	// close mappings for this internal address.
	for _, mapping := range nat.Mappings() {
		if mapping.InternalAddr().Equal(intaddr) {
			mapping.Close()
		}
	}
}

// nmgrNetNotifiee implements the network notification listening part
// of the NATManager. this is merely listening to Listen() and ListenClose()
// events.
type nmgrNetNotifiee NATManager

func (nn *nmgrNetNotifiee) natManager() *NATManager {
	return (*NATManager)(nn)
}

func (nn *nmgrNetNotifiee) Listen(n inet.Network, addr ma.Multiaddr) {
	if nn.natManager().NAT() == nil {
		return // not ready or doesnt exist.
	}

	addPortMapping(nn.natManager(), addr)
}

func (nn *nmgrNetNotifiee) ListenClose(n inet.Network, addr ma.Multiaddr) {
	if nn.natManager().NAT() == nil {
		return // not ready or doesnt exist.
	}

	rmPortMapping(nn.natManager(), addr)
}

func (nn *nmgrNetNotifiee) Connected(inet.Network, inet.Conn)      {}
func (nn *nmgrNetNotifiee) Disconnected(inet.Network, inet.Conn)   {}
func (nn *nmgrNetNotifiee) OpenedStream(inet.Network, inet.Stream) {}
func (nn *nmgrNetNotifiee) ClosedStream(inet.Network, inet.Stream) {}
