package util

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"ditto.co.jp/agentserver/config"
	"ditto.co.jp/agentserver/cx"
	"ditto.co.jp/agentserver/logger"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
)

//-
var (
	ErrNotFound = errors.New("not found")
	ErrOthers   = errors.New("error occurred")
)

//GetJob -
func GetJob(host string, jobid string) ([]cx.File, error) {
	url := fmt.Sprintf("http://%v/job/%v", host, jobid)

	client := resty.New()
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Get(url)
	if err != nil {
		return nil, err
	}

	// Not Found
	if resp.StatusCode() == http.StatusNotFound {
		return nil, ErrNotFound
	}

	var v []cx.File = make([]cx.File, 0)

	reader := bytes.NewReader(resp.Body())
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.NewDecoder(reader).Decode(&v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

//GetJobString -
func GetJobString(host string, jobid string) (string, error) {
	url := fmt.Sprintf("http://%v/job/%v", host, jobid)

	client := resty.New()
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Get(url)

	if err != nil {
		return "", err
	}

	// Not Found
	if resp.StatusCode() == http.StatusNotFound {
		return "", ErrNotFound
	}

	return string(resp.Body()), nil
}

//GetJobReader -
func GetJobReader(host string, jobid string) (io.Reader, error) {
	url := fmt.Sprintf("http://%v/job/%v", host, jobid)

	client := resty.New()
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Get(url)

	if err != nil {
		return nil, err
	}

	// Not Found
	if resp.StatusCode() == http.StatusNotFound {
		return nil, ErrNotFound
	}

	return bytes.NewReader(resp.Body()), nil
}

//ServerCheck -
func ServerCheck(host string) (alive string) {
	url := fmt.Sprintf("http://%v/health", host)

	client := resty.New()
	client.SetHeaders(map[string]string{
		"User-Agent": "s3transfer distributed",
	})
	resp, err := client.R().
		Get(url)

	// Not Found
	if resp.StatusCode() == http.StatusNotFound {
		return
	}

	if err != nil {
		return
	}

	alive = string(resp.Body())

	return
}

//ListJob -
func ListJob(host string) ([]cx.File, error) {
	url := fmt.Sprintf("http://%v/job/list", host)

	client := resty.New()
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Get(url)
	if err != nil {
		return nil, err
	}

	// Not Found
	if resp.StatusCode() == http.StatusNotFound {
		return nil, ErrNotFound
	}

	var v [][]cx.File = make([][]cx.File, 0)

	reader := bytes.NewReader(resp.Body())
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.NewDecoder(reader).Decode(&v)
	if err != nil {
		return nil, err
	}
	var result []cx.File = make([]cx.File, 0)
	for _, n := range v {
		for _, s := range n {
			result = append(result, s)
		}
	}

	return result, nil
}

//KickAgent -
func KickAgent(host string, conf *config.Config) error {
	//POST
	url := fmt.Sprintf("http://%v/job/exec?storage=%v&dir=%v", host, conf.Storage, conf.Dir)
	client := resty.New()
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Post(url)
	if err != nil {
		return err
	}

	logger.Debug(string(resp.Body()))

	return nil
}

//ListAgent -
func ListAgent(host string) ([]EventMessage, error) {
	url := fmt.Sprintf("http://%v/agent/list", host)

	client := resty.New()
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Get(url)
	if err != nil {
		return nil, err
	}

	// Not Found
	if resp.StatusCode() == http.StatusNotFound {
		return nil, ErrNotFound
	}

	var v []EventMessage = make([]EventMessage, 0)

	reader := bytes.NewReader(resp.Body())
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.NewDecoder(reader).Decode(&v)
	if err != nil {
		return nil, err
	}

	return v, nil
}
