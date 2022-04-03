package settings

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
)

type Settings struct {
	Ip                        string   `json:"ip"`
	Port                      uint16   `json:"port"`
	MaxPlayers                uint16   `json:"maxPlayers"`
	MasterUrl                 string   `json:"master"`
	Name                      string   `json:"name"`
	LoadOrder                 []string `json:"loadOrder"`
	DataDir                   string   `json:"dataDir"`
	IsPapyrusHotReloadEnabled bool     `json:"isPapyrusHotReloadEnabled"`
}

func Load(filePath string) *Settings {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	var settings Settings

	err = json.Unmarshal(fileContent, &settings)
	if err != nil {
		log.Fatal(err.Error())
	}

	if net.ParseIP(settings.Ip) == nil {
		settings.Ip = "127.0.0.1"
	}

	if settings.MaxPlayers == 0 {
		settings.MaxPlayers = 10
	}

	res, err := http.Get(settings.MasterUrl + "/api/servers")
	if err != nil || res.StatusCode != http.StatusOK {
		settings.MasterUrl = "https://skymp.io"
	}

	res.Body.Close()

	if settings.Name == "" {
		settings.Name = "My Server"
	}

	if len(settings.LoadOrder) == 0 {
		settings.LoadOrder = []string{
			"Skyrim.esm",
			"Update.esm",
			"Dawnguard.esm",
			"HearthFires.esm",
			"Dragonborn.esm",
		}
	}

	return &settings
}
