package agent

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

//RegisterRoutes -
func (s *Service) RegisterRoutes(e *echo.Echo, prefix string) {
	g := e.Group(prefix)
	s.prefix = prefix

	g.GET("/health", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "agent.alive")
	})

	g.POST("/:id", s.Upload)
	//Update upload server
	g.POST("/upsrv", s.SetUploadServer)
}
