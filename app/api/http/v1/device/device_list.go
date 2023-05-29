package device

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// DeviceList godoc
// @Description get the status of server.
// @Tags Device
// @Accept json
// @Produce json
// @Param limit query int	false	"limit"
// @Param offset query int	false	"offset"
// @Success 200 {array} ds.Device "ok"
// @Router /device/list [get]
func (e *DeviceEndpoint) DeviceList(ctx echo.Context) error {
	log.WithField("api", "http").Debugf("DeviceList")

	limit, _ := strconv.Atoi(getQueryOrDefault(ctx, "limit", "10"))
	offset, _ := strconv.Atoi(ctx.QueryParam("limit"))

	devices, err := e.d.DeviceList(ctx.Request().Context(), limit, offset)
	if err != nil {
		return err
	}
	ctx.JSON(http.StatusOK, devices)
	return nil
}

func getQueryOrDefault(ctx echo.Context, name, d string) string {
	if v := ctx.QueryParam("limit"); v != "" {
		return v
	}

	return d
}
