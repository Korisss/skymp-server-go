package manifest

import (
	"encoding/json"
	"hash/crc32"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type ManifestMod struct {
	Crc32    int32  `json:"crc32"`
	Filename string `json:"filename"`
	Size     int    `json:"size"`
}

type Manifest struct {
	VersionMajor int           `json:"versionMajor"`
	LoadOrder    []string      `json:"loadOrder"`
	Mods         []ManifestMod `json:"mods"`
}

func GenerateManifest(dataDir string, loadOrder []string) {
	manifest := Manifest{
		VersionMajor: 1,
		LoadOrder:    loadOrder,
		Mods:         []ManifestMod{},
	}

	for _, espm := range loadOrder {
		if !strings.HasSuffix(espm, ".esp") && !strings.HasSuffix(espm, ".esm") && !strings.HasSuffix(espm, ".esl") {
			logrus.Fatal(espm, "is not valid plugin name")
		}

		espmPath := dataDir + "/" + espm

		manifest.Mods = append(manifest.Mods, generateManifestMod(espm, espmPath))

		bsaName := getBsaName(espm)
		bsaPath := dataDir + "/" + bsaName

		if _, err := os.Open(bsaPath); err == nil {
			manifest.Mods = append(manifest.Mods, generateManifestMod(bsaName, bsaPath))
		}
	}

	jsonBuf, _ := json.Marshal(manifest)

	f, err := os.Create(dataDir + "/manifest.json")
	if err != nil {
		logrus.Errorln(err.Error())
	}

	defer f.Close()

	f.Write(jsonBuf)
}

func generateManifestMod(espmName, espmPath string) ManifestMod {
	buffer, err := os.ReadFile(espmPath)
	if err != nil {
		logrus.Fatalln(err.Error())
	}

	return ManifestMod{
		Crc32:    int32(crc32.ChecksumIEEE(buffer)),
		Filename: espmName,
		Size:     len(buffer),
	}
}

func getBsaName(espmName string) string {
	return espmName[0:len(espmName)-3] + "bsa"
}
