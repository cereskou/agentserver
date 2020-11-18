package master

import (
	"github.com/labstack/echo/v4"
)

//ServiceInterface -
type ServiceInterface interface {
	RegisterRoutes(e *echo.Echo, prefix string)
	Close()
}
