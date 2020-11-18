package cmd

import (
	"errors"
	"fmt"

	"ditto.co.jp/agentserver/config"
	"ditto.co.jp/agentserver/cx"
	"ditto.co.jp/agentserver/db"
	"ditto.co.jp/agentserver/logger"
	"ditto.co.jp/agentserver/s3"
	"ditto.co.jp/agentserver/util"
)

//RunWorker -
func RunWorker(conf *config.Config, db *db.Database) error {
	// quit := false
	// s := single.New("s3worker")
	// err := s.CheckLock()
	// if err != nil && err == single.ErrAlreadyRunning {
	// 	logger.Warn(err)
	// 	quit = true
	// }
	// defer s.TryUnlock()
	// if quit {
	// 	return err
	// }
	// worker
	// host:port -> agent's address

	host := conf.Host
	logger.Info(host)
	if len(host) == 0 || host == "0.0.0.0" {
		host = conf.LocalHost
	}
	port := conf.Port
	if port == 0 {
		//agent port
		port = 9091
	}
	//
	var err error
	var files []cx.File
	addr := fmt.Sprintf("%v:%v", host, port)
	if conf.JobID == "all" {
		files, err = util.ListJob(addr)
	} else {
		files, err = util.GetJob(addr, conf.JobID)
	}
	if err != nil {
		return err
	}

	service := s3.NewService(conf)
	if service == nil {
		return errors.New("Failed create s3 service")
	}

	service.Downlad(files)

	// //upload
	// upload := fmt.Sprintf("http://%v/job/%v?num=%v", host, conf.JobID, 1)
	// pr, pw := io.Pipe()
	// mw := multipart.NewWriter(pw)
	// go func() {
	// 	defer pw.Close()

	// 	fw, err := mw.CreateFormFile("file", "filename")
	// 	if err != nil {

	// 	}
	// 	_, err = io.Copy(fw, file)
	// 	if err != nil {

	// 	}
	// 	if err := mw.Close(); err != nil {

	// 	}
	// }()
	// contentType := mw.FormDataContentType()
	// res, err := http.Post(upload, contentType, pr)
	// if err != nil {
	// 	return err
	// }
	// defer res.Body.Close()
	// message, _ := ioutil.ReadAll(res.Body)
	// fmt.Printf(string(message))

	return nil
}
