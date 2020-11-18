package job

import (
	"ditto.co.jp/agentserver/config"
	"ditto.co.jp/agentserver/db"
	"ditto.co.jp/agentserver/logger"
)

//Service -
type Service struct {
	_config *config.Config
	_db     *db.Database
	prefix  string
}

//NewService -
func NewService(conf *config.Config, db *db.Database) (*Service, error) {
	err := db.CreateIndex("jobs", "job:*")
	if err != nil {
		logger.Error(err)

		return nil, err
	}
	return &Service{
		_config: conf,
		_db:     db,
	}, nil
}

//Close -
func (s *Service) Close() {
}
