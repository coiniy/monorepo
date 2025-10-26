package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"easy-node-backend/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type HealthService struct {
	db                  *gorm.DB
	dockerClientService *DockerClientService
	httpClient          *http.Client
}

type PrometheusMetrics struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

type NodeHealthInfo struct {
	NodeID          uint      `json:"node_id"`
	NodeName        string    `json:"node_name"`
	ContainerStatus string    `json:"container_status"`
	HealthStatus    string    `json:"health_status"`
	LastCheck       time.Time `json:"last_check"`
	ResponseTime    float64   `json:"response_time"`
	HealthMessage   string    `json:"health_message"`
	Metrics         map[string]interface{} `json:"metrics,omitempty"`
}

func NewHealthService(db *gorm.DB, dockerClientService *DockerClientService) *HealthService {
	return &HealthService{
		db:                  db,
		dockerClientService: dockerClientService,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// 检查单个节点健康状态
func (s *HealthService) CheckNodeHealth(nodeID uint) (*NodeHealthInfo, error) {
	var node model.Node
	if err := s.db.First(&node, nodeID).Error; err != nil {
		return nil, fmt.Errorf("node not found: %w", err)
	}

	healthInfo := &NodeHealthInfo{
		NodeID:   node.ID,
		NodeName: node.Name,
		LastCheck: time.Now(),
	}

	// 1. 检查容器状态（模拟）
	containerName := fmt.Sprintf("quilibrium-node-%d", node.ID)
	var containerStatus string
	var err error
	
	if s.dockerClientService != nil {
		containerStatus, err = s.dockerClientService.GetContainerStatus(containerName)
		if err != nil {
			healthInfo.ContainerStatus = "error"
			healthInfo.HealthStatus = "unhealthy"
			healthInfo.HealthMessage = fmt.Sprintf("Failed to get container status: %v", err)
			s.updateNodeHealth(&node, healthInfo)
			return healthInfo, nil
		}
	} else {
		// 模拟容器状态
		containerStatus = "running"
	}

	healthInfo.ContainerStatus = containerStatus

	// 2. 如果容器未运行，直接返回
	if containerStatus != "running" {
		healthInfo.HealthStatus = "unhealthy"
		healthInfo.HealthMessage = fmt.Sprintf("Container is %s", containerStatus)
		s.updateNodeHealth(&node, healthInfo)
		return healthInfo, nil
	}

	// 3. 检查端口连通性和响应时间
	if err := s.checkNodeConnectivity(&node, healthInfo); err != nil {
		healthInfo.HealthStatus = "unhealthy"
		healthInfo.HealthMessage = err.Error()
		s.updateNodeHealth(&node, healthInfo)
		return healthInfo, nil
	}

	// 4. 获取 Prometheus 指标
	if err := s.fetchPrometheusMetrics(&node, healthInfo); err != nil {
		zap.L().Warn("Failed to fetch prometheus metrics", zap.Uint("node_id", node.ID), zap.Error(err))
	}

	// 5. 综合判断健康状态
	healthInfo.HealthStatus = "healthy"
	healthInfo.HealthMessage = "Node is healthy"

	s.updateNodeHealth(&node, healthInfo)
	return healthInfo, nil
}

// 检查节点连通性
func (s *HealthService) checkNodeConnectivity(node *model.Node, healthInfo *NodeHealthInfo) error {
	start := time.Now()
	
	// 检查 Prometheus 端口 (BasePort+4 -> 8080)
	prometheusURL := fmt.Sprintf("http://localhost:%d/metrics", node.BasePort+4)
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	req, err := http.NewRequestWithContext(ctx, "GET", prometheusURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to prometheus endpoint: %w", err)
	}
	defer resp.Body.Close()
	
	healthInfo.ResponseTime = float64(time.Since(start).Nanoseconds()) / 1e6 // 转换为毫秒
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("prometheus endpoint returned status %d", resp.StatusCode)
	}
	
	return nil
}

// 获取 Prometheus 指标
func (s *HealthService) fetchPrometheusMetrics(node *model.Node, healthInfo *NodeHealthInfo) error {
	prometheusURL := fmt.Sprintf("http://localhost:%d/metrics", node.BasePort+4)
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	req, err := http.NewRequestWithContext(ctx, "GET", prometheusURL, nil)
	if err != nil {
		return err
	}
	
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("prometheus returned status %d", resp.StatusCode)
	}
	
	// 这里可以解析 Prometheus 指标，简化处理
	healthInfo.Metrics = map[string]interface{}{
		"prometheus_available": true,
		"last_scrape": time.Now(),
	}
	
	return nil
}

