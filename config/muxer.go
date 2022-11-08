package config

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	msmux "github.com/libp2p/go-libp2p/p2p/muxer/muxer-multistream"
)

type Muxer struct {
	ID          protocol.ID
	Multiplexer network.Multiplexer
}

func makeMuxer(muxers []Muxer) (network.Multiplexer, error) {
	muxMuxer := msmux.NewBlankTransport()
	transportSet := make(map[protocol.ID]struct{}, len(muxers))
	for _, m := range muxers {
		if _, ok := transportSet[m.ID]; ok {
			return nil, fmt.Errorf("duplicate muxer transport: %s", m.ID)
		}
		transportSet[m.ID] = struct{}{}
	}
	for _, m := range muxers {
		muxMuxer.AddTransport(string(m.ID), m.Multiplexer)
	}
	return muxMuxer, nil
}
