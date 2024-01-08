package watcher

import (
	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/rjeczalik/notify"
)

type FileSystemWatcher struct {
	events chan notify.EventInfo
	path   string
}

func NewFileSystemWatcher(cfg *config.Config) (*FileSystemWatcher, error) {
	c := make(chan notify.EventInfo, 1)

	return &FileSystemWatcher{
		events: c,
		path:   cfg.AssetsPath,
	}, nil
}

func (fsw *FileSystemWatcher) Start() error {
	return notify.Watch(fsw.path+"/...", fsw.events, notify.Write)
}

func (fsw *FileSystemWatcher) Stop() {
	notify.Stop(fsw.events)
}

func (fsw *FileSystemWatcher) Events() chan notify.EventInfo {
	return fsw.events
}
