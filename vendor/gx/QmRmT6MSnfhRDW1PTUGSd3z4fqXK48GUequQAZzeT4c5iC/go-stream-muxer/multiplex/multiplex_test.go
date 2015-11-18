package peerstream_multiplex

import (
	"testing"

	test "gx/QmRmT6MSnfhRDW1PTUGSd3z4fqXK48GUequQAZzeT4c5iC/go-stream-muxer/test"
)

func TestMultiplexTransport(t *testing.T) {
	test.SubtestAll(t, DefaultTransport)
}
