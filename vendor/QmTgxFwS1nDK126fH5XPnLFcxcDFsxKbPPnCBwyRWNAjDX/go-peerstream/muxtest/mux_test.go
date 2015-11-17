package muxtest

import (
	"testing"

	multiplex "QmPxuHs2NQjz16gnvndgkzHkm5PjtqbB5rwoSpLusBkQ7Q/go-stream-muxer/multiplex"
	multistream "QmPxuHs2NQjz16gnvndgkzHkm5PjtqbB5rwoSpLusBkQ7Q/go-stream-muxer/multistream"
	muxado "QmPxuHs2NQjz16gnvndgkzHkm5PjtqbB5rwoSpLusBkQ7Q/go-stream-muxer/muxado"
	spdy "QmPxuHs2NQjz16gnvndgkzHkm5PjtqbB5rwoSpLusBkQ7Q/go-stream-muxer/spdystream"
	yamux "QmPxuHs2NQjz16gnvndgkzHkm5PjtqbB5rwoSpLusBkQ7Q/go-stream-muxer/yamux"
)

func TestYamuxTransport(t *testing.T) {
	SubtestAll(t, yamux.DefaultTransport)
}

func TestSpdyStreamTransport(t *testing.T) {
	t.SkipNow()
	SubtestAll(t, spdy.Transport)
}

func TestMultiplexTransport(t *testing.T) {
	SubtestAll(t, multiplex.DefaultTransport)
}

func TestMuxadoTransport(t *testing.T) {
	SubtestAll(t, muxado.Transport)
}

func TestMultistreamTransport(t *testing.T) {
	SubtestAll(t, multistream.NewTransport())
}
