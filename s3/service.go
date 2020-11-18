package s3

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"ditto.co.jp/agentserver/config"
	"ditto.co.jp/agentserver/cx"
	"ditto.co.jp/agentserver/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cenkalti/backoff"
	"github.com/dustin/go-humanize"
)

//Service -
type Service struct {
	_config          *config.Config
	_client          *s3.S3
	_session         *session.Session
	_s3              *URI
	dirsSeen         map[string]bool
	newDirPermission os.FileMode
}

//NewService -
func NewService(cfg *config.Config) *Service {
	s3sv := &Service{
		_config:          cfg,
		_s3:              &URI{},
		dirsSeen:         make(map[string]bool),
		newDirPermission: os.ModePerm,
	}

	sess := s3sv.getSession()
	if sess == nil {
		return nil
	}

	s3sv._session = sess
	s3sv._client = s3.New(sess)
	if s3sv._client == nil {
		return nil
	}

	return s3sv
}

func (s *Service) getSession() *session.Session {
	//Proxy
	var httpClient *http.Client
	if len(s._config.Proxy) > 0 {
		httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: func(*http.Request) (*url.URL, error) {
					return url.Parse(s._config.Proxy)
				},
			},
		}
	}

	//認証情報を作成します。
	cred := credentials.NewStaticCredentials(
		s._config.AccessKey,
		s._config.SecretKey,
		s._config.Token)

	//セッション作成します
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(s._config.Region),
		Credentials: cred,
		HTTPClient:  httpClient,
		MaxRetries:  aws.Int(s._config.Retry),
	}))

	return sess
}

//Downlad -
func (s *Service) Downlad(files []cx.File) error {
	//Download
	keysChan := make(chan *File, s._config.ListSize)
	wg := new(sync.WaitGroup)

	var downloaded int64 = 0
	var downloadedfile int64 = 0

	var i uint
	for i = 0; i < s._config.Threads; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for fi := range keysChan {
				t := time.Now()

				if s._config.DryRun {
					fmt.Printf("\"%v\",\"%v\",%v\n", fi.S3, fi.SaveTo, fi.Size)
					continue
				}

				retry := s._config.Retry
				for {
					var err error
					err = s.downladPart(fi, fi.Num, fi.Offset, fi.Length)
					if err == nil {
						atomic.AddInt64(&downloaded, fi.Size)
						atomic.AddInt64(&downloadedfile, 1)

						break
					}

					fmt.Printf("\"%v\",\"%v\",%v,\"%v\"\n", fi.S3, fi.SaveTo, 0, err.Error())
					os.Remove(fi.SaveTo)

					time.Sleep(2 * time.Second)
					retry--
					if retry <= 0 {
						break
					}
				}

				elapsed := int64(time.Now().Sub(t)) / (int64(time.Millisecond) / int64(time.Nanosecond))
				if retry <= 0 {
					//Failed
					fmt.Printf("\"%v\",\"%v\",%v,\"Failed\"\n", fi.S3, fi.SaveTo, elapsed)
				} else {
					//Successed
					fmt.Printf("\"%v\",\"%v\",%v\n", fi.S3, fi.SaveTo, elapsed)
				}
			}
		}()
	}

	t := time.Now()
	var totalsize int64 = 0
	filecount := 0

	// Get download files from a list file.
	for _, f := range files {
		u, _ := url.Parse(f.S3)
		if u.Scheme != "s3" {
			continue
		}
		save := f.Local

		f := &File{
			Bucket: u.Host,
			Key:    u.Path,
			Size:   int64(f.Size),
			SaveTo: save,
			S3:     f.S3,
			Num:    f.Num,
			Offset: f.Offset,
			Length: f.Length,
		}

		totalsize += f.Size
		filecount++

		keysChan <- f
	}

	ft := time.Now().Sub(t)
	fmt.Println(fmt.Sprintf("Count: %v", filecount))

	close(keysChan)

	wg.Wait()
	dt := time.Now().Sub(t)
	fmt.Printf("Scan time: %v, TotalSize: %+v (%v), FileCount: %+v\n", ft, humanize.Bytes(uint64(totalsize)), humanize.Comma(totalsize), filecount)
	fmt.Printf("Down time: %v, TotalSize: %+v (%v), FileCount: %+v\n", dt, humanize.Bytes(uint64(downloaded)), humanize.Comma(downloaded), downloadedfile)
	return nil
}

//Downlad -
func (s *Service) downladPart(obj *File, num uint32, offset, size uint64) error {

	input := &s3.GetObjectInput{
		Bucket: aws.String(obj.Bucket),
		Key:    aws.String(obj.Key),
		Range:  aws.String(fmt.Sprintf("bytes=%d-%d", offset, offset+size-1)),
	}
	if offset == 0 && size == 0 {
		input.Range = nil
	}

	var output *s3.GetObjectOutput
	operation := func() error {
		var err error
		output, err = s._client.GetObject(input)
		return err
	}

	notify := func(err error, duration time.Duration) {
		logger.Printf("%v %v\n", duration, err)
	}

	err := backoff.RetryNotify(
		operation,
		backoff.WithMaxRetries(backoff.NewExponentialBackOff(), uint64(s._config.Retry)),
		notify)
	if err != nil {
		logger.Println(err)
		return err
	}

	defer output.Body.Close()

	// filename := filepath.Base(obj.SaveTo)
	filename := obj.SaveTo
	//Send data to central server
	err = s.transfer(output.Body, filename, num)
	if err != nil {
		logger.Println(err)
		return err
	}

	return nil
}

func (s *Service) transfer(file io.Reader, filename string, num uint32) error {
	//upload
	upload := fmt.Sprintf("http://%v:%v/%v?num=%v&dir=%v",
		s._config.Host,
		s._config.Port,
		s._config.JobID,
		num,
		s._config.Dir)

	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)
	go func() {
		defer pw.Close()

		fw, err := mw.CreateFormFile("file", filename)
		if err != nil {

		}
		_, err = io.Copy(fw, file)
		if err != nil {

		}
		if err := mw.Close(); err != nil {

		}
	}()
	contentType := mw.FormDataContentType()
	res, err := http.Post(upload, contentType, pr)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	message, _ := ioutil.ReadAll(res.Body)
	fmt.Printf(string(message))

	return nil
}
