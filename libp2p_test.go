package libp2p

import (
	"context"
	"fmt"
	"strings"
	"testing"

	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-tcp-transport"
)

func TestNewHost(t *testing.T) {
	h, err := makeRandomHost(t, 9000)
	if err != nil {
		t.Fatal(err)
	}
	h.Close()
}

func TestBadTransportConstructor(t *testing.T) {
	ctx := context.Background()
	h, err := New(ctx, Transport(func() {}))
	if err == nil {
		h.Close()
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "libp2p_test.go") {
		t.Error("expected error to contain debugging info")
	}
}

func TestInsecure(t *testing.T) {
	ctx := context.Background()
	h, err := New(ctx, NoSecurity)
	if err != nil {
		t.Fatal(err)
	}
	h.Close()
}

func TestDefaultListenAddrs(t *testing.T) {
	ctx := context.Background()

	// Test 1: Listen addr should not set if user defined transport is passed.
	h, err := New(
		ctx,
		Transport(tcp.TcpTransport{}),
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(h.Addrs()) != 0 {
		t.Error("expected zero listen addrs as none is set with user defined transport")
	}
	h.Close()

	// Test 2: User defined listener addrs should overwrite the default options.
	h, err = New(
		ctx,
		Transport(tcp.TcpTransport{}),
		ListenAddrStrings("/ip4/127.0.0.1/tcp/0"),
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(h.Addrs()) != 1 {
		t.Error("expected one listen addr")
	}
	h.Close()
}

func makeRandomHost(t *testing.T, port int) (host.Host, error) {
	ctx := context.Background()
	priv, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		t.Fatal(err)
	}

	opts := []Option{
		ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port)),
		Identity(priv),
		DefaultTransports,
		DefaultMuxers,
		DefaultSecurity,
		NATPortMap(),
	}

	return New(ctx, opts...)
}
