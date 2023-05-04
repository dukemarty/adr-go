package data

import (
	"encoding/json"
	"io/ioutil"
)

// {"language":"en","path":"docs/adr/","prefix":"abc","digits":3}

type Configuration struct {
	Language     string `json:"language"`
	Path         string `json:"path"`
	Prefix       string `json:"prefix"`
	Digits       int    `json:"digits"`
	TemplateName string `json:"template"`
}

func NewConfiguration(lang string, path string, prefix string, digits int) *Configuration {
	c := Configuration{
		Language: lang,
		Path:     path,
		Prefix:   prefix,
		Digits:   digits,
	}

	return &c
}

func (config Configuration) Store(filepath string) error {
	content, _ := json.MarshalIndent(config, "", " ")

	errorRes := ioutil.WriteFile(filepath, content, 0644)

	return errorRes
}
