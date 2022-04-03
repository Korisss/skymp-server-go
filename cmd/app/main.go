package main

import skymp_wrapper "github.com/Korisss/skymp-server-go/internal/skymp-wrapper"

func main() {
	skymp_wrapper.CreateServer()
	skymp_wrapper.Free()
}
