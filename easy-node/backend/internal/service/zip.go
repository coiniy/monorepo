package service

import (
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

type NodeConfigData struct {
	FolderName string
	ConfigYml  []byte
	KeysYml    []byte
}

func ExtractNodeConfigs(zipReader *zip.Reader) ([]NodeConfigData, error) {
	var nodeConfigs []NodeConfigData
	configMap := make(map[string]*NodeConfigData)

	for _, file := range zipReader.File {
		// 跳过目录和隐藏文件
		if file.FileInfo().IsDir() || strings.HasPrefix(filepath.Base(file.Name), ".") {
			continue
		}

		// 解析路径
		pathParts := strings.Split(file.Name, "/")
		if len(pathParts) < 2 {
			continue // 跳过根目录文件
		}

		folderName := pathParts[0]
		fileName := pathParts[len(pathParts)-1]

		// 只处理 config.yml 和 keys.yml
		if fileName != "config.yml" && fileName != "keys.yml" {
			continue
		}

		// 初始化节点配置数据
		if configMap[folderName] == nil {
			configMap[folderName] = &NodeConfigData{
				FolderName: folderName,
			}
		}

		// 读取文件内容
		rc, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", file.Name, err)
		}

		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", file.Name, err)
		}

		// 分配文件内容
		switch fileName {
		case "config.yml":
			configMap[folderName].ConfigYml = data
		case "keys.yml":
			configMap[folderName].KeysYml = data
		}
	}

	// 验证每个文件夹都包含必需的文件
	for folderName, config := range configMap {
		if len(config.ConfigYml) == 0 {
			return nil, fmt.Errorf("folder %s is missing config.yml", folderName)
		}
		if len(config.KeysYml) == 0 {
			return nil, fmt.Errorf("folder %s is missing keys.yml", folderName)
		}
		nodeConfigs = append(nodeConfigs, *config)
	}

	if len(nodeConfigs) == 0 {
		return nil, fmt.Errorf("no valid node configurations found in zip file")
	}

	return nodeConfigs, nil
}