package job

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

//RegisterRoutes -
func (s *Service) RegisterRoutes(e *echo.Echo, prefix string) {
	g := e.Group(prefix)
	s.prefix = prefix

	g.GET("/health", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "alive")
	})

	//POST JOB
	g.PUT("/:id", s.SubmitJob)
	g.GET("/:id", s.GetJob)
	g.GET("/list", s.ListJob)
	g.GET("/fetch", s.FetchJob)
	g.POST("/exec", s.ExecJobs)
}
