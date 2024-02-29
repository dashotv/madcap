// Code generated by oto; DO NOT EDIT.

package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type FileService interface {
	Index(echo.Context, *IndexRequest) (*FilesResponse, error)
	Stat(echo.Context, *KeyRequest) (*EmptyResponse, error)
	Update(echo.Context, *KeyRequest) (*EmptyResponse, error)
}

type fileServiceServer struct {
	fileService FileService
}

// Register adds the FileService to the otohttp.Server.
func RegisterFileService(e *echo.Group, fileService FileService) {
	handler := &fileServiceServer{
		fileService: fileService,
	}
	e.POST("/FileService.Index", handler.handleIndex)
	e.POST("/FileService.Stat", handler.handleStat)
	e.POST("/FileService.Update", handler.handleUpdate)
}

func (s *fileServiceServer) handleIndex(c echo.Context) error {
	request := &IndexRequest{}
	if err := c.Bind(request); err != nil {
		return fmt.Errorf("binding request: %w", err)
	}

	response, err := s.fileService.Index(c, request)
	if err != nil {
		return fmt.Errorf("handling request: %w", err)
	}

	return c.JSON(http.StatusOK, response)
}

func (s *fileServiceServer) handleStat(c echo.Context) error {
	request := &KeyRequest{}
	if err := c.Bind(request); err != nil {
		return fmt.Errorf("binding request: %w", err)
	}

	response, err := s.fileService.Stat(c, request)
	if err != nil {
		return fmt.Errorf("handling request: %w", err)
	}

	return c.JSON(http.StatusOK, response)
}

func (s *fileServiceServer) handleUpdate(c echo.Context) error {
	request := &KeyRequest{}
	if err := c.Bind(request); err != nil {
		return fmt.Errorf("binding request: %w", err)
	}

	response, err := s.fileService.Update(c, request)
	if err != nil {
		return fmt.Errorf("handling request: %w", err)
	}

	return c.JSON(http.StatusOK, response)
}
