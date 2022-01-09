module github.com/libp2p/go-libp2p/examples

go 1.16

require (
	github.com/gogo/protobuf v1.3.2
	github.com/google/uuid v1.3.0
	github.com/ipfs/go-datastore v0.5.0
	github.com/ipfs/go-log/v2 v2.5.0
	github.com/libp2p/go-libp2p v0.14.4
	github.com/libp2p/go-libp2p-connmgr v0.2.4
	github.com/libp2p/go-libp2p-core v0.13.1-0.20220104083644-a3dd401efe36
	github.com/libp2p/go-libp2p-discovery v0.6.0
	github.com/libp2p/go-libp2p-kad-dht v0.15.0
	github.com/libp2p/go-libp2p-noise v0.3.0
	github.com/libp2p/go-libp2p-swarm v0.9.1-0.20220104091227-f776b7e504b1
	github.com/libp2p/go-libp2p-tls v0.3.1
	github.com/multiformats/go-multiaddr v0.5.0
)

// Ensure that examples always use the go-libp2p version in the same git checkout.
replace github.com/libp2p/go-libp2p => ../
