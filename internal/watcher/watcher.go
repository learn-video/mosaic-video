// Package watcher provides a file system watcher, useful to collect events
// from the assets directory while also being able to catch errors.
package watcher

import (
	"github.com/fsnotify/fsnotify"
)

type Watcher interface {
	Run()
	Events() <-chan fsnotify.Event
	Errors() <-chan error
}
