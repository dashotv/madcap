package server

import (
	"sync/atomic"

	"github.com/robfig/cron/v3"
)

var walking uint32

func setupCron(s *Server) {
	s.Cron = cron.New()

	if s.Config.Production {
		s.Cron.AddFunc("*/15 * * * *", s.walkFiles)
	}
}

func (s *Server) walkFiles() {
	if !atomic.CompareAndSwapUint32(&walking, 0, 1) {
		s.Logger.Warnf("walkFiles: already running")
	}
	defer atomic.StoreUint32(&walking, 0)

	libs, err := s.Plex.GetLibraries()
	if err != nil {
		s.Logger.Errorw("libs", "error", err)
	}

	w := newWalker(s.db, s.Logger.Named("walker"), libs)
	if err := w.Walk(); err != nil {
		s.Logger.Errorw("walk", "error", err)
	}
}
