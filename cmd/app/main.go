package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Korisss/skymp-server-go/internal/manifest"
	serverpkg "github.com/Korisss/skymp-server-go/internal/server"
	"github.com/Korisss/skymp-server-go/internal/settings"
	skymp_wrapper "github.com/Korisss/skymp-server-go/internal/skymp-wrapper"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	settings := settings.Load("./server-settings.json")
	server := serverpkg.NewServer(settings)

	defer skymp_wrapper.Free()

	go manifest.GenerateManifest(settings.DataDir, settings.LoadOrder)

	go func() {
		if err := server.Run(); !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("Error occurred while running http server: %s", err.Error())
		}
	}()

	logrus.Print("SkyMP Go server started...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("SkyMP Go server shutting down...")

	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occurred on server shutting down: %s", err.Error())
	}

	logrus.Info("SkyMP Go server closed properly")
}
