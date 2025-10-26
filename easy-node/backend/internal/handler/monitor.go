package handler

import (
	"easy-node-backend/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MonitorHandler struct {
	schedulerService *service.SchedulerService
	healthService    *service.HealthService
}

func NewMonitorHandler(schedulerService *service.SchedulerService, healthService *service.HealthService) *MonitorHandler {
	return &MonitorHandler{
		schedulerService: schedulerService,
		healthService:    healthService,
	}
}

// 获取调度器状态
func (h *MonitorHandler) GetSchedulerStatus(c *gin.Context) {
	status := h.schedulerService.GetStatus()
	c.JSON(http.StatusOK, status)
}

// 启动调度器
func (h *MonitorHandler) StartScheduler(c *gin.Context) {
	h.schedulerService.Start()
	c.JSON(http.StatusOK, gin.H{
		"message": "Health check scheduler started",
		"status":  h.schedulerService.GetStatus(),
	})
}

// 停止调度器
func (h *MonitorHandler) StopScheduler(c *gin.Context) {
	h.schedulerService.Stop()
	c.JSON(http.StatusOK, gin.H{
		"message": "Health check scheduler stopped",
	})
}

// 更新检查间隔
func (h *MonitorHandler) UpdateInterval(c *gin.Context) {
	var request struct {
		Interval int `json:"interval" binding:"required,min=10"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.schedulerService.UpdateInterval(time.Duration(request.Interval) * time.Second)
	c.JSON(http.StatusOK, gin.H{
		"message":      "Health check interval updated",
		"new_interval": request.Interval,
	})
}

// 手动触发健康检查
func (h *MonitorHandler) TriggerHealthCheck(c *gin.Context) {
	h.schedulerService.TriggerHealthCheck()
	c.JSON(http.StatusOK, gin.H{
		"message": "Manual health check triggered",
	})
}

// 获取性能指标概览
func (h *MonitorHandler) GetPerformanceMetrics(c *gin.Context) {
	// 获取健康统计
	healthStats, err := h.healthService.GetHealthStats()
	if err != nil {
		zap.L().Error("Failed to get health stats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取所有节点健康详情
	healthDetails, err := h.healthService.GetNodeHealthDetails()
	if err != nil {
		zap.L().Error("Failed to get health details", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 计算性能指标
	var totalResponseTime float64
	var onlineNodes int
	var maxResponseTime float64
	var minResponseTime float64 = 999999
	
	for _, detail := range healthDetails {
		if detail.HealthStatus == "healthy" {
			onlineNodes++
			totalResponseTime += detail.ResponseTime
			
			if detail.ResponseTime > maxResponseTime {
				maxResponseTime = detail.ResponseTime
			}
			if detail.ResponseTime < minResponseTime && detail.ResponseTime > 0 {
				minResponseTime = detail.ResponseTime
			}
		}
	}
	
	if onlineNodes == 0 {
		minResponseTime = 0
	}

	metrics := map[string]interface{}{
		"overview": healthStats,
		"performance": map[string]interface{}{
			"total_nodes":         len(healthDetails),
			"online_nodes":        onlineNodes,
			"avg_response_time":   totalResponseTime / float64(onlineNodes),
			"max_response_time":   maxResponseTime,
			"min_response_time":   minResponseTime,
			"uptime_percentage":   float64(onlineNodes) / float64(len(healthDetails)) * 100,
		},
		"scheduler": h.schedulerService.GetStatus(),
		"timestamp": time.Now(),
	}

	c.JSON(http.StatusOK, metrics)
}

// 获取历史趋势数据
func (h *MonitorHandler) GetTrendData(c *gin.Context) {
	// 获取时间范围参数
	hoursParam := c.DefaultQuery("hours", "24")
	hours, err := strconv.Atoi(hoursParam)
	if err != nil || hours < 1 || hours > 168 { // 最多7天
		hours = 24
	}

	// 这里可以扩展为从数据库获取历史数据
	// 目前返回模拟数据结构
	trendData := map[string]interface{}{
		"time_range":   hours,
		"data_points":  []interface{}{}, // 这里可以添加历史数据点
		"message":      "Historical trend data feature is under development",
	}

	c.JSON(http.StatusOK, trendData)
}