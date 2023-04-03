package nat

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	logging "github.com/ipfs/go-log/v2"

	"github.com/libp2p/go-nat"
)

// ErrNoMapping signals no mapping exists for an address
var ErrNoMapping = errors.New("mapping not established")

var log = logging.Logger("nat")

// MappingDuration is a default port mapping duration.
// Port mappings are renewed every (MappingDuration / 3)
const MappingDuration = time.Second * 60

// CacheTime is the time a mapping will cache an external address for
const CacheTime = time.Second * 15

type entry struct {
	protocol string
	port     int
}

// DiscoverNAT looks for a NAT device in the network and
// returns an object that can manage port mappings.
func DiscoverNAT(ctx context.Context) (*NAT, error) {
	natInstance, err := nat.DiscoverGateway(ctx)
	if err != nil {
		return nil, err
	}

	// Log the device addr.
	addr, err := natInstance.GetDeviceAddress()
	if err != nil {
		log.Debug("DiscoverGateway address error:", err)
	} else {
		log.Debug("DiscoverGateway address:", addr)
	}

	ctx, cancel := context.WithCancel(context.Background())
	nat := &NAT{
		nat:       natInstance,
		mappings:  make(map[entry]int),
		ctx:       ctx,
		ctxCancel: cancel,
	}
	nat.refCount.Add(1)
	go func() {
		defer nat.refCount.Done()
		nat.background()
	}()
	return nat, nil
}

// NAT is an object that manages address port mappings in
// NATs (Network Address Translators). It is a long-running
// service that will periodically renew port mappings,
// and keep an up-to-date list of all the external addresses.
type NAT struct {
	natmu sync.Mutex
	nat   nat.NAT

	refCount  sync.WaitGroup
	ctx       context.Context
	ctxCancel context.CancelFunc

	mappingmu sync.RWMutex // guards mappings
	closed    bool
	mappings  map[entry]int
}

// Close shuts down all port mappings. NAT can no longer be used.
func (nat *NAT) Close() error {
	nat.mappingmu.Lock()
	nat.closed = true
	nat.mappingmu.Unlock()

	nat.ctxCancel()
	nat.refCount.Wait()
	return nil
}

// Mappings returns a slice of all NAT mappings
func (nat *NAT) Mappings() []Mapping {
	nat.mappingmu.Lock()
	defer nat.mappingmu.Unlock()
	maps2 := make([]Mapping, 0, len(nat.mappings))
	for e, extPort := range nat.mappings {
		maps2 = append(maps2, &mapping{
			nat:     nat,
			proto:   e.protocol,
			intport: e.port,
			extport: extPort,
		})
	}
	return maps2
}

// AddMapping attempts to construct a mapping on protocol and internal port
// It will also periodically renew the mapping.
//
// May not succeed, and mappings may change over time;
// NAT devices may not respect our port requests, and even lie.
func (nat *NAT) AddMapping(protocol string, port int) error {
	switch protocol {
	case "tcp", "udp":
	default:
		return fmt.Errorf("invalid protocol: %s", protocol)
	}

	nat.mappingmu.Lock()
	if nat.closed {
		nat.mappingmu.Unlock()
		return errors.New("closed")
	}

	// do it once synchronously, so first mapping is done right away, and before exiting,
	// allowing users -- in the optimistic case -- to use results right after.
	extPort := nat.establishMapping(protocol, port)
	nat.mappings[entry{protocol: protocol, port: port}] = extPort
	nat.mappingmu.Unlock()

	return nil
}

func (nat *NAT) RemoveMapping(protocol string, port int) error {
	nat.mappingmu.Lock()
	defer nat.mappingmu.Unlock()
	switch protocol {
	case "tcp", "udp":
		delete(nat.mappings, entry{protocol: protocol, port: port})
	default:
		return fmt.Errorf("invalid protocol: %s", protocol)
	}
	return nil
}

func (nat *NAT) background() {
	const tick = MappingDuration / 3
	t := time.NewTimer(tick) // don't use a ticker here. We don't know how long establishing the mappings takes.
	defer t.Stop()

	var in []entry
	var out []int // port numbers
	for {
		select {
		case <-t.C:
			in = in[:0]
			out = out[:0]
			nat.mappingmu.Lock()
			for e := range nat.mappings {
				in = append(in, e)
			}
			nat.mappingmu.Unlock()
			// Establishing the mapping involves network requests.
			// Don't hold the mutex, just save the ports.
			for _, e := range in {
				out = append(out, nat.establishMapping(e.protocol, e.port))
			}
			nat.mappingmu.Lock()
			for i, p := range in {
				if _, ok := nat.mappings[p]; !ok {
					continue // entry might have been deleted
				}
				nat.mappings[p] = out[i]
			}
			nat.mappingmu.Unlock()
			t.Reset(tick)
		case <-nat.ctx.Done():
			nat.mappingmu.Lock()
			for e := range nat.mappings {
				delete(nat.mappings, e)
			}
			nat.mappingmu.Unlock()
			return
		}
	}
}

func (nat *NAT) establishMapping(protocol string, internalPort int) (externalPort int) {
	log.Debugf("Attempting port map: %s/%d", protocol, internalPort)
	const comment = "libp2p"

	nat.natmu.Lock()
	var err error
	externalPort, err = nat.nat.AddPortMapping(protocol, internalPort, comment, MappingDuration)
	if err != nil {
		// Some hardware does not support mappings with timeout, so try that
		externalPort, err = nat.nat.AddPortMapping(protocol, internalPort, comment, 0)
	}
	nat.natmu.Unlock()

	if err != nil || externalPort == 0 {
		// TODO: log.Event
		if err != nil {
			log.Warnf("failed to establish port mapping: %s", err)
		} else {
			log.Warnf("failed to establish port mapping: newport = 0")
		}
		// we do not close if the mapping failed,
		// because it may work again next time.
		return 0
	}

	log.Debugf("NAT Mapping: %d --> %d (%s)", externalPort, internalPort, protocol)
	return externalPort
}