// 更新节点健康状态到数据库
func (s *HealthService) updateNodeHealth(node *model.Node, healthInfo *NodeHealthInfo) {
	now := time.Now()
	updates := map[string]interface{}{
		"health_status":     healthInfo.HealthStatus,
		"last_health_check": &now,
		"health_message":    healthInfo.HealthMessage,
		"response_time":     healthInfo.ResponseTime,
		"updated_at":        now,
	}
	
	if err := s.db.Model(node).Updates(updates).Error; err != nil {
		zap.L().Error("Failed to update node health", zap.Uint("node_id", node.ID), zap.Error(err))
	}
}

// 检查所有节点健康状态
func (s *HealthService) CheckAllNodesHealth() ([]NodeHealthInfo, error) {
	var nodes []model.Node
	if err := s.db.Where("enabled = ?", true).Find(&nodes).Error; err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
	}

	var healthInfos []NodeHealthInfo
	for _, node := range nodes {
		healthInfo, err := s.CheckNodeHealth(node.ID)
		if err != nil {
			zap.L().Error("Failed to check node health", zap.Uint("node_id", node.ID), zap.Error(err))
			// 创建错误状态信息
			healthInfo = &NodeHealthInfo{
				NodeID:        node.ID,
				NodeName:      node.Name,
				HealthStatus:  "error",
				LastCheck:     time.Now(),
				HealthMessage: err.Error(),
			}
		}
		healthInfos = append(healthInfos, *healthInfo)
	}

	return healthInfos, nil
}

// 获取节点健康统计
func (s *HealthService) GetHealthStats() (map[string]interface{}, error) {
	var totalCount int64
	var healthyCount int64
	var unhealthyCount int64
	var unknownCount int64

	// 获取总数
	if err := s.db.Model(&model.Node{}).Where("enabled = ?", true).Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// 获取健康状态统计
	if err := s.db.Model(&model.Node{}).Where("enabled = ? AND health_status = ?", true, "healthy").Count(&healthyCount).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&model.Node{}).Where("enabled = ? AND health_status = ?", true, "unhealthy").Count(&unhealthyCount).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&model.Node{}).Where("enabled = ? AND health_status = ?", true, "unknown").Count(&unknownCount).Error; err != nil {
		return nil, err
	}

	// 计算平均响应时间
	var avgResponseTime float64
	if err := s.db.Model(&model.Node{}).Where("enabled = ? AND health_status = ?", true, "healthy").Select("AVG(response_time)").Scan(&avgResponseTime).Error; err != nil {
		avgResponseTime = 0
	}

	return map[string]interface{}{
		"total_nodes":        totalCount,
		"healthy_nodes":      healthyCount,
		"unhealthy_nodes":    unhealthyCount,
		"unknown_nodes":      unknownCount,
		"health_percentage":  float64(healthyCount) / float64(totalCount) * 100,
		"avg_response_time":  avgResponseTime,
		"last_check":         time.Now(),
	}, nil
}

// 获取节点详细健康信息
func (s *HealthService) GetNodeHealthDetails() ([]NodeHealthInfo, error) {
	var nodes []model.Node
	if err := s.db.Where("enabled = ?", true).Find(&nodes).Error; err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
	}

	var healthInfos []NodeHealthInfo
	for _, node := range nodes {
		// 获取容器状态
		containerName := fmt.Sprintf("quilibrium-node-%d", node.ID)
		var containerStatus string
		if s.dockerClientService != nil {
			containerStatus, _ = s.dockerClientService.GetContainerStatus(containerName)
		} else {
			containerStatus = "running" // 模拟状态
		}

		healthInfo := NodeHealthInfo{
			NodeID:          node.ID,
			NodeName:        node.Name,
			ContainerStatus: containerStatus,
			HealthStatus:    node.HealthStatus,
			HealthMessage:   node.HealthMessage,
			ResponseTime:    node.ResponseTime,
		}

		if node.LastHealthCheck != nil {
			healthInfo.LastCheck = *node.LastHealthCheck
		}

		healthInfos = append(healthInfos, healthInfo)
	}

	return healthInfos, nil
}