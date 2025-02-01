package configs

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Configs struct {
	Keys struct {
		Me    string   `json:"me"`
		Nodes []string `json:"nodes"`
	} `json:"keys"`
	Websocket struct {
		Address string `json:"address"`
	} `json:"websocket"`
}

var (
	cfg  *Configs
	once sync.Once
)

func LoadConfigs() (*Configs, error) {
	var err error
	once.Do(func() {
		file, fileErr := os.Open("configs/configs.json")
		if fileErr != nil {
			err = fmt.Errorf("Failed to open configuration file: %w", fileErr)
			return
		}
		defer file.Close()

		cfg = &Configs{}
		if decodeErr := json.NewDecoder(file).Decode(cfg); decodeErr != nil {
			err = fmt.Errorf("Could not decode the JSON file: %w", decodeErr)
		}
	})
	return cfg, err
}

// GetConfig devuelve la configuración cargada
func GetConfig() *Configs {
	return cfg
}
