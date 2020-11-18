package master

import (
	"ditto.co.jp/agentserver/config"
	"ditto.co.jp/agentserver/db"
)

//Service -
type Service struct {
	_config *config.Config
	_db     *db.Database
	prefix  string
}

//NewService -
func NewService(conf *config.Config, db *db.Database) *Service {
	return &Service{
		_config: conf,
		_db:     db,
	}
}

//Close -
func (s *Service) Close() {
}
