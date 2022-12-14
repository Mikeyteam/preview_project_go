package http

import (
	"context"
	"net/http"
	"time"

	"github.com/Mikeyteam/preview_project_go/internal/dispatcher"
	"github.com/Mikeyteam/preview_project_go/internal/helpers"
	"github.com/sirupsen/logrus"
)

type Server struct {
	config          Config
	logger          *logrus.Logger
	httpServer      *http.Server
	imageDispatcher *dispatcher.ImageDispatcher
}

type Config struct {
	HTTPListen       string
	ImageMaxFileSize int
	ImageGetTimeout  int
}

func NewServer(config Config, logger *logrus.Logger, imageDispatcher *dispatcher.ImageDispatcher) Server {
	return Server{
		config:          config,
		logger:          logger,
		httpServer:      &http.Server{Addr: config.HTTPListen, Handler: handlerRequest(logger, imageDispatcher, config.ImageMaxFileSize, config.ImageGetTimeout)},
		imageDispatcher: imageDispatcher,
	}
}

func (s *Server) RunServer() {
	go func() {
		s.logger.Infoln("start HTTP server:", s.config.HTTPListen)
		err := s.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			helpers.FailOnError(err, "fail start HTTP Server")
		}
	}()
}

func (s *Server) Shutdown() {
	s.logger.Infoln("shutdown HTTP server... ")
	ctx, chancel := context.WithTimeout(context.Background(), time.Second*10)
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		s.logger.Errorln("fail Shutdown HTTP server")
		chancel()
		return
	}
	s.logger.Infoln("success Shutdown HTTP server")
	chancel()
}
