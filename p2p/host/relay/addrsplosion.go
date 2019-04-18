package relay

import (
	"encoding/binary"

	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

// This function cleans up a relay's address set to remove private addresses and curtail
// addrsplosion.
func cleanupAddressSet(pi pstore.PeerInfo) pstore.PeerInfo {
	var public, private []ma.Multiaddr

	for _, a := range pi.Addrs {
		if manet.IsPublicAddr(a) || isDNSAddr(a) {
			public = append(public, a)
			continue
		}

		// discard unroutable addrs
		if manet.IsPrivateAddr(a) {
			private = append(private, a)
		}
	}

	if !hasAddrsplosion(public) {
		return pstore.PeerInfo{ID: pi.ID, Addrs: public}
	}

	addrs := sanitizeAddrsplodedSet(public, private)
	return pstore.PeerInfo{ID: pi.ID, Addrs: addrs}
}

func isDNSAddr(a ma.Multiaddr) bool {
	dnsaddr := false
	ma.ForEach(a, func(c ma.Component) bool {
		switch c.Protocol().Code {
		case 54, 55, 56:
			dnsaddr = true
		}

		return false
	})

	return dnsaddr
}

// we have addrsplosion if for some protocol we advertise multiple ports on
// the same base address.
func hasAddrsplosion(addrs []ma.Multiaddr) bool {
	aset := make(map[string]int)

	for _, a := range addrs {
		key, port := addrKeyAndPort(a)
		xport, ok := aset[key]
		if ok && port != xport {
			return true
		}
		aset[key] = port
	}

	return false
}

func addrKeyAndPort(a ma.Multiaddr) (string, int) {
	var (
		key  string
		port int
	)

	ma.ForEach(a, func(c ma.Component) bool {
		switch c.Protocol().Code {
		case ma.P_TCP, ma.P_UDP:
			port = int(binary.BigEndian.Uint16(c.RawValue()))
			key += "/" + c.Protocol().Name
		default:
			val := c.Value()
			if val == "" {
				val = c.Protocol().Name
			}
			key += "/" + val
		}
		return true
	})

	return key, port
}

// clean up addrsplosion
// the following heuristic is used:
// - for each base address/protocol combination, if there are multiple ports advertised then
//   only accept the default port if present.
//   - If the default port is not present, we check for non-standard ports by tracking
//     private port bindings if present.
//   - If there is no default or private port binding, then we can't infer the correct
//     port and give up and return all addrs (for that base address)
func sanitizeAddrsplodedSet(public, private []ma.Multiaddr) []ma.Multiaddr {
	type portAndAddr struct {
		addr ma.Multiaddr
		port int
	}

	privports := make(map[int]struct{})
	pubaddrs := make(map[string][]portAndAddr)

	for _, a := range private {
		_, port := addrKeyAndPort(a)
		privports[port] = struct{}{}
	}

	for _, a := range public {
		key, port := addrKeyAndPort(a)
		pubaddrs[key] = append(pubaddrs[key], portAndAddr{addr: a, port: port})
	}

	var result []ma.Multiaddr
	for _, pas := range pubaddrs {
		if len(pas) == 1 {
			result = append(result, pas[0].addr)
			continue
		}

		for _, pa := range pas {
			if pa.port == 4001 || pa.port == 4002 {
				result = append(result, pa.addr)
				continue
			}
			if _, ok := privports[pa.port]; ok {
				result = append(result, pa.addr)
			}
		}
	}

	return result
}
