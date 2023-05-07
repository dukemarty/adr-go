package data

import (
	"encoding/json"
	"os"
)

// {"language":"en","path":"docs/adr/","prefix":"abc","digits":3}

type Configuration struct {
	Language     string `json:"language"`
	Path         string `json:"path"`
	Prefix       string `json:"prefix"`
	Digits       int    `json:"digits"`
	TemplateName string `json:"template"`
}

func NewConfiguration(lang string, path string, prefix string, digits int, template string) *Configuration {
	c := Configuration{
		Language:     lang,
		Path:         path,
		Prefix:       prefix,
		Digits:       digits,
		TemplateName: template,
	}

	return &c
}

func LoadConfiguration() (Configuration, error) {
	var config Configuration

	content, err := os.ReadFile(".adr.json")
	if err != nil {
		return config, err
	}
	parseErr := json.Unmarshal(content, &config)
	if parseErr != nil {
		return config, err
	}

	return config, nil
}

func (config Configuration) Store(filepath string) error {
	content, _ := json.MarshalIndent(config, "", " ")

	errorRes := os.WriteFile(filepath, content, 0644)

	return errorRes
}
