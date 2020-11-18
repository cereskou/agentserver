package master

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"ditto.co.jp/agentserver/cx"
	"ditto.co.jp/agentserver/logger"
	"ditto.co.jp/agentserver/util"
	"github.com/labstack/echo/v4"

	jsoniter "github.com/json-iterator/go"
)

// -
var (
	ErrAgentNotFound = errors.New("not found")
)

//ListAgent -
func (s *Service) ListAgent(c echo.Context) error {
	data := s._db.Map.List()

	return c.JSON(http.StatusOK, data)
}

//GetAgentIdle -
func (s *Service) GetAgentIdle(c echo.Context) error {
	result := cx.Result{
		StatusCode: 404,
	}
	id := s._db.Map.MinKey()
	agent := s._db.Map.Get(id)
	if agent != nil {
		resp := &util.EventMessage{
			Mode:      agent.Mode,
			IP:        agent.IP,
			Port:      agent.Port,
			TimeStamp: agent.TimeStamp,
			ID:        s._db.Map.GenerateID(),
		}
		agent.Status++
		agent.TimeStamp = time.Now()

		return c.JSON(http.StatusOK, resp)
	}

	return c.JSON(http.StatusNotFound, result)
}

//GetAgent -
func (s *Service) GetAgent(c echo.Context) error {
	id := c.Param("id")
	data := s._db.Map.Get(id)
	if data == nil {
		result := cx.Result{
			StatusCode: 404,
			Error:      ErrAgentNotFound,
		}

		return c.JSON(http.StatusNotFound, result)
	}

	return c.JSON(http.StatusOK, data)
}

//AddAgent -
func (s *Service) AddAgent(c echo.Context) error {
	result := cx.Result{
		StatusCode: 204,
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var v util.EventMessage
	err := json.NewDecoder(c.Request().Body).Decode(&v)
	if err != nil {
		logger.Error(err)
		return err
	}
	//ID
	id := fmt.Sprintf("%v:%v", v.IP, v.Port)
	logger.Tracef("Agent: %v", id)

	//Not Exists
	if !s._db.Map.Exists(id) {
		s._db.Map.Set(id, &v)
	}

	return c.JSON(http.StatusOK, result)
}

//DeleteAgent -
func (s *Service) DeleteAgent(c echo.Context) error {
	result := cx.Result{
		StatusCode: 200,
	}
	id := c.Param("id")
	logger.Tracef("delete %v", id)

	s._db.Map.Delete(id)

	return c.JSON(http.StatusOK, result)
}
