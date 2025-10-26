package service

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

type SchedulerService struct {
	healthService *HealthService
	interval      time.Duration
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	running       bool
	mu            sync.RWMutex
}

func NewSchedulerService(healthService *HealthService, interval time.Duration) *SchedulerService {
	if interval <= 0 {
		interval = 30 * time.Second // 默认30秒检查一次
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &SchedulerService{
		healthService: healthService,
		interval:      interval,
		ctx:           ctx,
		cancel:        cancel,
	}
}

// 启动监控调度器
func (s *SchedulerService) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		zap.L().Warn("Scheduler is already running")
		return
	}

	s.running = true
	s.wg.Add(1)

	go s.run()
	zap.L().Info("Health check scheduler started", zap.Duration("interval", s.interval))
}

// 停止监控调度器
func (s *SchedulerService) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	s.running = false
	s.cancel()
	s.wg.Wait()

	zap.L().Info("Health check scheduler stopped")
}

// 检查调度器是否运行中
func (s *SchedulerService) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// 更新检查间隔
func (s *SchedulerService) UpdateInterval(interval time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if interval <= 0 {
		interval = 30 * time.Second
	}

	s.interval = interval
	zap.L().Info("Health check interval updated", zap.Duration("new_interval", s.interval))
}

// 主运行循环
func (s *SchedulerService) run() {
	defer s.wg.Done()

	// 立即执行一次检查
	s.performHealthCheck()

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			zap.L().Info("Health check scheduler context cancelled")
			return
		case <-ticker.C:
			s.performHealthCheck()
			
			// 检查间隔是否有变化，如果有则重新设置ticker
			s.mu.RLock()
			currentInterval := s.interval
			s.mu.RUnlock()
			
			if ticker.C != time.NewTicker(currentInterval).C {
				ticker.Stop()
				ticker = time.NewTicker(currentInterval)
			}
		}
	}
}

// 执行健康检查
func (s *SchedulerService) performHealthCheck() {
	start := time.Now()
	
	healthInfos, err := s.healthService.CheckAllNodesHealth()
	if err != nil {
		zap.L().Error("Failed to perform scheduled health check", zap.Error(err))
		return
	}

	duration := time.Since(start)
	
	// 统计结果
	var healthyCount, unhealthyCount, errorCount int
	for _, info := range healthInfos {
		switch info.HealthStatus {
		case "healthy":
			healthyCount++
		case "unhealthy":
			unhealthyCount++
		default:
			errorCount++
		}
	}

	zap.L().Info("Scheduled health check completed",
		zap.Int("total", len(healthInfos)),
		zap.Int("healthy", healthyCount),
		zap.Int("unhealthy", unhealthyCount),
		zap.Int("error", errorCount),
		zap.Duration("duration", duration),
	)

	// 如果有不健康的节点，记录详细信息
	if unhealthyCount > 0 || errorCount > 0 {
		for _, info := range healthInfos {
			if info.HealthStatus != "healthy" {
				zap.L().Warn("Unhealthy node detected",
					zap.Uint("node_id", info.NodeID),
					zap.String("node_name", info.NodeName),
					zap.String("status", info.HealthStatus),
					zap.String("message", info.HealthMessage),
				)
			}
		}
	}
}

// 手动触发健康检查
func (s *SchedulerService) TriggerHealthCheck() {
	go s.performHealthCheck()
	zap.L().Info("Manual health check triggered")
}

// 获取调度器状态
func (s *SchedulerService) GetStatus() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"running":      s.running,
		"interval":     s.interval.String(),
		"next_check":   time.Now().Add(s.interval),
	}
}