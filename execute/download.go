package execute

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"ditto.co.jp/agentserver/config"
	"ditto.co.jp/agentserver/cx"
	"ditto.co.jp/agentserver/logger"
	"ditto.co.jp/agentserver/s3"
	"ditto.co.jp/agentserver/util"
	"github.com/dustin/go-humanize"
)

//
var (
	ErrS3Service = errors.New("Failed create s3 service")
)

//Download - download file from s3
func Download(conf *config.Config) error {
	defer util.TimeTrack(time.Now(), "s3transfer")

	s3svc := s3.NewService(conf)
	if s3svc == nil {
		return ErrS3Service
	}

	err := s3svc.SetURI(conf.S3)
	if err != nil {
		return err
	}

	t := time.Now()
	var downloaded int64 = 0
	var downloadedfile int64 = 0

	logger.Tracef("Scan s3 in list size %v", conf.ListSize)
	keysChan := make(chan *s3.File, conf.ListSize)
	wg := new(sync.WaitGroup)

	var i uint = 0
	for i = 0; i < conf.Threads; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			// files := make([]*cx.File, 0)
			for fi := range keysChan {
				atomic.AddInt64(&downloaded, fi.Size)
				atomic.AddInt64(&downloadedfile, 1)

				// logger.Println(fi.Key)
				cxf := &cx.File{
					S3:     fi.S3,
					Local:  fi.SaveTo,
					Size:   uint64(fi.Size),
					Num:    0,
					Offset: 0,
					Length: 0,
				}
				logger.Tracef("[%v] %v", downloadedfile, cxf.S3)
				_ = cxf
				//				files = append(files, cxf)
			}
			// jsontxt := cx.ToJSON(files)
			// logger.Println(jsontxt)
		}()
	}

	err = s3svc.List(keysChan)
	if err != nil {
		logger.Error(err)
	}

	close(keysChan)
	wg.Wait()

	logger.Info("Finished.")
	dt := time.Now().Sub(t)
	logger.Infof("Down time: %v, TotalSize: %+v (%v), FileCount: %+v\n", dt, humanize.Bytes(uint64(downloaded)), humanize.Comma(downloaded), downloadedfile)

	return nil
}
