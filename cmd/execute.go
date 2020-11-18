package cmd

import (
	"errors"

	"ditto.co.jp/agentserver/config"
	"ditto.co.jp/agentserver/db"
	"ditto.co.jp/agentserver/execute"
	"ditto.co.jp/agentserver/logger"
)

// -
var (
	ErrURLNotFound = errors.New("s3 not found")
)

//RunExecute -
func RunExecute(conf *config.Config, db *db.Database) error {
	if len(conf.S3) == 0 {
		logger.Error("Please sepcify a s3 path with --s3")
		return ErrURLNotFound
	}

	//分散処理
	if conf.Distributed {
		return execute.Distributed(conf)
	}

	return execute.Download(conf)
}
