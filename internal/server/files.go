package server

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type fileService struct {
	log *zap.SugaredLogger
}

func (s *fileService) Index(c echo.Context, req *IndexRequest) (*FilesResponse, error) {
	return &FilesResponse{}, nil
}
