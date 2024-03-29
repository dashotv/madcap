// Code generated by oto; DO NOT EDIT.

package server

import (
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/dashotv/grimoire"
)

func setupDatabase(s *Server) error {
	db := &connection{
		log: s.Logger.Named("db"),
	}

	dbname := s.Config.Name + "_development"
	if s.Config.Production {
		dbname = s.Config.Name + "_production"
	}

	if col, err := grimoire.New[*File](s.Config.Mongo, dbname, strings.ToLower("File")); err != nil {
		return fmt.Errorf("failed to create File collection: %w", err)
	} else {
		db.File = col
		grimoire.Indexes(col, &File{})
	}

	s.db = db
	return nil
}

type connection struct {
	log *zap.SugaredLogger

	File *grimoire.Store[*File]
}

type EmptyResponse struct {
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
}

// File represents a file in the file system.
type File struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Path       string `json:"path" bson:"path" grimoire:"index" `
	Size       int64  `json:"size" bson:"size" `
	ModifiedAt int64  `json:"modified_at" bson:"modified_at" grimoire:"index" `
}

type FilesResponse struct {
	Count  int64   `json:"count"`
	Result []*File `json:"files"`
	// Error is string explaining what went wrong. Empty if everything was fine.
	Error string `json:"error,omitempty"`
}

type IndexRequest struct {
	Limit int `json:"limit"`
	Skip  int `json:"skip"`
}

type KeyRequest struct {
	Key string `json:"id"`
}
