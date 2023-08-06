package data

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

const UserConfigFilename = ".adr-go"

// {"editor":"code", "centralstore":""}

type UserConfiguration struct {
	Editor          string `json:"editor"`
	CentralAdrStore string `json:"centralstore"`
}

func NewUserConfiguration(editor string, store string) *UserConfiguration {
	uc := UserConfiguration{
		Editor:          editor,
		CentralAdrStore: store,
	}

	return &uc
}

func LoadUserConfiguration() (UserConfiguration, error) {
	var config UserConfiguration

	home, err := os.UserHomeDir()
	if err != nil {
		return config, err
	}
	configPath := filepath.Join(home, UserConfigFilename)
	content, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}
	parseErr := json.Unmarshal(content, &config)
	if parseErr != nil {
		return config, err
	}

	return config, nil
}

func (config UserConfiguration) Store() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(home, UserConfigFilename)

	content, _ := json.MarshalIndent(config, "", "")

	errorRes := os.WriteFile(configPath, content, 0644)

	return errorRes
}

func LoadEditor(logger *log.Logger) string {
	config, err := LoadUserConfiguration()
	if err != nil {
		logger.Printf("Error loading user configuration: %v\n", err)
		return ""
	}

	return config.Editor
}
