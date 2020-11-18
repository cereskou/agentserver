package cmd

import (
	"fmt"
	"strings"

	"ditto.co.jp/agentserver/config"
	"ditto.co.jp/agentserver/logger"
	"ditto.co.jp/agentserver/util"
)

//RunList -
func RunList(conf *config.Config, cmd string) error {
	//
	host := fmt.Sprintf("%v:%v", conf.Host, conf.Port)
	switch strings.ToLower(cmd) {
	case "job":
		files, err := util.ListJob(host)
		if err != nil {
			return err
		}
		logger.Printf("Files: %v\n", len(files))
	case "agent":
		agents, err := util.ListAgent(host)
		if err != nil {
			return err
		}
		logger.Printf("Counts: %v\n", len(agents))
		for _, agent := range agents {
			logger.Printf("Agent: %v:%v\n", agent.IP, agent.Port)
		}
	case "status":
		agents, err := util.ListAgent(host)
		if err != nil {
			return err
		}
		for _, agent := range agents {
			addr := fmt.Sprintf("%v:%v", agent.IP, agent.Port)
			files, err := util.ListJob(addr)
			if err != nil {
				continue
			}
			logger.Printf("Agent: %v, Files: %v\n", addr, len(files))
		}

	default:
		logger.Printf("Usage:\n\tcommand list [agent,job]\n")
	}

	return nil
}
