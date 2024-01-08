package watcher

import "go.uber.org/fx"

var Module = fx.Provide(NewFileSystemWatcher)
