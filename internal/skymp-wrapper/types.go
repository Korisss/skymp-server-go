package skymp_wrapper

type OnConnectHandler = func(user_id UserId)
type OnDisonnectHandler = func(user_id UserId)
type OnCustomPacketHandler = func(user_id UserId, json_data uintptr)
type PacketHandlerFn = func(userId UserId, data uintptr, len uint64)

type MpForm = uintptr
type MpObjectReference = uintptr
type MpActor = uintptr
type FormDesc = uintptr

type ProfileId = int
type UserId = uint16
type ActorFormId = uint32
type FormId = uint32
type Angle = float32
type WorldOrCell = uint32

type Position struct {
	X float32
	Y float32
	Z float32
}

/*
fn (uid UserId) profile_id() ProfileId {
	return svr.profile_ids[uid.u16()]
}
fn (uid UserId) token() string {
	return svr.tokens[uid.u16()]
}*/

// func (uid UserId) ws_client() ?&lws.Client {
// 	client_id := svr.ws_client_ids[uid.u16()]
// 	client := svr.ui_ws.clients[client_id] or { return error("websocket client not found") }
// 	return client
// }

type InvEntry struct {
	Base_id FormId `json:"baseId"`
	Count   uint32 `json:"count"`
	Name    string `json:"name"`
	Worn    bool   `json:"worn"`
}

type Inventory struct {
	Entries []InvEntry
}

var RunMode = map[string]uint8{
	"standing":  0,
	"walking":   1,
	"running":   2,
	"sprinting": 3,
}

type PacketType uint8

type GamePacket struct {
	Typ PacketType `json:"t"` // = .custom_packet
	//typ2	string			[json: 'type'] = 'CustomPacket' //or 'customPacket'
	Data GamePacketData `json:"content"`
}

type GamePacketData = GPScale //| Empty

type GPScale struct {
	Typ    string  `json:"type"` // = 'function'
	Name   string  `json:"name"` // = "setScale"
	Target FormId  `json:"target"`
	Value  float32 `json:"value"`
}

// struct Empty {}
