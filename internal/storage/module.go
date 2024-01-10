package storage

import (
	"github.com/mauricioabreu/mosaic-video/internal/storage/s3"
	"go.uber.org/fx"
)

var Module = fx.Provide(s3.NewClient)
