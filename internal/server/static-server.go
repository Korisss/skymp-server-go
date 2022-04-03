package server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type StaticServer struct {
	server *http.Server
}

func NewStaticServer() *StaticServer {
	return &StaticServer{}
}

func (s *StaticServer) Run(errChan chan error, port uint16, dataDir string) {
	router := gin.Default()
	router.Static("/", "./"+dataDir)

	s.server = &http.Server{
		Addr:           ":" + strconv.Itoa(int(port+1)),
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	errChan <- s.server.ListenAndServe()
}

func (s *StaticServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
