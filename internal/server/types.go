package server

type UserId = uint16
type ActorFormId = uint32
type ProfileId = int

//используется для кодирования send_server_info
type JsonServerInfo struct {
	Name       string `json:"name"`
	MaxPlayers int    `json:"maxPlayers"`
	Online     int    `json:"online"`
}

//используется в get_user_profile_id
type JsonResponseProfileId struct {
	User JsonUser `json:"user"`
}

//используется только в JsonResponseProfileId
type JsonUser struct {
	Id ProfileId `json:"id"` // = -1
}

//используется в on_custom_packet
type JsonClientVersion struct {
	Typ string `json:"customPacketType"`
	Src string `json:"src"`
}

//используется в on_ws_message
type JsonSocketMessage struct {
	Typ       string      `json:"type"`
	Token     string      `json:"token"`
	MsgRaw    interface{} `json:"msg"` //[raw]
	ChannelId int         `json:"channelIdx"`
	Text      string      `json:"text"`
}

//используется в on_ws_message
type JsonChatMessage struct {
	Typ        string `json:"type"`
	AuthorName string `json:"author"`
	ChannelId  int    `json:"channelIdx"`
	Text       string `json:"text"`
}

type JsonUiEvent struct {
	Typ  string `json:"type"`
	Data string `json:"data"`
}
