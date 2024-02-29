package server

import (
	"os"
	"sync/atomic"
	"time"

	"github.com/robfig/cron/v3"
	"golang.org/x/sync/errgroup"
)

var walking uint32
var updating uint32

func setupCron(s *Server) {
	s.Cron = cron.New()

	s.Cron.AddFunc("*/15 * * * *", s.walkFiles)
	s.Cron.AddFunc("5 * * *", s.updateFiles)
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

func (s *Server) updateFiles() {
	l := s.Logger.Named("updater")

	if !atomic.CompareAndSwapUint32(&updating, 0, 1) {
		l.Warnf("already running")
	}
	defer atomic.StoreUint32(&updating, 0)

	c := &counter{v: 0}
	start := time.Now()
	defer func() { l.Warnw("complete", "duration", time.Since(start), "count", c.v) }()

	total, err := s.db.File.Query().Count()
	if err != nil {
		l.Errorw("total", "error", err)
		return
	}
	l.Debugf("total %d", total)

	files := make(chan *File, 10)

	eg := new(errgroup.Group)
	eg.Go(func() error {
		defer close(files)

		for i := 0; i < int(total); i += 100 {
			list, err := s.db.File.Query().Skip(i).Limit(100).Run()
			if err != nil {
				return err
			}

			for _, file := range list {
				files <- file
			}
		}
		return nil
	})

	for i := 0; i < 10; i++ {
		eg.Go(func() error {
			for file := range files {
				c.Inc()
				if err := s.updateFile(file); err != nil {
					return err
				}
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		l.Errorw("updateFiles", "error", err)
	}
}

func (s *Server) updateFile(file *File) error {
	info, err := os.Lstat(file.Path)
	if err != nil {
		return err
	}

	file.Size = info.Size()
	file.ModifiedAt = info.ModTime().Format(time.RFC3339)

	// TODO: if video, ffmpeg info

	if err := s.db.File.Save(file); err != nil {
		return err
	}

	return nil
}
