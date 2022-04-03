package server

import "C"
import (
	"encoding/json"
	"unsafe"

	skymp_wrapper "github.com/Korisss/skymp-server-go/internal/skymp-wrapper"
	"github.com/Korisss/skymp-server-go/pkg/utils"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func (s *Server) OnWsMessage(client *websocket.Conn, msg *JsonSocketMessage) {
	logrus.Println("New ws message:", msg)

	switch msg.Typ {
	case "token":
		logrus.Println("Try token", msg.Token)
		if msg.Token == "" {
			logrus.Errorln("Invalid token")
			return
		}

		userId := s.getUserIdFromToken(msg.Token)
		s.wsClients[userId] = client
	case "uiEvent":
		userId := utils.IndexOfWsClient(s.wsClients, client)
		actorId, err := s.scampServer.GetUserActor(userId)
		if err != nil {
			logrus.Errorln(err.Error())
			return
		}

		s.onUiEvent(actorId, msg.Msg)
	}
}

// На момент подключения актёр ещё не доступен
func (s *Server) onConnect(userId UserId) {}

// Срабатывает не сразу
func (s *Server) onDisconnect(userId UserId) {
	s.profileIds[userId] = -1

	actorId, err := s.scampServer.GetUserActor(userId)
	if err != nil {
		return
	}

	s.scampServer.SetEnabled(actorId, false)
}

func (s *Server) onCustomPacket(userId UserId, jsonData uintptr) {
	jsonStr := C.GoString((*C.char)(unsafe.Pointer(jsonData)))
	var data skymp_wrapper.CustomPacket
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		logrus.Errorln("[Error onCustomPacket] failed to decode json")
		logrus.Errorln(jsonData)
		logrus.Errorln(err.Error())
		return
	}

	switch data.Typ {
	case "browserToken":
		logrus.Println("User", userId, "sets browser token to", data.Token)
		s.tokens[userId] = data.Token
	case "loginWithSkympIo":
		if data.GameData.Session == "" {
			logrus.Errorln("Error on login")
			return
		}

		profileId := s.getUserProfileId(data.GameData.Session)
		if profileId < 0 {
			logrus.Errorln(err)
			return
		}

		s.profileIds[userId] = profileId

		s.onSpawnAllowed(userId, profileId)

		logrus.Println("Logged as", profileId)
	default:
		logrus.Errorln("invalid json data")
	}
}

func (s *Server) onSpawnAllowed(userId UserId, profileId ProfileId) {
	actorIds := s.scampServer.GetActorsByProfileId(profileId)

	if len(actorIds) > 0 {
		actorId := actorIds[0]

		logrus.Println("Loading character", actorId)
		s.scampServer.SetEnabled(actorId, true)
		s.scampServer.SetUserActor(userId, actorId)
	} else {
		pos := skymp_wrapper.Position{
			X: 22106.24609375,
			Y: -44752.68359375,
			Z: -140.59170532226562,
		}

		actorId := s.scampServer.CreateActor(0, pos, 47, 0x3c, profileId)
		s.scampServer.SetUserActor(userId, actorId)
		s.scampServer.SetRaceMenuOpen(actorId, true)
	}
}

func (s *Server) onUiEvent(formId ActorFormId, jsonData string) {
	var data JsonUiEvent
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		logrus.Errorln("error read json")
		logrus.Errorln(err.Error())
		return
	}

	switch data.Typ {
	case "cef::chat:send":
		logrus.Println("onUiEvent", data)
		//	tokens := strings.Split(data.Data, " ")

		// if strings.HasPrefix(tokens[0], "/") {
		//	TODO call module.onChatCommand
		//	D:\_projects\RH-workspace\functions-lib\src\events\index.ts#266
		// } else {
		//	TODO call module.onChatInput
		//	D:\_projects\RH-workspace\functions-lib\src\events\index.ts#282
		// }
	case "server::msg:send":
		logrus.Println("onUiEvent", data)
		//TODO call handleServerMsg
		//D:\_projects\RH-workspace\functions-lib\src\events\server-msg.ts#22
	case "socketOpen":
		// видимо не используется
	case "error":
		logrus.Errorln("Error onUiEvent", data)
	default:
		logrus.Errorln("invalid json data")
	}
}
