package server

import (
	"fmt"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/dashotv/madcap/internal/plex"
	"github.com/dashotv/madcap/internal/workers"
)

func newWalker(db *connection, logger *zap.SugaredLogger, libs []*plex.PlexLibrary) *Walker {
	return &Walker{
		db:          db,
		logger:      logger,
		Libraries:   libs,
		directories: make(chan string, 10),
		files:       make(chan string, 10),
	}
}

type Walker struct {
	db          *connection
	bg          *workers.Client
	logger      *zap.SugaredLogger
	Libraries   []*plex.PlexLibrary
	directories chan string
	files       chan string
}

type counter struct {
	sync.Mutex
	v int
}

func (w *Walker) Walk() error {
	start := time.Now()
	defer func() { w.logger.Infow("walk", "duration", time.Since(start)) }()

	c := counter{v: 0}
	defer func() { w.logger.Infow("files", "count", c.v) }()

	eg := new(errgroup.Group)

	eg.Go(func() error {
		defer close(w.directories)

		if err := w.getDirectories(); err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		defer close(w.files)

		for dir := range w.directories {
			if err := w.getDirFiles(dir); err != nil {
				return err
			}
		}
		return nil
	})

	for i := 0; i < 10; i++ {
		eg.Go(func() error {
			for line := range w.files {
				c.Lock()
				c.v++
				c.Unlock()
				if err := w.createFile(line); err != nil {
					return err
				}
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func (w *Walker) createFile(line string) error {
	file, err := w.db.FileByPath(line)
	if err != nil {
		return err
	}

	if file.ID != primitive.NilObjectID {
		return nil
	}

	if err := w.db.File.Save(file); err != nil {
		return err
	}

	return nil
}

func (w *Walker) updateFile(line string) error {
	info, err := os.Stat(line)
	if err != nil {
		w.logger.Errorw("stat", "error", err)
		return err
	}

	file, err := w.db.FileByPath(line)
	if err != nil {
		w.logger.Errorw("file", "error", err)
		return err
	}

	file.ModifiedAt = info.ModTime().Format(time.RFC3339)
	file.Size = info.Size()

	if err := w.db.File.Save(file); err != nil {
		w.logger.Errorw("save", "error", err)
		return err
	}

	return nil
}

func (w *Walker) getDirectories() error {
	out := func(name, line string) {
		// w.logger.Debugf("sending directory: %s", line)
		w.directories <- line
	}
	err := func(name, line string) {
		w.logger.Errorw("shell", "error", line)
	}

	for _, lib := range w.Libraries {
		if lib.Type != "show" && lib.Type != "movie" {
			continue
		}
		for _, loc := range lib.Locations {
			// w.logger.Infow("walking", "library", lib.Title, "path", loc.Path)
			// workers.Shell(w.logger.Named("shell"), loc.Path, "find", loc.Path, "-type", "f", "-exec", "stat", "{}", "+")
			status, err := workers.Shell(fmt.Sprintf("find '%s' -maxdepth 1 -mindepth 1 -type d", loc.Path), workers.ShellOptions{Out: out, Err: err})
			if err != nil {
				return err
			}
			if status.Exit != 0 {
				return fmt.Errorf("command failed with exit code %d", status.Exit)
			}
		}
	}

	return nil
}

// func (w *Walker) getFiles() error {
// 	for _, dir := range w.Directories {
// 		if err := w.getDirFiles(dir); err != nil {
// 			return err
// 		}
// 	}
//
// 	return nil
// }

func (w *Walker) getDirFiles(dir string) error {
	out := func(name, line string) {
		w.files <- line
	}
	err := func(name, line string) {
		w.logger.Errorw("shell", "error", line)
	}

	// status, err := workers.Shell(out, err, "find", dir, "-type", "f", "-exec", "stat", "{}", "+")
	status, e := workers.Shell(fmt.Sprintf("find '%s' -type f", dir), workers.ShellOptions{Out: out, Err: err})
	if e != nil {
		return e
	}
	if status.Exit != 0 {
		return fmt.Errorf("command failed with exit code %d", status.Exit)
	}
	return nil
}
