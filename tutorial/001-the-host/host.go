package main

import (
	"context"
	"fmt"

	pstore "github.com/libp2p/go-libp2p-peerstore"
	swarm "github.com/libp2p/go-libp2p-swarm"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	testutil "github.com/libp2p/go-testutil"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	// For toy applications, this is an easy way to get an identity
	ident, err := testutil.RandIdentity()
	if err != nil {
		panic(err)
	}

	// We've created the identity, now we need to store it.
	// A peerstore holds information about peers, including your own
	ps := pstore.NewPeerstore()

	// An identity is essentially a public/private keypair
	ps.AddPrivKey(ident.ID(), ident.PrivateKey())
	ps.AddPubKey(ident.ID(), ident.PublicKey())

	maddr, err := ma.NewMultiaddr("/ip4/0.0.0.0/tcp/9000")
	if err != nil {
		panic(err)
	}

	// Make a context to govern the lifespan of the swarm
	ctx := context.Background()

	// Put all this together
	netw, err := swarm.NewNetwork(ctx, []ma.Multiaddr{maddr}, ident.ID(), ps, nil)
	if err != nil {
		panic(err)
	}

	myhost := bhost.New(netw)
	fmt.Printf("Hello World, my hosts ID is %s\n", myhost.ID())
}
