package worker

import "github.com/mauricioabreu/mosaic-video/internal/mosaic"

const (
	TypeStartMosaic = "mosaic:start"
	TypeStopMosaic  = "mosaic:stop"
)

type StartMosaicPayload struct {
	Mosaic mosaic.Mosaic
}
