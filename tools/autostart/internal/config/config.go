package config

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

const ConfigDir = "/etc/autostart-manager"

type ServiceConfig struct {
	Name         string            `json:"name"`
	ExecStart    string            `json:"exec_start"`    // 完整的启动命令
	WorkDir      string            `json:"work_dir"`      // 工作目录
	User         string            `json:"user"`          // 运行用户
	Group        string            `json:"group"`         // 运行组
	Description  string            `json:"description"`   // 服务描述
	Env          map[string]string `json:"env"`           // 环境变量
	Restart      string            `json:"restart"`       // 重启策略
	RestartSec   int               `json:"restart_sec"`   // 重启间隔
	KillMode     string            `json:"kill_mode"`     // 终止模式
	KillSignal   string            `json:"kill_signal"`   // 终止信号
	TimeoutStart int               `json:"timeout_start"` // 启动超时
	TimeoutStop  int               `json:"timeout_stop"`  // 停止超时
	After        []string          `json:"after"`         // 依赖服务
	Wants        []string          `json:"wants"`         // 期望服务
	Requires     []string          `json:"requires"`      // 必需服务
}

// 保存服务配置
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
