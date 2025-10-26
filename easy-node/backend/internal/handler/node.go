package handler

import (
	"archive/zip"
	"bytes"
	"easy-node-backend/internal/model"
	"easy-node-backend/internal/service"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type NodeHandler struct {
	nodeService         *service.NodeService
	dockerService       *service.DockerService
	dockerClientService *service.DockerClientService
	healthService       *service.HealthService
}

func NewNodeHandler(nodeService *service.NodeService, dockerService *service.DockerService, dockerClientService *service.DockerClientService, healthService *service.HealthService) *NodeHandler {
	return &NodeHandler{
		nodeService:         nodeService,
		dockerService:       dockerService,
		dockerClientService: dockerClientService,
		healthService:       healthService,
	}
}

func (h *NodeHandler) GetNodes(c *gin.Context) {
	nodes, err := h.nodeService.GetAllNodes()
	if err != nil {
		zap.L().Error("Failed to get nodes", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nodes)
}

func (h *NodeHandler) CreateNode(c *gin.Context) {
	var node model.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.nodeService.CreateNode(&node); err != nil {
		zap.L().Error("Failed to create node", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.dockerService.GenerateDockerCompose(); err != nil {
		zap.L().Error("Failed to generate docker-compose", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Node created but failed to update docker-compose"})
		return
	}

	c.JSON(http.StatusCreated, node)
}

func (h *NodeHandler) UpdateNode(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node ID"})
		return
	}

	var node model.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	node.ID = uint(id)
	if err := h.nodeService.UpdateNode(&node); err != nil {
		zap.L().Error("Failed to update node", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.dockerService.GenerateDockerCompose(); err != nil {
		zap.L().Error("Failed to generate docker-compose", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Node updated but failed to update docker-compose"})
		return
	}

	c.JSON(http.StatusOK, node)
}

func (h *NodeHandler) DeleteNode(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node ID"})
		return
	}

	if err := h.nodeService.DeleteNode(uint(id)); err != nil {
		zap.L().Error("Failed to delete node", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.dockerService.GenerateDockerCompose(); err != nil {
		zap.L().Error("Failed to generate docker-compose", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Node deleted but failed to update docker-compose"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Node deleted successfully"})
}

func (h *NodeHandler) RegenerateDockerCompose(c *gin.Context) {
	if err := h.dockerService.GenerateDockerCompose(); err != nil {
		zap.L().Error("Failed to generate docker-compose", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Docker compose regenerated successfully"})
}

func (h *NodeHandler) CreateNodeFromConfig(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	configFiles := form.File["config"]
	keysFiles := form.File["keys"]
	
	if len(configFiles) == 0 || len(keysFiles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both config.yml and keys.yml files are required"})
		return
	}

	configFile := configFiles[0]
	keysFile := keysFiles[0]

	configData, err := h.readUploadedFile(configFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read config file"})
		return
	}

	keysData, err := h.readUploadedFile(keysFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read keys file"})
		return
	}

	if err := service.ValidateConfigFiles(configData, keysData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config, err := service.ParseQuilibriumConfig(configData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var nodeRequest struct {
		Name     string `json:"name" binding:"required"`
		BasePort int    `json:"base_port" binding:"required"`
	}

	if err := c.ShouldBind(&nodeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	node := &model.Node{
		Name:      nodeRequest.Name,
		BasePort:  nodeRequest.BasePort,
		PeerKey:   config.GetPeerPrivKey(),
		ConfigYml: string(configData),
		KeysYml:   string(keysData),
		Enabled:   true,
	}

	if err := h.nodeService.CreateNodeFromConfig(node); err != nil {
		zap.L().Error("Failed to create node from config", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.dockerService.GenerateDockerCompose(); err != nil {
		zap.L().Error("Failed to generate docker-compose", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Node created but failed to update docker-compose"})
		return
	}

	c.JSON(http.StatusCreated, node)
}

func (h *NodeHandler) readUploadedFile(header *multipart.FileHeader) ([]byte, error) {
	file, err := header.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

func (h *NodeHandler) BatchUploadNodes(c *gin.Context) {
	file, _, err := c.Request.FormFile("zip")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get zip file"})
		return
	}
	defer file.Close()

	// 读取 ZIP 文件内容
	zipData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read zip file"})
		return
	}

	// 创建 ZIP 读取器
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid zip file"})
		return
	}

	// 解析节点配置
	nodeConfigs, err := service.ExtractNodeConfigs(zipReader)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 批量处理节点配置
	results := make([]map[string]interface{}, 0)
	successCount := 0
	errorCount := 0

	for _, nodeConfig := range nodeConfigs {
		result := map[string]interface{}{
			"folder_name": nodeConfig.FolderName,
		}

		// 解析配置文件
		config, err := service.ParseQuilibriumConfig(nodeConfig.ConfigYml)
		if err != nil {
			result["status"] = "error"
			result["error"] = "Invalid config.yml: " + err.Error()
			errorCount++
			results = append(results, result)
			continue
		}

		// 验证 keys.yml
		if err := service.ValidateConfigFiles(nodeConfig.ConfigYml, nodeConfig.KeysYml); err != nil {
			result["status"] = "error"
			result["error"] = err.Error()
			errorCount++
			results = append(results, result)
			continue
		}

		// 创建或更新节点
		node := &model.Node{
			Name:      nodeConfig.FolderName,
			PeerKey:   config.GetPeerPrivKey(),
			ConfigYml: string(nodeConfig.ConfigYml),
			KeysYml:   string(nodeConfig.KeysYml),
			Enabled:   true,
		}

		// 检查 BasePort 是否需要分配
		if node.BasePort == 0 {
			basePort, err := h.nodeService.GetNextAvailableBasePort()
			if err != nil {
				result["status"] = "error"
				result["error"] = "Failed to allocate base port: " + err.Error()
				errorCount++
				results = append(results, result)
				continue
			}
			node.BasePort = basePort
		}

		// 检查 SortNo 是否需要分配
		if node.SortNo == 0 {
			sortNo, err := h.nodeService.GetNextAvailableSortNo()
			if err != nil {
				result["status"] = "error"
				result["error"] = "Failed to allocate sort number: " + err.Error()
				errorCount++
				results = append(results, result)
				continue
			}
			node.SortNo = sortNo
		}

		if err := h.nodeService.UpsertNodeByPeerKey(node); err != nil {
			result["status"] = "error"
			result["error"] = err.Error()
			errorCount++
		} else {
			result["status"] = "success"
			result["node_id"] = node.ID
			result["base_port"] = node.BasePort
			successCount++
		}

		results = append(results, result)
	}

	// 重新生成 Docker Compose 配置
	if successCount > 0 {
		if err := h.dockerService.GenerateDockerCompose(); err != nil {
			zap.L().Error("Failed to generate docker-compose", zap.Error(err))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Batch upload completed",
		"total":         len(nodeConfigs),
		"success_count": successCount,
		"error_count":   errorCount,
		"results":       results,
	})
}

// 获取所有节点状态
func (h *NodeHandler) GetNodesStatus(c *gin.Context) {
	nodes, err := h.nodeService.GetAllNodes()
	if err != nil {
		zap.L().Error("Failed to get nodes", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取容器状态
	var containerStatus map[string]string
	if h.dockerClientService != nil {
		containerStatus, err = h.dockerClientService.GetAllNodeStatus("easy-node")
		if err != nil {
			zap.L().Error("Failed to get container status", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 模拟容器状态
		containerStatus = make(map[string]string)
		for _, node := range nodes {
			containerName := fmt.Sprintf("quilibrium-node-%d", node.SortNo)
			containerStatus[containerName] = "running"
		}
	}

	// 合并节点信息和容器状态
	var nodeStatus []map[string]interface{}
	for _, node := range nodes {
		containerName := fmt.Sprintf("quilibrium-node-%d", node.SortNo)
		status := "unknown"
		if containerStatus[containerName] != "" {
			status = containerStatus[containerName]
		}

		nodeStatus = append(nodeStatus, map[string]interface{}{
			"id":              node.ID,
			"name":            node.Name,
			"sort_no":         node.SortNo,
			"container_name":  containerName,
			"container_status": status,
			"base_port":       node.BasePort,
			"enabled":         node.Enabled,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"nodes": nodeStatus,
	})
}

// 启动节点
func (h *NodeHandler) StartNode(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node ID"})
		return
	}

	node, err := h.nodeService.GetNodeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Node not found"})
		return
	}

	containerName := fmt.Sprintf("quilibrium-node-%d", node.SortNo)
	if h.dockerClientService != nil {
		if err := h.dockerClientService.StartContainer(containerName); err != nil {
			zap.L().Error("Failed to start container", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 模拟启动成功
		zap.L().Info("Simulated node start", zap.String("container", containerName))
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Node %s started successfully", node.Name),
	})
}

// 停止节点
func (h *NodeHandler) StopNode(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node ID"})
		return
	}

	node, err := h.nodeService.GetNodeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Node not found"})
		return
	}

	containerName := fmt.Sprintf("quilibrium-node-%d", node.SortNo)
	if h.dockerClientService != nil {
		if err := h.dockerClientService.StopContainer(containerName); err != nil {
			zap.L().Error("Failed to stop container", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 模拟停止成功
		zap.L().Info("Simulated node stop", zap.String("container", containerName))
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Node %s stopped successfully", node.Name),
	})
}

// 重启节点
func (h *NodeHandler) RestartNode(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node ID"})
		return
	}

	node, err := h.nodeService.GetNodeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Node not found"})
		return
	}

	containerName := fmt.Sprintf("quilibrium-node-%d", node.SortNo)
	if h.dockerClientService != nil {
		if err := h.dockerClientService.RestartContainer(containerName); err != nil {
			zap.L().Error("Failed to restart container", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 模拟重启成功
		zap.L().Info("Simulated node restart", zap.String("container", containerName))
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Node %s restarted successfully", node.Name),
	})
}

// 获取节点日志
func (h *NodeHandler) GetNodeLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node ID"})
		return
	}

	node, err := h.nodeService.GetNodeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Node not found"})
		return
	}

	// 获取日志行数参数，默认100行
	lines := 100
	if linesParam := c.Query("lines"); linesParam != "" {
		if parsedLines, err := strconv.Atoi(linesParam); err == nil && parsedLines > 0 {
			lines = parsedLines
		}
	}

	containerName := fmt.Sprintf("quilibrium-node-%d", node.SortNo)
	var logs string
	if h.dockerClientService != nil {
		logs, err = h.dockerClientService.GetContainerLogs(containerName, lines)
		if err != nil {
			zap.L().Error("Failed to get container logs", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 模拟日志
		logs = fmt.Sprintf("Simulated logs for %s\nNode is running normally\nLatest activity: %s", node.Name, time.Now().Format("2006-01-02 15:04:05"))
	}

	c.JSON(http.StatusOK, gin.H{
		"node_name": node.Name,
		"lines":     lines,
		"logs":      logs,
	})
}

// 启动所有节点
func (h *NodeHandler) StartAllNodes(c *gin.Context) {
	if h.dockerClientService != nil {
		if err := h.dockerClientService.StartComposeProject("easy-node"); err != nil {
			zap.L().Error("Failed to start all nodes", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 模拟启动所有节点
		zap.L().Info("Simulated start all nodes")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All nodes started successfully",
	})
}

// 停止所有节点
func (h *NodeHandler) StopAllNodes(c *gin.Context) {
	if h.dockerClientService != nil {
		if err := h.dockerClientService.StopComposeProject("easy-node"); err != nil {
			zap.L().Error("Failed to stop all nodes", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 模拟停止所有节点
		zap.L().Info("Simulated stop all nodes")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All nodes stopped successfully",
	})
}

// 获取健康统计信息
func (h *NodeHandler) GetHealthStats(c *gin.Context) {
	stats, err := h.healthService.GetHealthStats()
	if err != nil {
		zap.L().Error("Failed to get health stats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// 获取所有节点健康详情
func (h *NodeHandler) GetNodesHealth(c *gin.Context) {
	healthInfos, err := h.healthService.GetNodeHealthDetails()
	if err != nil {
		zap.L().Error("Failed to get nodes health", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"nodes": healthInfos,
	})
}

// 检查单个节点健康状态
func (h *NodeHandler) CheckNodeHealth(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node ID"})
		return
	}

	healthInfo, err := h.healthService.CheckNodeHealth(uint(id))
	if err != nil {
		zap.L().Error("Failed to check node health", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, healthInfo)
}

// 检查所有节点健康状态
func (h *NodeHandler) CheckAllNodesHealth(c *gin.Context) {
	healthInfos, err := h.healthService.CheckAllNodesHealth()
	if err != nil {
		zap.L().Error("Failed to check all nodes health", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"nodes": healthInfos,
		"total": len(healthInfos),
	})
}

// 批量创建节点
func (h *NodeHandler) BatchCreateNodes(c *gin.Context) {
	var req struct {
		Count      int    `json:"count" binding:"required,min=1,max=1000"`
		NamePrefix string `json:"name_prefix" binding:"required"`
		StartPort  int    `json:"start_port" binding:"required,min=1024,max=65000"`
		Enabled    bool   `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证端口范围不会超出
	if req.StartPort+(req.Count*5) > 65535 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Port range exceeds maximum (65535)"})
		return
	}

	// 批量创建节点
	results := make([]map[string]interface{}, 0)
	successCount := 0
	errorCount := 0

	for i := 0; i < req.Count; i++ {
		result := map[string]interface{}{
			"index": i + 1,
			"name":  fmt.Sprintf("%s-%d", req.NamePrefix, i+1),
		}

		// 为每个节点获取新的 sortNo（Redis 原子递增）
		sortNo, err := h.nodeService.GetNextAvailableSortNo()
		if err != nil {
			result["status"] = "error"
			result["error"] = "Failed to get next sort number: " + err.Error()
			errorCount++
			results = append(results, result)
			continue
		}

		// 生成密钥
		peerPrivKey, err := service.GeneratePeerPrivKey()
		if err != nil {
			result["status"] = "error"
			result["error"] = "Failed to generate peer private key: " + err.Error()
			errorCount++
			results = append(results, result)
			continue
		}

		encryptionKey, err := service.GenerateEncryptionKey()
		if err != nil {
			result["status"] = "error"
			result["error"] = "Failed to generate encryption key: " + err.Error()
			errorCount++
			results = append(results, result)
			continue
		}

		// 生成默认配置
		configYml := service.GenerateDefaultConfigYml(peerPrivKey, encryptionKey)
		keysYml := service.GenerateDefaultKeysYml()

		// 创建节点
		node := &model.Node{
			Name:      fmt.Sprintf("%s-%d", req.NamePrefix, i+1),
			SortNo:    sortNo,
			BasePort:  req.StartPort + (i * 5),
			PeerKey:   peerPrivKey,
			ConfigYml: configYml,
			KeysYml:   keysYml,
			Enabled:   req.Enabled,
		}

		if err := h.nodeService.CreateNode(node); err != nil {
			result["status"] = "error"
			result["error"] = err.Error()
			errorCount++
		} else {
			result["status"] = "success"
			result["node_id"] = node.ID
			result["sort_no"] = node.SortNo
			result["base_port"] = node.BasePort
			successCount++
		}

		results = append(results, result)
	}

	// 重新生成 Docker Compose 配置
	if successCount > 0 {
		if err := h.dockerService.GenerateDockerCompose(); err != nil {
			zap.L().Error("Failed to generate docker-compose", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":         "Nodes created but failed to update docker-compose",
				"success_count": successCount,
				"error_count":   errorCount,
				"results":       results,
			})
			return
		}

		// 启动容器 (如果 enabled=true)
		if req.Enabled && h.dockerClientService != nil {
			if err := h.dockerClientService.StartComposeProject("easy-node"); err != nil {
				zap.L().Error("Failed to start containers", zap.Error(err))
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Batch create completed",
		"total":         req.Count,
		"success_count": successCount,
		"error_count":   errorCount,
		"results":       results,
	})
}