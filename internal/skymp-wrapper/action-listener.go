package skymp_wrapper

import (
	"encoding/json"
	"fmt"
)

// const (
// 	movement_message_header_byte = 'M'
// )

type Packet interface{}

var PacketTypes = map[string]uint8{
	"invalid":               0,
	"custom_packet":         1,
	"update_movement":       2,
	"update_animation":      3,
	"update_appearance":     4,
	"update_equipment":      5,
	"activate":              6,
	"update_property":       7,
	"put_item":              8,
	"take_item":             9,
	"finish_sp_snippet":     10,
	"on_equip":              11,
	"console_command":       12,
	"craft_item":            13,
	"host":                  14,
	"custom_event":          15,
	"change_values":         16,
	"on_hit":                17,
	"death_state_container": 18,
	"max":                   19,
}

type UpdateMovement struct {
	Idx              uint32
	WorldOrCell      uint32
	Pos              Position // = Position{}
	Rot              Position // = Position{}
	Direction        float32
	HealthPercentage float32

	RunMode       uint8 //RunMode // = .standing
	IsInJumpState bool
	IsSneaking    bool
	IsBlocking    bool
	IsWeaponDrawn bool
	IsDead        bool
	LookAt        Position // = Position{}
}

type CustomPacket struct {
	Typ      string       `json:"customPacketType"`
	Src      string       `json:"src"`
	Token    string       `json:"token"`
	GameData JsonGameData `json:"gameData"`
}

type UpdateAnimation struct {
	NumChanges uint32 `json:"numChanges"`
	Name       string `json:"animEventName"`
}

type UpdateEquipment struct {
	Inv        Inventory `json:"inv"`
	NumChanges uint32    `json:"numChanges"`
}

type ChangeValuesPacket struct {
	Health  float32 `json:"health"`
	Stamina float32 `json:"stamina"`
	Magicka float32 `json:"magicka"`
}

type ActivatePacket struct {
	Caster FormId `json:"caster"`
	Target FormId `json:"target"`
}

type EquipPacket struct {
	BaseId FormId `json:"base_id"`
}

type HitPacket struct {
	Aggressor     FormId `json:"aggressor"`
	IsBashAttack  bool   `json:"isBashAttack"`
	IsHitBlocked  bool   `json:"isHitBlocked"`
	IsPowerAttack bool   `json:"isPowerAttack"`
	IsSneakAttack bool   `json:"isSneakAttack"`
	Projectile    int64  `json:"projectile"` //TODO wtf??????????
	Source        FormId `json:"source"`
	Target        FormId `json:"target"`
}

type PutItemPacket struct {
	Target FormId `json:"target"`
	BaseId FormId `json:"baseId"`
	Count  uint32 `json:"count"`
}

type TakeItemPacket struct {
	Target FormId `json:"target"`
	BaseId FormId `json:"baseId"`
	Count  uint32 `json:"count"`
}

type ConsoleCommandPacketArg = string // | int64

type ConsoleCommandPacket struct {
	CommandName string                    `json:"commandName"`
	Args        []ConsoleCommandPacketArg `json:"args"`
}

type UpdateAppearanceTint struct {
	TexturePath string `json:"texturePath"`
	Typ         uint32 `json:"type"`
	Argb        int    `json:"argb"`
}

type UpdateAppearance struct {
	IsFemale         bool                   `json:"isFemale"`
	RaceId           uint32                 `json:"raceId"`
	Weight           float32                `json:"weight"`
	HairColor        uint32                 `json:"hairColor"` //uint32?
	HeadpartIds      []uint32               `json:"headpartIds"`
	HeadTextureSetId uint32                 `json:"headTextureSetId"`
	Options          []float64              `json:"options"`
	Presets          []int                  `json:"presets"`
	Tints            []UpdateAppearanceTint `json:"tints"`
	SkinColor        uint32                 `json:"skinColor"` //uint32?
	Name             string                 `json:"name"`
}

type JsonPacket struct {
	Typ     PacketType `json:"t"`
	Idx     uint32     `json:"idx"`     // .update_animation .update_appearance .update_equipment
	Data    string     `json:"data"`    //[raw] // typ != .custom_packet
	Content string     `json:"content"` //[raw] // .custom_packet
	BaseId  FormId     `json:"baseId"`  // .on_equip
}

