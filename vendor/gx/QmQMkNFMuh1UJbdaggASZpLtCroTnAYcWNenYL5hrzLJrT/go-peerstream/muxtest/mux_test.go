package muxtest

import (
	"testing"

	multiplex "gx/QmRmT6MSnfhRDW1PTUGSd3z4fqXK48GUequQAZzeT4c5iC/go-stream-muxer/multiplex"
	multistream "gx/QmRmT6MSnfhRDW1PTUGSd3z4fqXK48GUequQAZzeT4c5iC/go-stream-muxer/multistream"
	muxado "gx/QmRmT6MSnfhRDW1PTUGSd3z4fqXK48GUequQAZzeT4c5iC/go-stream-muxer/muxado"
	spdy "gx/QmRmT6MSnfhRDW1PTUGSd3z4fqXK48GUequQAZzeT4c5iC/go-stream-muxer/spdystream"
	yamux "gx/QmRmT6MSnfhRDW1PTUGSd3z4fqXK48GUequQAZzeT4c5iC/go-stream-muxer/yamux"
)

func TestYamuxTransport(t *testing.T) {
	SubtestAll(t, yamux.DefaultTransport)
}

func TestSpdyStreamTransport(t *testing.T) {
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
