/*
Copyright © 2026
*/
// Package server provides the HTTP server interface.
//
package server

import "context"

// Server defines the lifecycle interface for the HTTP server.
// cmd/run.go uses this interface to start and stop the server
// without depending on any specific web framework.
type Server interface {
	// Run starts the HTTP server in a goroutine and sends any fatal
	// startup error to errChan. It returns immediately (non-blocking).
	Run(errChan chan<- error)

	// Stop performs a graceful shutdown, respecting the context deadline.
	Stop(ctx context.Context) error
}
