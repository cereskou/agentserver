package s3

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"ditto.co.jp/agentserver/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

//
var (
	ErrInvalidArgument     = errors.New("Invalid argument type")
	ErrParameterValidation = errors.New("Parameter validation failed")
)

func (s *Service) SetURI(uri string) error {
	u, err := url.Parse(uri)
	if err != nil {
		logger.Errorf("Invalide S3 URI '%s'", uri)
		return err
	}
	if u.Scheme != "s3" {
		logger.Errorf("Invalid argument type '%s'", uri)
		return ErrInvalidArgument
	}
	//マッチチェック
	matched, err := regexp.MatchString(`^[a-zA-Z0-9.\-_]{1,255}$`, u.Host)
	if (err != nil) || !matched {
		logger.Errorf("Parameter validation failed:\n"+`Invalid bucket name "%+v": Bucket name must match the regex "^[a-zA-Z0-9.\-_]{1,255}$"`, u.Host)
		return ErrParameterValidation
	}
	bucket := u.Host
	prefix := u.Path
	if len(prefix) > 0 && prefix[0] == '/' {
		prefix = prefix[1:]
	}

	s._s3.Bucket = bucket
	s._s3.Prefix = prefix

	return nil
}

//List -
func (s *Service) List(keysChan chan *File) error {

	req := &s3.ListObjectsV2Input{
		Bucket: aws.String(s._s3.Bucket),
		Prefix: aws.String(s._s3.Prefix),
	}

	err := s._client.ListObjectsV2Pages(req, func(resp *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, content := range resp.Contents {
			var md5 string
			if s._config.Md5 {
				md5 = string(*content.ETag)
				md5 = strings.ReplaceAll(md5, "\n", "")
			}
			key := string(*content.Key)
			s3 := fmt.Sprintf("s3://%v/%v", s._s3.Bucket, key)
			fi := &File{
				Bucket:  s._s3.Bucket,
				Key:     key,
				SaveTo:  key,
				Size:    int64(*content.Size),
				Md5:     md5,
				ModTime: time.Time(*content.LastModified),
				S3:      s3,
			}

			keysChan <- fi
		}

		return true
	})
	if err != nil {
		return err
	}

	return nil
}

//ListV2 -
func (s *Service) ListV2(keysChan chan *File) error {
	ch := make(chan int)

	go func() {
		token := s.lists3(keysChan, nil)
		for token != nil {
			token = s.lists3(keysChan, token)
		}
		ch <- 1
	}()
	<-ch

	return nil
}

func (s *Service) lists3(out chan *File, token *string) *string {
	req := &s3.ListObjectsV2Input{
		Bucket:            aws.String(s._s3.Bucket),
		Prefix:            aws.String(s._s3.Prefix),
		ContinuationToken: token,
	}

	list, err := s._client.ListObjectsV2(req)
	if err != nil {
		out <- &File{
			Err: err,
		}
		return nil
	}

	for _, obj := range list.Contents {

		md5 := *obj.ETag

		s3 := fmt.Sprintf("s3://%v/%v", s._s3.Bucket, *obj.Key)

		out <- &File{
			Bucket:  s._s3.Bucket,
			Key:     *obj.Key,
			SaveTo:  *obj.Key,
			Size:    int64(*obj.Size),
			Md5:     md5,
			ModTime: *obj.LastModified,
			S3:      s3,
		}
	}

	return list.NextContinuationToken
}

//ListWithSplit - split big file to small parts
func (s *Service) ListWithSplit(keysChan chan *File) error {

	req := &s3.ListObjectsV2Input{
		Bucket: aws.String(s._s3.Bucket),
		Prefix: aws.String(s._s3.Prefix),
	}

	err := s._client.ListObjectsV2Pages(req, func(resp *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, content := range resp.Contents {
			var md5 string
			if s._config.Md5 {
				md5 = string(*content.ETag)
				md5 = strings.ReplaceAll(md5, "\n", "")
			}
			key := string(*content.Key)
			s3 := fmt.Sprintf("s3://%v/%v", s._s3.Bucket, key)

			// Part計算
			size := s._config.DistPartSize
			objsize := uint64(*content.Size)
			parts := (objsize / size)
			if objsize%size > 0 {
				parts++
			}
			if parts > 1 {
				//複数分割
				var x uint64 = 0
				for x = 0; x < parts; x++ {
					offset := x * size
					if offset+size > objsize {
						size = objsize - offset
					}
					fi := &File{
						Bucket:  s._s3.Bucket,
						Key:     key,
						SaveTo:  key,
						Size:    int64(size),
						Md5:     md5,
						ModTime: time.Time(*content.LastModified),
						S3:      s3,
						Num:     uint32(x + 1),
						Offset:  offset,
						Length:  size,
					}

					keysChan <- fi
				}
			} else {
				//分割不要
				fi := &File{
					Bucket:  s._s3.Bucket,
					Key:     key,
					SaveTo:  key,
					Size:    int64(*content.Size),
					Md5:     md5,
					ModTime: time.Time(*content.LastModified),
					S3:      s3,
				}

				keysChan <- fi
			}
		}

		return true
	})
	if err != nil {
		return err
	}

	return nil
}
