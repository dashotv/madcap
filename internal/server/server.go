package server

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
	"github.com/streamingfast/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/term"

	"github.com/dashotv/madcap/internal/plex"
)

type Server struct {
	Config *Config

	Router *echo.Echo
	Cron   *cron.Cron
	Logger *zap.SugaredLogger
	Plex   *plex.Client

	db *connection

	// Services
	page *fileService
}

func New() (*Server, error) {
	logger := setupLogger()

	s := &Server{
		Logger: logger,
	}

	if err := setupConfig(s); err != nil {
		return nil, err
	}
	if err := setupDatabase(s); err != nil {
		return nil, err
	}
	if err := setupPlex(s); err != nil {
		return nil, err
	}
	setupRouter(s)
	setupCron(s)

	file := &fileService{log: logger.Named("services.file"), server: s}

	g := s.Router.Group("/api")
	RegisterFileService(g, file)

	return s, nil
}

func (s *Server) Start() error {
	if s.Cron != nil {
		go s.Cron.Run()
	}

	return s.Router.Start(":" + s.Config.Port)
}

func setupLogger() *zap.SugaredLogger {
	isTTY := term.IsTerminal(int(os.Stderr.Fd()))
	verbosity := 1
	logStdoutWriter := zapcore.Lock(os.Stderr)
	log := zap.New(zapcore.NewCore(logging.NewEncoder(verbosity, isTTY), logStdoutWriter, zapcore.DebugLevel))
	return log.Named("madcap").Sugar()
}
