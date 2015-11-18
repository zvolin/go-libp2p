package peerstream_spdystream

import (
	"testing"

	test "gx/QmRmT6MSnfhRDW1PTUGSd3z4fqXK48GUequQAZzeT4c5iC/go-stream-muxer/test"
)

func TestSpdyStreamTransport(t *testing.T) {
	test.SubtestAll(t, Transport)
}
