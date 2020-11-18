package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ditto.co.jp/agentserver/config"
	"ditto.co.jp/agentserver/db"
	"ditto.co.jp/agentserver/util"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

func initConfig(stage string) (*config.AppSetting, *db.Database, error) {
	conf, err := config.GetConfig(stage)
	if err != nil {
		return nil, nil, err
	}
	db := db.NewDatabase()

	return conf, db, nil
}

func dirWindows() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// Prefer standard environment variable USERPROFILE
	if home := os.Getenv("USERPROFILE"); home != "" {
		return home, nil
	}

	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, or USERPROFILE are blank")
	}

	return home, nil
}

func updateConfig(appsetting *config.AppSetting, opts *options) *config.Config {
	ip, err := util.ExternalIP()
	if err != nil {
		ip = "127.0.0.1"
	}

	conf := config.Config{
		Cmd:           strings.ToLower(opts.Cmd),
		Dir:           filepath.ToSlash(opts.Dir),
		S3:            opts.S3,
		Host:          appsetting.Server.Host, // Server hostname
		Port:          appsetting.Server.Port, // Server port
		LocalHost:     ip,
		AccessKey:     appsetting.Aws.AccessKeyID,
		SecretKey:     appsetting.Aws.SecretAccessKey,
		Region:        appsetting.Aws.Region,
		Token:         opts.Token,
		Proxy:         appsetting.Common.Proxy,
		BroadcastPort: appsetting.Common.BroadcastPort,
		Retry:         appsetting.Common.Retry,
		Timeout:       appsetting.Common.Timeout,
		JobID:         opts.JobID,
		Log:           appsetting.Server.LogLevel,
		Env:           opts.Env,
		MasterServer:  opts.MasterServer,
		DryRun:        opts.DryRun,
		Distributed:   opts.Distributed,
		ListSize:      appsetting.File.ListSize,
		Threads:       appsetting.File.Threads,
		PartSize:      uint64(appsetting.File.PartSize * config.MB),
		DistListSize:  appsetting.Dist.ListSize,
		DistThreads:   appsetting.Dist.Threads,
		DistPartSize:  uint64(appsetting.Dist.PartSize * config.MB),
		DistJobSize:   appsetting.Dist.JobSize,
		Storage:       opts.Storage,
		Wait:          opts.Wait,
	}

	//Server
	if len(opts.Host) > 0 {
		conf.Host = opts.Host
	}
	if opts.Port != 0 {
		conf.Port = opts.Port
	} else {
		if conf.Cmd == "agent" {
			conf.Port++
		}
	}
	if len(opts.Log) > 0 {
		conf.Log = opts.Log
	}
	if len(conf.MasterServer) == 0 {
		conf.MasterServer = fmt.Sprintf("%v:9090", conf.LocalHost)
	}
	//AWS params -> environment -> app.settings ->
	//environment
	aws := config.AwsSetting{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Region:          os.Getenv("AWS_DEFAULT_REGION"),
		Token:           os.Getenv("AWS_SESSION_TOKEN"),
	}
	if len(aws.AccessKeyID) > 0 {
		conf.AccessKey = aws.AccessKeyID
	}
	if len(aws.SecretAccessKey) > 0 {
		conf.SecretKey = aws.SecretAccessKey
	}
	if len(aws.Region) > 0 {
		conf.Region = aws.Region
	}
	if len(aws.Token) > 0 {
		conf.Token = aws.Token
	}
	// parameter
	if len(opts.AccessKey) > 0 {
		conf.AccessKey = opts.AccessKey
	}
	if len(opts.SecretKey) > 0 {
		conf.SecretKey = opts.SecretKey
	}
	if len(opts.Region) > 0 {
		conf.Region = opts.Region
	}
	if len(opts.Token) > 0 {
		conf.Token = opts.Token
	}

	if len(conf.AccessKey) == 0 || len(conf.SecretKey) == 0 {
		homedir, err := dirWindows()
		if err != nil {
			homedir = "~"
		} else {
			homedir = filepath.ToSlash(homedir)
		}

		credsfile := fmt.Sprintf("%s/.aws/credentials", homedir)
		creds := credentials.NewSharedCredentials(credsfile, "default")
		credValue, err := creds.Get()
		if err == nil {
			conf.AccessKey = credValue.AccessKeyID
			conf.SecretKey = credValue.SecretAccessKey
		}
	}

	if len(opts.Proxy) > 0 {
		conf.Proxy = opts.Proxy
	}
	if opts.Retry > 0 {
		conf.Retry = opts.Retry
	}
	if opts.Port > 0 {
		conf.Port = opts.Port
	}

	if conf.ListSize == 0 {
		conf.ListSize = 1000
	}
	if conf.Threads == 0 {
		conf.Threads = 64
	}

	if conf.DistListSize == 0 {
		conf.DistListSize = 1000
	}
	if conf.DistThreads == 0 {
		conf.DistThreads = 5
	}
	if conf.DistJobSize == 0 {
		conf.DistJobSize = 300
	}

	return &conf
}
