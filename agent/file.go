package agent

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"ditto.co.jp/agentserver/cx"
	"ditto.co.jp/agentserver/logger"
	"github.com/labstack/echo/v4"
)

//Upload - upload file
func (s *Service) Upload(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// post url/job/id?num=1
	// id := c.Param("id")
	num := c.QueryParam("num")
	dir := c.QueryParam("dir")
	// fname := fmt.Sprintf("%v-%v.part", id, num)
	fname := file.Filename
	n, _ := strconv.ParseInt(num, 10, 32)
	if n > 0 {
		fname = fmt.Sprintf("%v-%v", file.Filename, num)
	}
	if err := s.MkDirIfNeeded(dir); err != nil {
		logger.Warn(err)
	}
	fname = filepath.Join(dir, fname)

	dst, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	result := cx.Result{
		StatusCode: 200,
	}

	return c.JSON(http.StatusOK, result)
}

//MkDirIfNeeded -
func (s *Service) MkDirIfNeeded(dir string) error {
	if lastIdx := strings.LastIndex(dir, "/"); lastIdx != -1 {
		dirPath := dir[:lastIdx]
		if _, ok := s.dirsSeen[dirPath]; !ok {
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return err
			}
			s.dirsSeen[dirPath] = true
		}
	}
	return nil
}
