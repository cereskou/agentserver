package main

import (
	"os"
	"runtime"
	"strings"

	"ditto.co.jp/agentserver/cmd"
	"ditto.co.jp/agentserver/logger"
	flags "github.com/jessevdk/go-flags"
)

//VERSION -
var version string

// options -
type options struct {
	Cmd          string `long:"cmd" description:"Command [agent, master, execute]" default:"execute"`
	Dir          string `long:"dir" description:"work directory"`
	S3           string `long:"s3" description:"s3 url"`
	Host         string `short:"H" long:"host" description:"master host or ip"`
	Port         uint   `short:"p" long:"port" description:"listen port"`
	AccessKey    string `long:"accesskey" description:"AWS access key id."`
	SecretKey    string `long:"secretkey" description:"AWS sercret key"`
	Region       string `long:"region" description:"AWS region" default:"ap-northeast-1"`
	Token        string `long:"token" description:"aws credentials token"`
	Proxy        string `long:"proxy" description:"Proxy server"`
	Retry        int    `long:"retry" description:"Retry count" default:"3"`
	JobID        string `long:"jobid" description:"Job ID when cmd=worker"`
	Log          string `long:"log" description:"Log level (trace,debug,info,warn,error,fatal)"`
	Distributed  bool   `short:"d" long:"distributed" description:"work with execute."`
	DryRun       bool   `long:"dryrun" description:"debug mode"`
	Env          string `short:"e" long:"env" description:"Stage [deve]"`
	AgentServer  string `short:"a" long:"agent" description:"agent server [host:port]"`
	MasterServer string `short:"m" long:"master" description:"master server [host:port]"`
	Storage      string `short:"s" long:"storage" description:"storage [host:port]"` //最後データ格納サーバー
	Wait         bool   `long:"wait" description:"wait for download finish"`         //only in execute mode
	Version      bool   `short:"v" long:"version" description:"Show version"`        //バージョン表示
}

// @title AgentServer API
// @version 1.0
// @description This is a s3transfer agent server.
// @termsOfService http://www.ditto.co.jp

// @contact.name API Support
// @contact.url ditto.co.jp/support
// @contact.email support@ditto.co.jp

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host www.ditto.co.jp
// @BasePath /v2
func main() {
	var optional string
	var opts options
	if cmd, err := flags.Parse(&opts); err != nil {
		os.Exit(-1)
	} else {
		if len(cmd) > 0 {
			opts.Cmd = cmd[0]
		}
		if len(cmd) > 1 {
			optional = cmd[1]
		}
	}
	if opts.Version {
		logger.Printf("s3transfer agent %+v\n", version)
		os.Exit(0)
	}
	appsetting, db, err := initConfig(opts.Env)
	if err != nil {
		logger.Printf("Failed initialize configurations or db.")
		os.Exit(-1)
	}

	//Setup Config
	conf := updateConfig(appsetting, &opts)

	// Log
	logger.SetLevel(conf.Log)

	//to lower
	conf.Cmd = strings.ToLower(conf.Cmd)

	if conf.Cmd == "worker" {
		if len(conf.JobID) == 0 {
			logger.Error("Please specify the jobid when worker mode")
			os.Exit(-1)
		}
		conf.Host = opts.Host
		conf.Port = opts.Port
		if len(conf.Host) == 0 || conf.Port == 0 {
			logger.Error("Please specify the agent server")
			os.Exit(-1)
		}

		logger.Info("Run Worker")
		err := cmd.RunWorker(conf, db)
		if err != nil {
			logger.Error(err)
			os.Exit(-1)
		}
	}

	if conf.Cmd == "agent" {
		logger.Info("Run Agent")
		err := cmd.RunAgentServer(conf, db)
		if err != nil {
			logger.Error(err)
			os.Exit(-1)
		}
	}

	if conf.Cmd == "master" {
		logger.Info("Run Master")
		err := cmd.RunMasterServer(conf, db)
		if err != nil {
			logger.Error(err)
			os.Exit(-1)
		}
	}

	if opts.Cmd == "execute" {
		err := cmd.RunExecute(conf, db)
		if err != nil {
			logger.Error(err)
			os.Exit(-1)
		}
	}

	if opts.Cmd == "list" {
		err := cmd.RunList(conf, optional)
		if err != nil {
			logger.Error(err)
			os.Exit(-1)
		}
	}

	if opts.Cmd == "start" {
		err := cmd.RunStart(conf)
		if err != nil {
			logger.Error(err)
			os.Exit(-1)
		}
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
