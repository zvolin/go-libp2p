package config

import (
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/sec"
	"github.com/libp2p/go-libp2p/core/sec/insecure"
	csms "github.com/libp2p/go-libp2p/p2p/net/conn-security-multistream"
)

func makeInsecureTransport(id peer.ID, privKey crypto.PrivKey) sec.SecureMuxer {
	secMuxer := new(csms.SSMuxer)
	secMuxer.AddTransport(insecure.ID, insecure.NewWithIdentity(insecure.ID, id, privKey))
	return secMuxer
}

func makeSecurityMuxer(tpts []sec.SecureTransport) sec.SecureMuxer {
	secMuxer := new(csms.SSMuxer)
	for _, tpt := range tpts {
		secMuxer.AddTransport(string(tpt.ID()), tpt)
	}
	return secMuxer
}
