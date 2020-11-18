package job

import (
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"

	"ditto.co.jp/agentserver/cx"
	"ditto.co.jp/agentserver/logger"
	"ditto.co.jp/agentserver/util"
	"github.com/kardianos/osext"
	"github.com/labstack/echo/v4"
)

//ExecJobs -
func (s *Service) ExecJobs(c echo.Context) error {
	storage := c.QueryParam("storage")
	if len(storage) == 0 {
		storage = fmt.Sprintf("%v:%v", s._config.LocalHost, s._config.Port)
	}
	dir := c.QueryParam("dir")
	dir = filepath.ToSlash(dir)
	logger.Debugf("dir: %v", dir)

	err := s.execute("all", storage, dir)
	if err != nil {
		return err
	}

	result := cx.Result{
		StatusCode: 200,
	}

	return c.JSON(http.StatusOK, result)
}

//FetchJob -
func (s *Service) FetchJob(c echo.Context) error {
	val, err := s._db.Fetch("jobs")
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, val)
}

//ListJob -
func (s *Service) ListJob(c echo.Context) error {
	val, err := s._db.List("jobs")
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, val)
}

//GetJob -
func (s *Service) GetJob(c echo.Context) error {
	id := c.Param("id")
	jobid := fmt.Sprintf("job:%v", id)

	val, err := s._db.Get(jobid)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, val)
}

//SubmitJob -
func (s *Service) SubmitJob(c echo.Context) error {
	// var v []File = make([]File, 0)
	// err := json.NewDecoder(c.Request().Body).Decode(&v)
	// if err != nil {
	// 	return err
	// }
	// id := c.Param("id")
	// fmt.Println("ID: ", id)
	// for _, f := range v {
	// 	fmt.Println(f.S3)
	// }
	body := util.StreamToString(c.Request().Body)

	id := c.Param("id")
	jobid := fmt.Sprintf("job:%v", id)

	err := s._db.Set(jobid, body)
	if err != nil {
		return err
	}
	// storage := c.QueryParam("storage")
	// if len(storage) == 0 {
	// 	storage = fmt.Sprintf("%v:%v", s._config.LocalHost, s._config.Port)
	// }
	// dir := c.QueryParam("dir")
	// err = s.execute(id, storage, dir)
	// if err != nil {
	// 	return err
	// }

	result := cx.Result{
		StatusCode: 200,
	}

	return c.JSON(http.StatusOK, result)
}

func (s *Service) execute(jobid string, storage string, dir string) error {
	port := fmt.Sprintf("%v", s._config.Port)
	exename, _ := osext.Executable()
	// Command set
	command := exec.Command(
		exename,
		"worker",
		"--jobid",
		jobid,
		"--host",
		s._config.Host,
		"--port",
		port,
		"--storage",
		storage,
		"--dir",
		dir,
	)
	command.Dir = filepath.Dir(exename)
	logger.Debug(command)

	err := command.Start()
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
