package watcher

import (
	"github.com/fsnotify/fsnotify"
	"github.com/mauricioabreu/mosaic-video/internal/config"
)

type FileSystemWatcher struct {
	watcher *fsnotify.Watcher
	events  chan fsnotify.Event
	errors  chan error
}

func NewFileSystemWatcher(cfg config.Config) (*FileSystemWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := watcher.Add(cfg.AssetsPath); err != nil {
		return nil, err
	}

	return &FileSystemWatcher{
		watcher: watcher,
		events:  make(chan fsnotify.Event),
		errors:  make(chan error),
	}, nil
}

func (fsw *FileSystemWatcher) Run() {
	go func() {
		for {
			select {
			case event := <-fsw.watcher.Events:
				fsw.events <- event
			case err := <-fsw.watcher.Errors:
				fsw.errors <- err
			}
		}
	}()
}

func (fsw *FileSystemWatcher) Events() <-chan fsnotify.Event {
	return fsw.events
}

func (fsw *FileSystemWatcher) Errors() <-chan error {
	return fsw.errors
}
