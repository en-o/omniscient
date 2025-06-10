package config

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

const ConfigDir = "/etc/autostart-manager"

type ServiceConfig struct {
	Name         string            `json:"name"`
	ExecStart    string            `json:"exec_start"`
	WorkDir      string            `json:"work_dir"`
	User         string            `json:"user"`
	Group        string            `json:"group"`
	Description  string            `json:"description"`
	Env          map[string]string `json:"env"`
	Restart      string            `json:"restart"`
	RestartSec   int               `json:"restart_sec"`
	KillMode     string            `json:"kill_mode"`
	KillSignal   string            `json:"kill_signal"`
	TimeoutStart int               `json:"timeout_start"`
	TimeoutStop  int               `json:"timeout_stop"`
	After        []string          `json:"after"`
	Wants        []string          `json:"wants"`
	Requires     []string          `json:"requires"`
}

func SaveServiceConfig(config *ServiceConfig) error {
	configFile := filepath.Join(ConfigDir, config.Name+".json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configFile, data, 0644)
}

func LoadServiceConfig(name string) (*ServiceConfig, error) {
	configFile := filepath.Join(ConfigDir, name+".json")
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var config ServiceConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
