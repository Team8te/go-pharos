package store

import (
	"net/http"

	"github.com/go-pharos/app/ds"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type MoveFileRequest struct {
	Files  []string  `json:"files"`
	Target ds.Target `json:"target"`
}

// DeviceList godoc
// @Description get the status of server.
// @Tags Store
// @Accept json
// @Produce json
// @Param account body store.MoveFileRequest true "Request"
// @Success 200 {string} string	"ok"
// @Router /store/file/move [post]
func (e *StoreEndpoint) MoveFile(ctx echo.Context) error {
	log.WithField("api", "http").Debugf("StoreList")

	req := &MoveFileRequest{}
	err := ctx.Bind(req)
	if err != nil {
		return err
	}

	e.st.MoveFiles(ctx.Request().Context(), req.Files, &req.Target)
	ctx.JSON(http.StatusOK, nil)
	return nil
}
