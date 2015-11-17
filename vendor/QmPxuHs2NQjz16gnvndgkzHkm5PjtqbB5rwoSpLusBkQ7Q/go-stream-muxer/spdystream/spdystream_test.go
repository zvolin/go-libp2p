package peerstream_spdystream

import (
	"testing"

	test "QmPxuHs2NQjz16gnvndgkzHkm5PjtqbB5rwoSpLusBkQ7Q/go-stream-muxer/test"
)

func TestSpdyStreamTransport(t *testing.T) {
	test.SubtestAll(t, Transport)
}
