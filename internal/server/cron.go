package server

import (
	"io/fs"
	"path/filepath"
	"time"

	"github.com/robfig/cron/v3"
)

var walking bool

func setupCron(s *Server) {
	s.Cron = cron.New()

	// s.Cron.AddFunc("*/30 * * * *", s.walkFiles)
}

func (s *Server) walkFiles() {
	if walking {
		s.Logger.Infow("already walking")
		return
	}

	walking = true
	defer func() {
		walking = false
	}()

	libs, err := s.Plex.GetLibraries()
	if err != nil {
		s.Logger.Errorw("failed to get libraries", "error", err)
		return
	}

	start := time.Now()
	defer func() {
		s.Logger.Infow("walked", "duration", time.Since(start))
	}()

	for _, lib := range libs {
		for _, loc := range lib.Locations {
			s.Logger.Infow("walking", "library", lib.Title, "location", loc)
			err := filepath.WalkDir(loc.Path, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					s.Logger.Errorw("failed to walk", "path", path, "error", err)
					return err
				}
				if d.IsDir() {
					return nil
				}
				s.Logger.Infow("found file", "path", path)
				return nil
			})
			if err != nil {
				s.Logger.Errorw("failed to walk", "path", loc.Path, "error", err)
			}
		}
	}
	return
}
