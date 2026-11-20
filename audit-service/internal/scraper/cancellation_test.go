package scraper

import (
	"context"
	"errors"
	"net/netip"
	"testing"
)

func TestCancelledCallerDoesNotWaitForBrowserSlot(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	scraper := &Scraper{semaphore: make(chan struct{})}
	_, err := scraper.Scrape(ctx, "https://192.0.2.1/")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
}

func TestPublicAddressPolicy(t *testing.T) {
	for _, addr := range []string{"127.0.0.1", "::1", "10.0.0.1", "169.254.1.1", "100.64.1.1", "::ffff:127.0.0.1"} {
		if publicAddress(netip.MustParseAddr(addr)) {
			t.Errorf("unexpected public classification for %s", addr)
		}
	}
	if !publicAddress(netip.MustParseAddr("1.1.1.1")) {
		t.Error("public DNS address rejected")
	}
}