func ActionListener(userId UserId, data *Packet) {
	switch (*data).(type) {
	case UpdateMovement:
	case CustomPacket:
		fmt.Println("[Packet] Custom Packet")
	case UpdateAnimation:
		fmt.Println("[Packet] Update Animation Event", data)
	case UpdateEquipment:
		// Могут попадаться локальные предметы
		fmt.Println("[Packet] Update Equipment Event", data)
	case ChangeValuesPacket:
		fmt.Println("[Packet] Change Values Event", data)
	case ActivatePacket:
		fmt.Println("[Packet] ActivatePacket Event", data)
		println("[Packet] ActivatePacket Event ${data as ActivatePacket}")
		// ac := svr.actor_by_user(user_id) or {
		// 	eprintln('Can`t do this without Actor attached')
		// 	return
		// }
		// if data.caster != 0x14 {
		// 	hoster_id := svr.find_hoster(data.caster) or {
		// 		return
		// 	}
		// 	if hoster_id != ac.get_form_id() {
		// 		return
		// 	}
		// }
		// caster := ac.cast<ObjectReference>() or {
		// 	println('onActivate - caster not found')
		// 	return
		// }
		// target := svr.get_form<ObjectReference>(data.target) or {
		// 	println('onActivate - target not found')
		// 	return
		// }
		// svr.js_modules.send_event(GmActivateEventData{
		// 	target: target
		// 	caster: caster
		// })
		// svr.send_game_packet(caster, true, GPScale{
		// 	target: caster.get_form_id()
		// 	value: 0.5
		// })
	case EquipPacket:
		fmt.Println("[Packet] On Equip Event: ", data)
		// ac := svr.actor_by_user(user_id) or {
		// 	eprintln('Can`t do this without Actor attached')
		// 	return
		// }
		// //TODO по базовому id нельзя получить форму
		// target := svr.get_form<Form>(data.base_id) or {
		// 	println('onEquip - target not found')
		// 	return
		// }
		// svr.js_modules.send_event(GmEquipEventData{
		// 	actor: ac
		// 	target: target
		// })
	case HitPacket:
		fmt.Println("[Packet] One HitPacket Event", data)
	case TakeItemPacket:
		fmt.Println("[Packet] Take Item Event", data)
	case PutItemPacket:
		fmt.Println("[Packet] Put Item Event", data)
	case ConsoleCommandPacket:
		fmt.Println("[Packet] Console Command", data)
	case UpdateAppearance:
		fmt.Println("[Packet] Update Appearance", data)
	default:
		fmt.Println("[Packet] Unsupported type")
	}
}

func ParsePacketFromJson(jsonData string) (*Packet, error) {
	var packet Packet

	err := json.Unmarshal([]byte(jsonData), &packet)

	if err != nil {
		return nil, err
	}

	return &packet, nil
}

// [inline]
// fn read_movement_from_bitstream(data []byte) ?&UpdateMovement {
// 	// 2      6       10                22                          34                38      42
// 	//     ___idx__ _w_o_c__ ___________pos____________ ___________rot____________ ___dir__ __hp____
// 	//864d 00000000 0000003c 46bdd72d c71aca51 41d980f2 40b7949e 00000000 43a4deaa 00000000 3f800000 0000 6f006d0000000000000000004abbf82600db
// 	//864d 00000000 0000003c 46bdb67f c71aada9 41e2c200 40eefa7b 00000000 43a5ba3f 00000000 3f800000 1000 6f006d000000000000000000d2ba7020003d //is_sneaking = true
// 	//864d 00000000 0000003c 46bdb67f c71aada9 41e2c200 40b5ff71 00000000 43a5ef15 00000000 3f800000 0000 00000000000000000000000066baec210020 //is_sneaking = false
// 	//864d 00000000 0000003c 46bdb635 c71aadbe 41e28e91 40c113d3 00000000 43a5256c 00000000 3f800000 0000 000000000000000000000000dc187b060041 //is_sneaking = false
// 	//864d 00000000 0000003c 46bdb635 c71aadbe 41e28eb0 4089adfc 00000000 43a55ebd 00000000 3f800000 1000 6f006d000000000000000000fc185b060049 //is_sneaking = true
// 	//864d 00000000 0000003c 46be3c94 c71abb69 41f8c51e 41303c75 00000000 43a36019 00000000 3f800000 2000 00000000000087e670020000919172b700b7 //is_in_jump_state = true
// 	//864d 00000000 0000003c 46be3c94 c71abb69 41f8c519 41303c75 00000000 43a36019 00000000 3f800000 3000 0000000000000000000000005991bab70089 //is_in_jump_state = true, is_sneaking = true
// 	mut msg := UpdateMovement{}

// 	mut b := bitstream.create(data.data, usize(data.len))

// 	b.read<uint32>(mut msg.idx) ?
// 	b.read<uint32>(mut msg.world_or_cell) ?
// 	b.read<Position>(mut msg.pos) ?
// 	b.read<Position>(mut msg.rot) ?
// 	b.read<float32>(mut msg.direction) ?
// 	b.read<float32>(mut msg.health_percentage) ?

// 	mut val := u8(0)
// 	val |= byte(b.read_bool() ?)
// 	val <<= 1
// 	val |= byte(b.read_bool() ?)
// 	msg.run_mode = RunMode(int(val))

// 	msg.is_in_jump_state = b.read_bool() ?
// 	msg.is_sneaking = b.read_bool()  ?
// 	msg.is_blocking = b.read_bool()  ?
// 	msg.is_weap_drawn = b.read_bool()  ?
// 	msg.is_dead = b.read_bool()  ?
// 	b.read_optional<Position>(mut msg.look_at) or { eprintln(err.msg) }

// 	return &msg
// }

// fn transform_packet_into_action(user_id UserId, dataptr byteptr, len usize) ? {
// 	bytes := unsafe { dataptr.vbytes(int(len)) }

// 	if len > 1 {
// 		header_byte := bytes[1].ascii_str()

// 		if header_byte == movement_message_header_byte {
// 			data := bytes[2..len]
// 			msg := read_movement_from_bitstream(data) or {
// 				return error(err.msg)
// 			}

// 			svr.action_listener(user_id, msg)
// 			return
// 		}
// 	}

// 	json_data := unsafe { byteptr(dataptr + 1).vstring_literal_with_len(int(len - 1)) }

// 	data := parse_packet_from_json(json_data) ?
// 	svr.action_listener(user_id, data)
// }
