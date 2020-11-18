package master

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"ditto.co.jp/agentserver/cx"
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
