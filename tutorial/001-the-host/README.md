# The libp2p 'host'
For most applications, the host is the basic building block you'll need to get
started. This guide will show how to construct and use a simple host.
The host is an abstraction that manages services on top of a swarm. It provides
a clean interface to connect to a service on a given remote peer.

First, you'll need an ID, and a place to store that ID. To generate a 'test'
ID, you can do the following:
```go
import (
	testutil "github.com/libp2p/go-testutil"
	pstore "github.com/libp2p/go-libp2p-peerstore"
)

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
```

Next, you'll need at least one address that you want to listen on. You can go
from a string to a multiaddr like this:
```go
import ma "github.com/multiformats/go-multiaddr"

...

maddr, err := ma.NewMultiaddr("/ip4/0.0.0.0/tcp/9000")
if err != nil {
	panic(err)
}
```

Now you know who you are, and where you live (in a manner of speaking). The
next step is setting up a 'swarm network' to handle all the peers you will
connect to. The swarm handles incoming connections from other peers, and
handles the negotiation of new outbound connections.
```go
import (
	"context"
	swarm "github.com/libp2p/go-libp2p-swarm"
)

// Make a context to govern the lifespan of the swarm
ctx := context.Background()

// Put all this together
netw, err := swarm.NewNetwork(ctx, []ma.Multiaddr{maddr}, ident.ID(), ps, nil)
if err != nil {
	panic(err)
}
```

At this point, we have everything needed to finally construct a host. That call
is the simplest one so far:
```go
import bhost "github.com/libp2p/go-libp2p/p2p/host/basic"

myhost := bhost.New(netw)
```

And thats it, you have a libp2p host and you're ready to start doing some
awesome p2p networking!

In future guides we will go over ways to use hosts, configure them differently
(hint: there are a huge number of ways to set these up), and interesting ways
to apply this technology to various applications you might want to build.

To see this code all put together, take a look at the `host.go` file in this
directory.
