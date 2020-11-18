package services

import (
	"reflect"

	"ditto.co.jp/agentserver/agent"
	"ditto.co.jp/agentserver/config"
	"ditto.co.jp/agentserver/db"
	"ditto.co.jp/agentserver/job"
	"ditto.co.jp/agentserver/master"
)

var (
	//AgentService ...
	AgentService agent.ServiceInterface

	//JobServer ...
	JobService job.ServiceInterface

	//MasterService ...
	MasterService master.ServiceInterface
)

// InitService -
func InitService(mode string, conf *config.Config, db *db.Database) (err error) {
	switch mode {
	case "agent":
		if nil == reflect.TypeOf(AgentService) {
			AgentService = agent.NewService(conf, db)
		}

		if nil == reflect.TypeOf(JobService) {
			JobService, err = job.NewService(conf, db)
		}

	case "master":
		if nil == reflect.TypeOf(MasterService) {
			MasterService = master.NewService(conf, db)
		}

	default:
	}

	return
}

//Close -
func Close() {
	if AgentService != nil {
		AgentService.Close()
	}
	if JobService != nil {
		JobService.Close()
	}
	if MasterService != nil {
		MasterService.Close()
	}
}
