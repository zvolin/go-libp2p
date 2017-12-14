package libp2p

import (
	"context"
	"crypto/rand"

	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	pnet "github.com/libp2p/go-libp2p-interface-pnet"
	metrics "github.com/libp2p/go-libp2p-metrics"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	swarm "github.com/libp2p/go-libp2p-swarm"
	transport "github.com/libp2p/go-libp2p-transport"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	mux "github.com/libp2p/go-stream-muxer"
	ma "github.com/multiformats/go-multiaddr"
	mplex "github.com/whyrusleeping/go-smux-multiplex"
	msmux "github.com/whyrusleeping/go-smux-multistream"
	yamux "github.com/whyrusleeping/go-smux-yamux"
)

// Config describes a set of settings for a libp2p node
type Config struct {
	Transports   []transport.Transport
	Muxer        mux.Transport
	ListenAddrs  []ma.Multiaddr
	PeerKey      crypto.PrivKey
	Peerstore    pstore.Peerstore
	Protector    pnet.Protector
	Reporter     metrics.Reporter
	DisableSecio bool
}

func New(ctx context.Context) (host.Host, error) {
	return NewWithCfg(ctx, DefaultConfig())
}

// Construct instantiates a libp2p host using information from the given
// config. `nil` may be passed to use default options.
func NewWithCfg(ctx context.Context, cfg *Config) (host.Host, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// If no key was given, generate a random 2048 bit RSA key
	if cfg.PeerKey == nil {
		priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
		if err != nil {
			return nil, err
		}
		cfg.PeerKey = priv
	}

	// Obtain Peer ID from public key
	pid, err := peer.IDFromPublicKey(cfg.PeerKey.GetPublic())
	if err != nil {
		return nil, err
	}

	// Create a new blank peerstore if none was passed in
	ps := cfg.Peerstore
	if ps == nil {
		ps = pstore.NewPeerstore()
	}

	// If secio is disabled, don't add our private key to the peerstore
	if !cfg.DisableSecio {
		ps.AddPrivKey(pid, cfg.PeerKey)
		ps.AddPubKey(pid, cfg.PeerKey.GetPublic())
	}

	swrm, err := swarm.NewSwarmWithProtector(ctx, cfg.ListenAddrs, pid, ps, cfg.Protector, cfg.Muxer, cfg.Reporter)
	if err != nil {
		return nil, err
	}

	netw := (*swarm.Network)(swrm)

	return bhost.New(netw), nil
}

func DefaultMuxer() mux.Transport {
	// Set up stream multiplexer
	tpt := msmux.NewBlankTransport()

	// By default, support yamux and multiplex
	tpt.AddTransport("/yamux/1.0.0", yamux.DefaultTransport)
	tpt.AddTransport("/mplex/6.3.0", mplex.DefaultTransport)

	return tpt
}

func DefaultConfig() *Config {
	// Create a multiaddress that listens on a random port on all interfaces
	addr, err := ma.NewMultiaddr("/ip4/0.0.0.0/tcp/0")
	if err != nil {
		panic(err)
	}

	return &Config{
		ListenAddrs: []ma.Multiaddr{addr},
		Peerstore:   pstore.NewPeerstore(),
		Muxer:       DefaultMuxer(),
	}
}
