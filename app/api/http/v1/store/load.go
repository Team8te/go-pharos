package store

import (
	"io"
	"mime/multipart"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type LoadRequest struct {
	Name string                `json:"name"`
	File *multipart.FileHeader `form:"file" binding:"required"`
}

// DeviceList godoc
// @Description get the status of server.
// @Tags Store
// @Accept mpfd
// @Produce json
// @Param name formData string true "file_name"
// @Param file formData file true "file"
// @Success 200 {string} string	"ok"
// @Router /store/file/load [post]
func (e *StoreEndpoint) Load(ctx echo.Context) error {
	log.WithField("api", "http").Debugf("Load")
	req := &LoadRequest{}
	ctx.Bind(req)
	file, err := ctx.FormFile("file")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}

	defer src.Close()
	io.Copy()

	return nil
}
