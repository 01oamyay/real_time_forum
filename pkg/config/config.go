package config

import (
	"encoding/json"
	"os"
)

type (
	Conf struct {
		API      API      `json:"api"`
		Database Database `json:"database"`
	}

	API struct {
		Host string `json:"host"`
		Port string `json:"port"`
	}
	Database struct {
		Driver    string `json:"driver"`
		FileName  string `json:"fileName"`
		SchemeDir string `json:"schemeDir"`
	}
)

func NewConfig() (*Conf, error) {
	var newConfig Conf
	file, err := os.Open("./config/config.json")
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(file).Decode(&newConfig); err != nil {
		return nil, err
	}
	return &newConfig, nil
}
