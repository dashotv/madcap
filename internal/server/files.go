package server

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type fileService struct {
	log    *zap.SugaredLogger
	server *Server
}

func (s *fileService) Index(c echo.Context, req *IndexRequest) (*FilesResponse, error) {
	return &FilesResponse{}, nil
}
func (s *fileService) Update(c echo.Context, req *KeyRequest) (*EmptyResponse, error) {
	go s.server.walkFiles()
	return &EmptyResponse{}, nil
}
