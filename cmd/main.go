package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Mikeyteam/preview_project_go/internal/dispatcher"
	"github.com/Mikeyteam/preview_project_go/internal/http"
	"github.com/Mikeyteam/preview_project_go/internal/storage"

	"github.com/Mikeyteam/preview_project_go/internal/config"
	"github.com/Mikeyteam/preview_project_go/internal/logger"
)

func main() {
	conf := config.ConfigFromEnv()
	fmt.Printf("current settings: %+v\n", conf)
	log := logger.NewLogger(logger.Config{
		Level: conf.LogLevel,
	})
	st := storage.Create(conf.CacheType, conf.CachePath, log)
	imageDispatcher := dispatcher.New(st, conf.CacheSize, log)

	server := http.NewServer(http.Config{
		HTTPListen:       conf.HTTPListen,
		ImageMaxFileSize: conf.ImageMaxFileSize,
		ImageGetTimeout:  conf.ImageGetTimeout,
	}, log, &imageDispatcher)

	server.RunServer()
	defer server.Shutdown()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	log.Infof("got signal from OS: %v. Exit...", <-osSignals)
}
