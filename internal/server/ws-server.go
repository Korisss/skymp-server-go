package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type WsServer struct {
	server     *Server
	httpServer *http.Server
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWsServer(server *Server) *WsServer {
	wsServer := &WsServer{
		server: server,
	}
	return wsServer
}

func (s *WsServer) onOpen(client *websocket.Conn) {
}

func (s *WsServer) onClose(client *websocket.Conn) {
	defer client.Close()
}

func (s *WsServer) onMessage(client *websocket.Conn, msg *JsonSocketMessage) {
	s.server.OnWsMessage(client, msg)
}

func (s *WsServer) handle(c *gin.Context) {
	client, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Println("error get connection")
		logrus.Fatal(err)
	}

	s.onOpen(client)

	defer s.onClose(client)

	for {
		var data JsonSocketMessage
		err = client.ReadJSON(&data)
		if err != nil {
			logrus.Println("error read json")
			logrus.Errorln(err.Error())
			break
		}
		fmt.Println(data)
		s.onMessage(client, &data)
	}
}

func (s *WsServer) Run(errChan chan error, port uint16) {
	router := gin.Default()
	router.GET("/", s.handle)

	s.httpServer = &http.Server{
		Addr:           ":" + strconv.Itoa(int(port+2)),
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	errChan <- s.httpServer.ListenAndServe()
}

func (s *WsServer) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
