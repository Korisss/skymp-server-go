package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	args := os.Args

	for _, arg := range args {
		switch arg {
		case "build":
			build()
		case "build-make":
			buildMake()
		}
	}
}

func buildMake() {
	cmd := exec.Command("go", "build", "-o", "./make.exe", "./tools/make/cmd/app/main.go")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func build() {
	cmd := exec.Command("go", "build", "-o", "./build/server.exe", "./cmd/app/main.go")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
