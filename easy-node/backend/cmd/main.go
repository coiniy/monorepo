package main

import (
	"easy-node-backend/internal/config"
	"easy-node-backend/internal/db"
	"easy-node-backend/internal/handler"
	redisClient "easy-node-backend/internal/redis"
	"easy-node-backend/internal/service"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	database, err := db.InitDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 初始化 Redis
	redis, err := redisClient.InitRedis(cfg)
	if err != nil {
		log.Fatal("Failed to initialize Redis:", err)
	}

	nodeService := service.NewNodeService(database, redis)

	// 初始化 sortNo 计数器
	if err := nodeService.InitializeSortNoCounter(); err != nil {
		zap.L().Warn("Failed to initialize sortNo counter", zap.Error(err))
	}

	dockerService := service.NewDockerService(database, cfg)
	
	// 初始化健康检查服务（暂时去掉Docker依赖）
	healthService := service.NewHealthService(database, nil)
	
	// 初始化监控调度器 (30秒间隔)
	schedulerService := service.NewSchedulerService(healthService, 30*time.Second)
	
	// 启动健康监控调度器
	schedulerService.Start()
	defer schedulerService.Stop()
	
	nodeHandler := handler.NewNodeHandler(nodeService, dockerService, nil, healthService)
	monitorHandler := handler.NewMonitorHandler(schedulerService, healthService)

	gin.SetMode(cfg.Server.Mode)
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		nodes := api.Group("/nodes")
		{
			// 节点信息管理
			nodes.GET("", nodeHandler.GetNodes)
			nodes.POST("", nodeHandler.CreateNode)
			nodes.POST("/batch-create", nodeHandler.BatchCreateNodes)
			nodes.POST("/upload", nodeHandler.CreateNodeFromConfig)
			nodes.POST("/batch-upload", nodeHandler.BatchUploadNodes)
			nodes.PUT("/:id", nodeHandler.UpdateNode)
			nodes.DELETE("/:id", nodeHandler.DeleteNode)
			
			// 节点状态管理
			nodes.GET("/status", nodeHandler.GetNodesStatus)
			nodes.POST("/:id/start", nodeHandler.StartNode)
			nodes.POST("/:id/stop", nodeHandler.StopNode)
			nodes.POST("/:id/restart", nodeHandler.RestartNode)
			nodes.GET("/:id/logs", nodeHandler.GetNodeLogs)
			
			// 健康监控
			nodes.GET("/health", nodeHandler.GetNodesHealth)
			nodes.GET("/health/stats", nodeHandler.GetHealthStats)
			nodes.GET("/:id/health", nodeHandler.CheckNodeHealth)
			nodes.POST("/health/check", nodeHandler.CheckAllNodesHealth)
		}
		
		// Docker 操作
		docker := api.Group("/docker")
		{
			docker.POST("/regenerate", nodeHandler.RegenerateDockerCompose)
			docker.POST("/start-all", nodeHandler.StartAllNodes)
			docker.POST("/stop-all", nodeHandler.StopAllNodes)
		}
		
		// 监控和性能指标
		monitor := api.Group("/monitor")
		{
			monitor.GET("/scheduler/status", monitorHandler.GetSchedulerStatus)
			monitor.POST("/scheduler/start", monitorHandler.StartScheduler)
			monitor.POST("/scheduler/stop", monitorHandler.StopScheduler)
			monitor.PUT("/scheduler/interval", monitorHandler.UpdateInterval)
			monitor.POST("/health/trigger", monitorHandler.TriggerHealthCheck)
			monitor.GET("/metrics", monitorHandler.GetPerformanceMetrics)
			monitor.GET("/trends", monitorHandler.GetTrendData)
		}
	}

	zap.L().Info("Server starting", zap.String("port", cfg.Server.Port))
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}