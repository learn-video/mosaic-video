// Package watcher provides a watcher that can be used to watch for changes
// in the local filesystem.
package watcher

type Watcher interface {
	Start() error
	Stop()
}
