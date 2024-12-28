package store

import (
	"mime/multipart"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type LoadRequest struct {
	Path string                `form:"path"`
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
func (e *StoreEndpoint) Load(ctx echo.Context) (err error) {
	log.WithField("api", "http").Debugf("Load")
	req := &LoadRequest{}
	ctx.Bind(req)
	req.File, err = ctx.FormFile("file")
	if err != nil {
		return err
	}

	src, err := req.File.Open()
	if err != nil {
		return err
	}

	defer src.Close()
	e.st.UploadFile(ctx.Request().Context(), src, req.Path)

	return nil
}
