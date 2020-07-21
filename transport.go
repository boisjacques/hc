package hc

import "sync"

// Transport provides accessories over a network.
type Transport interface {
	// Start starts the transport
	Start(wg *sync.WaitGroup)

	// Stop stops the transport
	// Use the returned channel to wait until the transport is fully stopped.
	Stop() <-chan struct{}
}
