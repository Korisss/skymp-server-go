package server

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Korisss/skymp-server-go/internal/settings"
	skymp_wrapper "github.com/Korisss/skymp-server-go/internal/skymp-wrapper"
	"github.com/Korisss/skymp-server-go/pkg/utils"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Server struct {
	settings     *settings.Settings
	staticServer *StaticServer
	wsServer     *WsServer
	wsClients    []*websocket.Conn
	profileIds   []ProfileId
	tokens       []string
	scampServer  skymp_wrapper.ScampServer
	endpoint     string
}

func NewServer(settings *settings.Settings) *Server {
	server := &Server{
		settings:     settings,
		staticServer: NewStaticServer(),
		scampServer:  skymp_wrapper.CreateServer(settings.Port, settings.MaxPlayers),
		endpoint:     settings.MasterUrl + "/api/servers/" + settings.Ip + ":" + strconv.Itoa(int(settings.Port)),
		wsClients:    make([]*websocket.Conn, settings.MaxPlayers),
		profileIds:   make([]ProfileId, settings.MaxPlayers),
		tokens:       make([]string, settings.MaxPlayers),
	}

	server.wsServer = NewWsServer(server)

	return server
}

func (s *Server) Run() error {
	s.scampServer.SetConnectHandler(s.onConnect)
	s.scampServer.SetDisconnectHandler(s.onDisconnect)
	s.scampServer.SetCustomPacketHandler(s.onCustomPacket)

	// svr.set_packet_handler(
	// 	fn(user_id UserId, data byteptr, len usize) {
	// 		transform_packet_into_action(user_id, data, len) or {
	// 			eprintln(err)
	// 			return
	// 		}
	// 	}
	// )

	errChan := make(chan error, 1)

	go s.staticServer.Run(errChan, s.settings.Port, s.settings.DataDir)
	go s.wsServer.Run(errChan, s.settings.Port)

	err := <-errChan

	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	err1 := s.staticServer.Shutdown(ctx)
	err2 := s.wsServer.Shutdown(ctx)

	if err1 != nil {
		return err1
	}

	return err2
}

func (s *Server) getUserActor(userId UserId) (ActorFormId, error) {
	return s.scampServer.GetUserActor(userId)
}

func (s *Server) getUserIdFromToken(token string) UserId {
	return uint16(utils.IndexOfString(s.tokens, token))
}

func (s *Server) getUserProfileId(session string) ProfileId {
	url := s.endpoint + "/sessions/" + session
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != http.StatusOK {
		logrus.Errorln("Error when get id from session")
		return -1
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorln("Error when get id from session")
		return -1
	}

	var data JsonResponseProfileId
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		logrus.Errorln("Error when get id from session")
		return -1
	}

	return data.User.Id
}
