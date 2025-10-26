package service

import (
	"context"
	"easy-node-backend/internal/model"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	RedisSortNoKey = "quilibrium:node:sortno:counter"
)

type NodeService struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewNodeService(db *gorm.DB, redisClient *redis.Client) *NodeService {
	return &NodeService{
		db:    db,
		redis: redisClient,
	}
}

func (s *NodeService) GetAllNodes() ([]model.Node, error) {
	var nodes []model.Node
	err := s.db.Where("enabled = ?", true).Find(&nodes).Error
	return nodes, err
}

func (s *NodeService) CreateNode(node *model.Node) error {
	return s.db.Create(node).Error
}

func (s *NodeService) UpdateNode(node *model.Node) error {
	return s.db.Save(node).Error
}

func (s *NodeService) DeleteNode(id uint) error {
	return s.db.Delete(&model.Node{}, id).Error
}

func (s *NodeService) GetNodeByID(id uint) (*model.Node, error) {
	var node model.Node
	err := s.db.First(&node, id).Error
	return &node, err
}

func (s *NodeService) CreateNodeFromConfig(node *model.Node) error {
	// 检查 PeerKey 是否已存在
	var existingNode model.Node
	err := s.db.Where("peer_key = ?", node.PeerKey).First(&existingNode).Error
	if err == nil {
		return fmt.Errorf("node with peer_key %s already exists", node.PeerKey)
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing peer_key: %w", err)
	}

	// 检查节点名是否已存在
	err = s.db.Where("name = ?", node.Name).First(&existingNode).Error
	if err == nil {
		return fmt.Errorf("node with name %s already exists", node.Name)
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing name: %w", err)
	}

	// 检查端口是否已被使用
	err = s.db.Where("base_port = ?", node.BasePort).First(&existingNode).Error
	if err == nil {
		return fmt.Errorf("base_port %d already in use", node.BasePort)
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing base_port: %w", err)
	}

	// 创建节点
	return s.db.Create(node).Error
}

func (s *NodeService) UpsertNodeByPeerKey(node *model.Node) error {
	var existingNode model.Node
	err := s.db.Where("peer_key = ?", node.PeerKey).First(&existingNode).Error
	
	if err == gorm.ErrRecordNotFound {
		// 节点不存在，创建新节点
		// 检查名称和端口冲突
		if err := s.checkNameConflict(node.Name, 0); err != nil {
			return err
		}
		if err := s.checkPortConflict(node.BasePort, 0); err != nil {
			return err
		}
		return s.db.Create(node).Error
	} else if err != nil {
		return fmt.Errorf("failed to check existing peer_key: %w", err)
	}
	
	// 节点存在，更新配置
	// 检查名称和端口冲突（排除当前节点）
	if existingNode.Name != node.Name {
		if err := s.checkNameConflict(node.Name, existingNode.ID); err != nil {
			return err
		}
	}
	if existingNode.BasePort != node.BasePort {
		if err := s.checkPortConflict(node.BasePort, existingNode.ID); err != nil {
			return err
		}
	}
	
	// 更新现有节点
	node.ID = existingNode.ID
	node.CreatedAt = existingNode.CreatedAt
	return s.db.Save(node).Error
}

func (s *NodeService) checkNameConflict(name string, excludeID uint) error {
	var existingNode model.Node
	query := s.db.Where("name = ?", name)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	
	err := query.First(&existingNode).Error
	if err == nil {
		return fmt.Errorf("node with name %s already exists", name)
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing name: %w", err)
	}
	return nil
}

func (s *NodeService) checkPortConflict(basePort int, excludeID uint) error {
	var existingNode model.Node
	query := s.db.Where("base_port = ?", basePort)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	
	err := query.First(&existingNode).Error
	if err == nil {
		return fmt.Errorf("base_port %d already in use", basePort)
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing base_port: %w", err)
	}
	return nil
}

func (s *NodeService) GetNextAvailableBasePort() (int, error) {
	var maxPort int
	err := s.db.Model(&model.Node{}).Select("COALESCE(MAX(base_port), 5000)").Scan(&maxPort).Error
	if err != nil {
		return 0, fmt.Errorf("failed to get max base_port: %w", err)
	}

	// 下一个可用端口，每个节点间隔10个端口
	nextPort := maxPort + 10

	// 确保端口没有被使用
	for {
		var count int64
		err := s.db.Model(&model.Node{}).Where("base_port = ?", nextPort).Count(&count).Error
		if err != nil {
			return 0, fmt.Errorf("failed to check port availability: %w", err)
		}
		if count == 0 {
			break
		}
		nextPort += 10
	}

	return nextPort, nil
}

// GetNextAvailableSortNo 获取下一个可用的排序号（使用Redis原子递增）
func (s *NodeService) GetNextAvailableSortNo() (int, error) {
	ctx := context.Background()

	// 使用Redis INCR命令原子递增
	sortNo, err := s.redis.Incr(ctx, RedisSortNoKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment sortNo in Redis: %w", err)
	}

	return int(sortNo), nil
}

// InitializeSortNoCounter 初始化Redis中的sortNo计数器（从数据库的最大值开始）
func (s *NodeService) InitializeSortNoCounter() error {
	ctx := context.Background()

	// 检查Redis中是否已经有计数器
	exists, err := s.redis.Exists(ctx, RedisSortNoKey).Result()
	if err != nil {
		return fmt.Errorf("failed to check Redis key existence: %w", err)
	}

	// 如果计数器已存在，不需要初始化
	if exists > 0 {
		return nil
	}

	// 从数据库获取当前最大的sortNo
	var maxSortNo int
	err = s.db.Model(&model.Node{}).Select("COALESCE(MAX(sort_no), 0)").Scan(&maxSortNo).Error
	if err != nil {
		return fmt.Errorf("failed to get max sort_no: %w", err)
	}

	// 设置Redis计数器为当前最大值
	if err := s.redis.Set(ctx, RedisSortNoKey, maxSortNo, 0).Err(); err != nil {
		return fmt.Errorf("failed to initialize sortNo counter in Redis: %w", err)
	}

	return nil
}