package agent

import (
	"fmt"
	"net/http"

	"ditto.co.jp/agentserver/cx"
	"github.com/labstack/echo/v4"
)

//SetUploadServer -
func (s *Service) SetUploadServer(c echo.Context) error {
	// /upsrv?ip=&port=
	ip := c.QueryParam("ip")
	port := c.QueryParam("port")

	addr := fmt.Sprintf("%v:%v", ip, port)
	s._config.Storage = addr

	result := cx.Result{
		StatusCode: 200,
		Message:    addr,
	}

	return c.JSON(http.StatusOK, result)
}
