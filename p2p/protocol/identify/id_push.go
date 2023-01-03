package identify

import (
	"github.com/libp2p/go-libp2p/core/network"
)

// IDPush is the protocol.ID of the Identify push protocol.
// It sends full identify messages containing the current state of the peer.
const IDPush = "/ipfs/id/push/1.0.0"

// pushHandler handles incoming identify push streams. The behaviour is identical to the ordinary identify protocol.
func (ids *idService) pushHandler(s network.Stream) {
	ids.handleIdentifyResponse(s, true)
}
