package model

import (
	"gorm.io/gorm"
	"time"
)

type Node struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Name           string         `json:"name" gorm:"not null;unique"`
	Status         string         `json:"status" gorm:"default:stopped"`
	BasePort       int            `json:"base_port" gorm:"not null"`
	Enabled        bool           `json:"enabled" gorm:"default:true"`
	PeerKey        string         `json:"peer_key" gorm:"unique;index"`
	ConfigYml      string         `json:"config_yml" gorm:"type:text"`
	KeysYml        string         `json:"keys_yml" gorm:"type:text"`
	
	// 健康检查相关字段
	HealthStatus   string         `json:"health_status" gorm:"default:unknown"`    // unknown, healthy, unhealthy, checking
	LastHealthCheck *time.Time    `json:"last_health_check"`
	HealthMessage  string         `json:"health_message"`
	ResponseTime   float64        `json:"response_time"`                           // 响应时间（毫秒）
	
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type NodePortMapping struct {
	NodeID       uint   `json:"node_id"`
	ContainerPort int   `json:"container_port"`
	HostPort     int    `json:"host_port"`
}

func (Node) TableName() string {
	return "nodes"
}