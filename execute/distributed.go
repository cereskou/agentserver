package execute

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"ditto.co.jp/agentserver/config"
	"ditto.co.jp/agentserver/cx"
	"ditto.co.jp/agentserver/logger"
	"ditto.co.jp/agentserver/s3"
	"ditto.co.jp/agentserver/util"
	"github.com/dustin/go-humanize"
	"gopkg.in/resty.v1"
)

//
var (
	ErrMasterServerNotFound = errors.New("Master server not found")
)

//Distributed -
func Distributed(conf *config.Config) error {
	defer util.TimeTrack(time.Now(), "s3transfer")

	//Find master
	logger.Info("Scan master node...")
	logger.Info(conf.LocalHost)
	host := conf.Host
	if len(host) == 0 || host == "0.0.0.0" {
		host = conf.LocalHost
	}
	port := conf.Port
	if port == 0 {
		port = 9090
	}
	//
	storage := conf.Storage
	if len(storage) == 0 {
		storage = fmt.Sprintf("%v:%v", conf.LocalHost, conf.Port)
	}

	dir := conf.Dir
	addr := fmt.Sprintf("%v:%v", host, port)
	logger.Tracef("master note : %v", addr)
	logger.Tracef("storage node: %v", storage)
	logger.Tracef("save to: %v", dir)

	alive := util.ServerCheck(addr)
	if len(alive) == 0 || alive != "master.alive" {
		return ErrMasterServerNotFound
	}

	logger.Debug("Scan started")
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
	var download2 int64 = 0
	var splited int64 = 0

	keysChan := make(chan *s3.File, conf.DistListSize)
	wg := new(sync.WaitGroup)

	var i uint = 0
	for i = 0; i < conf.DistThreads; i++ {
		wg.Add(1)

		go func(chid uint) {
			defer wg.Done()

			// JSON
			buf := make([]*cx.File, 0)
			for fi := range keysChan {
				atomic.AddInt64(&downloaded, fi.Size)
				atomic.AddInt64(&downloadedfile, 1)

				cxf := &cx.File{
					S3:     fi.S3,
					Local:  fi.SaveTo,
					Size:   uint64(fi.Size),
					Num:    fi.Num,
					Offset: fi.Offset,
					Length: fi.Length,
				}
				if cxf.Num > 0 {
					atomic.AddInt64(&splited, 1)
				}
				buf = append(buf, cxf)
				if len(buf) >= conf.DistJobSize {
					atomic.AddInt64(&download2, int64(len(buf)))
					logger.Tracef("Count[%v]: %v", chid, len(buf))

					agent, err := getAgent(addr)
					if err == nil {
						logger.Trace(agent)

						err = submitJob(storage, dir, agent, buf)
						if err != nil {
							logger.Error(err)
						}
					}

					//logger.Trace(cx.ToJSON(buf))
					buf = make([]*cx.File, 0)
				}
			}
			if len(buf) > 0 {
				atomic.AddInt64(&download2, int64(len(buf)))
				logger.Tracef("Count[%v]: %v", chid, len(buf))

				agent, err := getAgent(addr)
				if err == nil {
					err = submitJob(storage, dir, agent, buf)
					if err != nil {
						logger.Error(err)
					}
				} else {
					logger.Error(err)
				}
				//logger.Trace(cx.ToJSON(buf))
				buf = make([]*cx.File, 0)
			}
		}(i)
	}

	err = s3svc.ListWithSplit(keysChan)
	if err != nil {
		logger.Error(err)
	}

	close(keysChan)
	wg.Wait()

	logger.Info("Finished.")
	dt := time.Now().Sub(t)
	logger.Infof("Down time: %v, TotalSize: %+v (%v), FileCount: %+v\n", dt, humanize.Bytes(uint64(downloaded)), humanize.Comma(downloaded), downloadedfile)
	logger.Infof("Download 2# %v (splitted: %v)", download2, splited)

	if conf.Wait {
		ch := make(chan int)
		<-ch
	}

	return nil
}

func getAgent(addr string) (*util.EventMessage, error) {
	//Alive check
	client := resty.New()
	client.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "s3transfer distributed",
	})
	client.SetTimeout(3 * time.Second)
	url := fmt.Sprintf("http://%v/agent/alloc", addr)

	resp, _ := client.R().Get(url)

	var v util.EventMessage
	reader := bytes.NewReader(resp.Body())
	err := json.NewDecoder(reader).Decode(&v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func submitJob(storage string, dir string, v *util.EventMessage, files []*cx.File) error {
	client := resty.New()
	client.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "s3transfer distributed",
	})

	//POST -> submitjob
	url := fmt.Sprintf("http://%v:%v/job/%v?storage=%v&dir=%v", v.IP, v.Port, v.ID, storage, dir)
	logger.Trace(url)
	resp, _ := client.R().
		SetBody(files).
		Put(url)

	logger.Trace(resp)

	return nil
}
