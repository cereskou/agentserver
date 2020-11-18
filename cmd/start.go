package cmd

import (
	"fmt"
	"path/filepath"

	"ditto.co.jp/agentserver/config"
	"ditto.co.jp/agentserver/logger"
	"ditto.co.jp/agentserver/util"
)

//RunStart -
func RunStart(conf *config.Config) error {
	//Find master
	logger.Info("Scan master node...")
	host := conf.Host
	logger.Info(host)
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
		nport := conf.Port
		nport++
		storage = fmt.Sprintf("%v:%v", conf.LocalHost, nport)
	}

	dir := filepath.Clean(conf.Dir)
	dir = filepath.ToSlash(dir)

	addr := fmt.Sprintf("%v:%v", host, port)
	logger.Tracef("master note : %v", addr)
	logger.Tracef("storage node: %v", storage)
	logger.Tracef("save to: %v", dir)
	//
	agents, err := util.ListAgent(addr)
	if err != nil {
		return err
	}
	for _, agent := range agents {
		addr := fmt.Sprintf("%v:%v", agent.IP, agent.Port)
		//
		err = util.KickAgent(addr, conf)
		if err != nil {
			logger.Error(err)
			continue
		}
	}

	return nil
}
