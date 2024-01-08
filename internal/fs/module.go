package fs

import "go.uber.org/fx"

var Module = fx.Provide(NewFileSystemWatcher)
