package store

import (
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// DeviceList godoc
// @Description get the status of server.
// @Tags Store
// @Accept json
// @Produce json
// @Param path query string	false	"path"
// @Success 200 {array} ds.File "ok"
// @Router /store/file/list [get]
func (e *StoreEndpoint) StoreList(ctx echo.Context) error {
	log.WithField("api", "http").Debugf("StoreList")

	files, err := e.st.ListDir(ctx.Request().Context(), ctx.QueryParam("path"))
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, files)

	return nil
}
