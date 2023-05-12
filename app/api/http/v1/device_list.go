package v1

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// DeviceList godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /device/list [get]
func (e *HTTPEndpoint) DeviceList(ctx echo.Context) error {
	last_id := ctx.FormValue("last_id")
	log.WithField("api", "http").Debugf("DeviceList. Last_id %v", last_id)

	return nil
}
