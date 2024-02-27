package server

import (
	"github.com/robfig/cron/v3"
)

func setupCron(s *Server) {
	s.Cron = cron.New()
}
