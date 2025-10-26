package service

import (
	"fmt"
	"strings"
	"gopkg.in/yaml.v3"
)

type QuilibriumConfig struct {
	P2P struct {
		PeerPrivKey string `yaml:"peerPrivKey"`
	} `yaml:"p2p"`
}

type QuilibriumKeys struct {
	// 根据实际的keys.yml结构定义字段
	Keys map[string]interface{} `yaml:",inline"`
}

// ParseQuilibriumConfig 使用宽松模式解析配置文件
// 即使存在重复的键也能正常工作，因为我们只关心特定字段
func ParseQuilibriumConfig(configData []byte) (*QuilibriumConfig, error) {
	// 首先尝试标准解析
	var config QuilibriumConfig
	err := yaml.Unmarshal(configData, &config)

	// 如果解析失败，尝试使用 map 方式提取 peerPrivKey
	if err != nil {
		// 检查是否是重复键错误
		if strings.Contains(err.Error(), "already defined") {
			// 使用 map 解析，重复的键会被后面的值覆盖
			var configMap map[string]interface{}
			decoder := yaml.NewDecoder(strings.NewReader(string(configData)))
			decoder.KnownFields(false) // 允许未知字段

			if mapErr := decoder.Decode(&configMap); mapErr != nil {
				return nil, fmt.Errorf("failed to parse config.yml: %w", err)
			}

			// 从 map 中提取 p2p.peerPrivKey
			if p2pMap, ok := configMap["p2p"].(map[string]interface{}); ok {
				if peerKey, ok := p2pMap["peerPrivKey"].(string); ok {
					config.P2P.PeerPrivKey = peerKey
				} else {
					return nil, fmt.Errorf("p2p.peerPrivKey not found or invalid in config.yml")
				}
			} else {
				return nil, fmt.Errorf("p2p section not found in config.yml")
			}
		} else {
			return nil, fmt.Errorf("failed to parse config.yml: %w", err)
		}
	}

	if config.P2P.PeerPrivKey == "" {
		return nil, fmt.Errorf("peerPrivKey is required in config.yml")
	}

	return &config, nil
}

// GetPeerPrivKey 返回节点的 peer private key
func (c *QuilibriumConfig) GetPeerPrivKey() string {
	return c.P2P.PeerPrivKey
}

func ParseQuilibriumKeys(keysData []byte) (*QuilibriumKeys, error) {
	var keys QuilibriumKeys
	// 使用宽松模式
	decoder := yaml.NewDecoder(strings.NewReader(string(keysData)))
	decoder.KnownFields(false)

	if err := decoder.Decode(&keys); err != nil {
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
