package utils

import "github.com/gorilla/websocket"

func IndexOfWsClient(array []*websocket.Conn, element *websocket.Conn) uint16 {
	for i, v := range array {
		if v == element {
			return uint16(i)
		}
	}
	return 0
}

func IndexOfString(array []string, element string) int {
	for i, v := range array {
		if v == element {
			return i
		}
	}
	return 0
}
