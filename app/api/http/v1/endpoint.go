package v1

import (
	"github.com/labstack/echo/v4"
)

type HTTPEndpoint struct {
	handlerMap map[string]func(c echo.Context) error
}

func NewHTTPEndpoint() *HTTPEndpoint {
	return &HTTPEndpoint{}
}
