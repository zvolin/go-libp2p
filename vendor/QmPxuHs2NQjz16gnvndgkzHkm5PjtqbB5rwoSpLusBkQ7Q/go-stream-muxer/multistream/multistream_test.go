package multistream

import (
	"testing"

	test "QmPxuHs2NQjz16gnvndgkzHkm5PjtqbB5rwoSpLusBkQ7Q/go-stream-muxer/test"
)

func TestMultiStreamTransport(t *testing.T) {
	test.SubtestAll(t, NewTransport())
}
