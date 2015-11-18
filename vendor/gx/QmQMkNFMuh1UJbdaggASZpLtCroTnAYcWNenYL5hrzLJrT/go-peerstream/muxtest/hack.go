package muxtest

import (
	multiplex "gx/QmRmT6MSnfhRDW1PTUGSd3z4fqXK48GUequQAZzeT4c5iC/go-stream-muxer/multiplex"
	multistream "gx/QmRmT6MSnfhRDW1PTUGSd3z4fqXK48GUequQAZzeT4c5iC/go-stream-muxer/multistream"
	muxado "gx/QmRmT6MSnfhRDW1PTUGSd3z4fqXK48GUequQAZzeT4c5iC/go-stream-muxer/muxado"
	spdy "gx/QmRmT6MSnfhRDW1PTUGSd3z4fqXK48GUequQAZzeT4c5iC/go-stream-muxer/spdystream"
	yamux "gx/QmRmT6MSnfhRDW1PTUGSd3z4fqXK48GUequQAZzeT4c5iC/go-stream-muxer/yamux"
)

var _ = multiplex.DefaultTransport
var _ = multistream.NewTransport
var _ = muxado.Transport
var _ = spdy.Transport
var _ = yamux.DefaultTransport
