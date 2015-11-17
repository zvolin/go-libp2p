package muxtest

import (
	multiplex "QmPxuHs2NQjz16gnvndgkzHkm5PjtqbB5rwoSpLusBkQ7Q/go-stream-muxer/multiplex"
	multistream "QmPxuHs2NQjz16gnvndgkzHkm5PjtqbB5rwoSpLusBkQ7Q/go-stream-muxer/multistream"
	muxado "QmPxuHs2NQjz16gnvndgkzHkm5PjtqbB5rwoSpLusBkQ7Q/go-stream-muxer/muxado"
	spdy "QmPxuHs2NQjz16gnvndgkzHkm5PjtqbB5rwoSpLusBkQ7Q/go-stream-muxer/spdystream"
	yamux "QmPxuHs2NQjz16gnvndgkzHkm5PjtqbB5rwoSpLusBkQ7Q/go-stream-muxer/yamux"
)

var _ = multiplex.DefaultTransport
var _ = multistream.NewTransport
var _ = muxado.Transport
var _ = spdy.Transport
var _ = yamux.DefaultTransport
