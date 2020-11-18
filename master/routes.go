package master

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

//RegisterRoutes -
func (s *Service) RegisterRoutes(e *echo.Echo, prefix string) {
	g := e.Group(prefix)
	s.prefix = prefix

	g.GET("/health", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "master.alive")
	})

	e.GET("/agent/alloc", s.GetAgentIdle)
	//全Agent取得
	g.GET("/agent/list", s.ListAgent)
	//指定されたAgent取得
	g.GET("/agent/:id", s.GetAgent)
	//新規追加
	g.POST("/agent", s.AddAgent)
	//削除
	g.DELETE("/agent/:id", s.DeleteAgent)
	//アップロード
	g.POST("/:id", s.Upload)
}
