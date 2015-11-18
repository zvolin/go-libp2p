package multistream

import (
	"testing"

	test "gx/QmRmT6MSnfhRDW1PTUGSd3z4fqXK48GUequQAZzeT4c5iC/go-stream-muxer/test"
)

func TestMultiStreamTransport(t *testing.T) {
	test.SubtestAll(t, NewTransport())
}
