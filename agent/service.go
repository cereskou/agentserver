package agent

import (
	"ditto.co.jp/agentserver/config"
	"ditto.co.jp/agentserver/db"
)

//Service -
type Service struct {
	_config  *config.Config
	_db      *db.Database
	prefix   string
	dirsSeen map[string]bool
}

//NewService -
func NewService(conf *config.Config, db *db.Database) *Service {
	return &Service{
		_config:  conf,
		_db:      db,
		dirsSeen: make(map[string]bool),
	}
}

//Close -
func (s *Service) Close() {
}
