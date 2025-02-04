package configs

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"CivicCoinMain/pkg/models"
)

var (
	cfg  *models.Configs
	once sync.Once
)

func LoadConfigs() (*models.Configs, error) {
	var err error
	once.Do(func() {
		file, fileErr := os.Open("configs/configs.json")
		if fileErr != nil {
			err = fmt.Errorf("failed to open configuration file: %w", fileErr)
			return
		}
		defer file.Close()

		cfg = &models.Configs{}
		if decodeErr := json.NewDecoder(file).Decode(cfg); decodeErr != nil {
			err = fmt.Errorf("could not decode the JSON file: %w", decodeErr)
		}
	})
	return cfg, err
}

// GetConfig devuelve la configuración cargada
func GetConfig() *models.Configs {
	return cfg
}
