package sm_yamux

import (
	"testing"

	test "QmPxuHs2NQjz16gnvndgkzHkm5PjtqbB5rwoSpLusBkQ7Q/go-stream-muxer/test"
)

func TestYamuxTransport(t *testing.T) {
	test.SubtestAll(t, DefaultTransport)
}
