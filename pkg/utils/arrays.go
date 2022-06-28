package utils

import (
	"fmt"

	"github.com/gorilla/websocket"
)

func IndexOfWsClient(array []*websocket.Conn, element *websocket.Conn) uint16 {
	for i, v := range array {
		if v == element {
			fmt.Println(i)
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
