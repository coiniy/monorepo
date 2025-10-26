package service

import (
	"easy-node-backend/internal/model"
	"fmt"
	"gorm.io/gorm"
)

type NodeService struct {
	db *gorm.DB
}

func NewNodeService(db *gorm.DB) *NodeService {
	return &NodeService{db: db}
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