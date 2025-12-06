package config

import (
	"encoding/json"
	"os"
)

type (
    // Conf describes API and Database configuration loaded from disk.
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

// NewConfig reads config/config.json and populates the Conf struct.
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
