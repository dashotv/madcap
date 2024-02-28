package server

import (
	"github.com/robfig/cron/v3"
)

var walking bool

func setupCron(s *Server) {
	s.Cron = cron.New()

	// s.Cron.AddFunc("*/30 * * * *", s.walkFiles)
}

func (s *Server) walkFiles() {
	if walking {
		s.Logger.Info("already walking")
		return
	}

	walking = true
	defer func() {
		walking = false
	}()

	libs, err := s.Plex.GetLibraries()
	if err != nil {
		s.Logger.Errorw("libs", "error", err)
	}

	w := newWalker(s.db, s.Logger.Named("walker"), libs)
	if err := w.Walk(); err != nil {
		s.Logger.Errorw("walk", "error", err)
	}
}
