package service

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type QuilibriumConfig struct {
	PeerPrivKey string `yaml:"peerPrivKey"`
	// 可以根据需要添加其他配置字段
}

type QuilibriumKeys struct {
	// 根据实际的keys.yml结构定义字段
	Keys map[string]interface{} `yaml:",inline"`
}

func ParseQuilibriumConfig(configData []byte) (*QuilibriumConfig, error) {
	var config QuilibriumConfig
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config.yml: %w", err)
	}
	
	if config.PeerPrivKey == "" {
		return nil, fmt.Errorf("peerPrivKey is required in config.yml")
	}
	
	return &config, nil
}

func ParseQuilibriumKeys(keysData []byte) (*QuilibriumKeys, error) {
	var keys QuilibriumKeys
	if err := yaml.Unmarshal(keysData, &keys); err != nil {
		return nil, fmt.Errorf("failed to parse keys.yml: %w", err)
	}
	
	return &keys, nil
}

func ValidateConfigFiles(configData, keysData []byte) error {
	_, err := ParseQuilibriumConfig(configData)
	if err != nil {
		return fmt.Errorf("invalid config.yml: %w", err)
	}
	
	_, err = ParseQuilibriumKeys(keysData)
	if err != nil {
		return fmt.Errorf("invalid keys.yml: %w", err)
	}
	
	return nil
}