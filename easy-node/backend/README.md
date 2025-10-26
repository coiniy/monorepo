# Easy Node Backend

基于 Gin + GORM + PostgreSQL + Zap + Viper 的 Quilibrium 节点管理后端，用于动态管理节点配置和部署环境。

## 功能特性

- **节点信息管理** - 完整的 CRUD 操作
- **配置文件管理** - 支持 config.yml 和 keys.yml 上传解析
- **批量节点部署** - ZIP 压缩包批量上传节点配置
- **智能端口分配** - 自动分配端口避免冲突
- **动态部署生成** - 根据数据库生成完整部署环境
- **负载均衡代理** - 使用 Envoy 代理 TCP/UDP 流量
- **配置文件同步** - 数据库与文件系统配置同步
- **节点状态管理** - 启动/停止/重启容器控制
- **Docker 容器操作** - 集成 Docker API 管理容器生命周期
- **基础监控系统** - 节点健康检查和状态显示
- **性能指标监控** - 响应时间、健康状态统计
- **自动化调度器** - 定时健康检查和状态更新

## 架构组件

- **API 服务器** - Gin 框架 RESTful API
- **数据库** - PostgreSQL 存储节点配置
- **负载均衡** - Envoy 代理处理 TCP/UDP 流量
- **容器编排** - Docker Compose 管理节点容器
- **配置管理** - 动态生成节点配置文件
- **Docker 集成** - Docker SDK 容器生命周期管理
- **健康监控** - Prometheus 指标收集和健康检查
- **调度系统** - 自动化监控任务调度器

## API 接口

### 节点管理

- `GET /api/v1/nodes` - 获取所有节点列表
- `POST /api/v1/nodes` - 创建新节点（JSON）
- `POST /api/v1/nodes/upload` - 上传配置文件创建节点
- `POST /api/v1/nodes/batch-upload` - 批量上传节点（ZIP）
- `PUT /api/v1/nodes/:id` - 更新指定节点
- `DELETE /api/v1/nodes/:id` - 删除指定节点

### 节点状态管理

- `GET /api/v1/nodes/status` - 获取所有节点状态
- `POST /api/v1/nodes/:id/start` - 启动指定节点
- `POST /api/v1/nodes/:id/stop` - 停止指定节点
- `POST /api/v1/nodes/:id/restart` - 重启指定节点
- `GET /api/v1/nodes/:id/logs` - 获取节点日志

### 健康监控

- `GET /api/v1/nodes/health` - 获取所有节点健康详情
- `GET /api/v1/nodes/health/stats` - 获取健康统计信息
- `GET /api/v1/nodes/:id/health` - 检查单个节点健康状态
- `POST /api/v1/nodes/health/check` - 手动检查所有节点

### 监控管理

- `GET /api/v1/monitor/scheduler/status` - 获取调度器状态
- `POST /api/v1/monitor/scheduler/start` - 启动调度器
- `POST /api/v1/monitor/scheduler/stop` - 停止调度器
- `PUT /api/v1/monitor/scheduler/interval` - 更新检查间隔
- `POST /api/v1/monitor/health/trigger` - 手动触发健康检查
- `GET /api/v1/monitor/metrics` - 获取性能指标概览
- `GET /api/v1/monitor/trends` - 获取历史趋势数据

### 部署管理

- `POST /api/v1/docker/regenerate` - 重新生成完整部署环境
- `POST /api/v1/docker/start-all` - 启动所有节点
- `POST /api/v1/docker/stop-all` - 停止所有节点

## 安装和运行

### 1. 环境准备
```bash
# 安装 Go 依赖
go mod tidy

# 启动 PostgreSQL 数据库
docker run -d --name postgres \
  -e POSTGRES_DB=easynode \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 postgres:15
```

### 2. 配置设置
修改 `config.yaml` 配置文件：
```yaml
server:
  port: "8080"
  mode: "debug"

database:
  host: "localhost"
  port: 5432
  username: "postgres"
  password: "password"
  dbname: "easynode"
  sslmode: "disable"

docker:
  deployment_path: "../deployment"
  node_image: "."
```

### 3. 启动服务
```bash
# 启动后端服务
go run cmd/main.go
```

## 使用示例

### 单节点上传
```bash
curl -X POST http://localhost:8080/api/v1/nodes/upload \
  -F "config=@config.yml" \
  -F "keys=@keys.yml" \
  -F "name=node-test" \
  -F "base_port=5000"
```

### 批量节点上传
```bash
# ZIP 文件结构：
# nodes.zip
# ├── node1/config.yml
# ├── node1/keys.yml
# ├── node2/config.yml
# └── node2/keys.yml

curl -X POST http://localhost:8080/api/v1/nodes/batch-upload \
  -F "zip=@nodes.zip"
```

### 重新生成部署配置
```bash
curl -X POST http://localhost:8080/api/v1/docker/regenerate
```

### 节点状态管理
```bash
# 启动节点
curl -X POST http://localhost:8080/api/v1/nodes/1/start

# 停止节点
curl -X POST http://localhost:8080/api/v1/nodes/1/stop

# 重启节点
curl -X POST http://localhost:8080/api/v1/nodes/1/restart

# 查看节点日志
curl "http://localhost:8080/api/v1/nodes/1/logs?lines=50"
```

### 健康监控
```bash
# 获取所有节点健康状态
curl http://localhost:8080/api/v1/nodes/health

# 获取健康统计信息
curl http://localhost:8080/api/v1/nodes/health/stats

# 检查单个节点健康
curl http://localhost:8080/api/v1/nodes/1/health

# 手动触发健康检查
curl -X POST http://localhost:8080/api/v1/nodes/health/check
```

### 监控管理
```bash
# 获取性能指标概览
curl http://localhost:8080/api/v1/monitor/metrics

# 手动触发健康检查
curl -X POST http://localhost:8080/api/v1/monitor/health/trigger

# 更新检查间隔为60秒
curl -X PUT http://localhost:8080/api/v1/monitor/scheduler/interval \
  -H "Content-Type: application/json" \
  -d '{"interval": 60}'
```

## 数据库模型

节点模型包含以下字段：
- **ID**: 节点唯一标识
- **Name**: 节点名称（唯一）
- **Status**: 节点运行状态
- **BasePort**: 基础端口（自动映射 +0~+4 的端口）
- **Enabled**: 是否启用节点
- **PeerKey**: 节点对等密钥（唯一索引）
- **ConfigYml**: config.yml 配置内容
- **KeysYml**: keys.yml 密钥内容
- **HealthStatus**: 健康状态（unknown, healthy, unhealthy, checking）
- **LastHealthCheck**: 最后健康检查时间
- **HealthMessage**: 健康检查消息
- **ResponseTime**: 响应时间（毫秒）

## 部署文件生成

每次调用重新生成接口时，系统会自动创建：

```
deployment/
├── docker-compose.yml      # 容器编排配置
├── envoy.yaml             # 负载均衡配置
├── node1-config/
│   ├── config.yml         # 节点1配置
│   └── keys.yml          # 节点1密钥
├── node2-config/
│   ├── config.yml         # 节点2配置
│   └── keys.yml          # 节点2密钥
└── ...
```

## 负载均衡

使用 Envoy 代理提供：
- **TCP 负载均衡**: 8340 端口（带健康检查）
- **UDP 负载均衡**: 8336 端口（QUIC 流量）
- **管理界面**: 9901 端口

## 端口分配规则

每个节点占用5个端口：
- `BasePort+0`: 8340 端口映射
- `BasePort+1`: 8336 端口映射（UDP/QUIC）
- `BasePort+2`: 8337 端口映射
- `BasePort+3`: 8338 端口映射
- `BasePort+4`: 8080 端口映射（Prometheus）

## 项目结构

```
backend/
├── cmd/main.go                    # 主程序入口
├── internal/
│   ├── config/config.go           # 配置管理
│   ├── db/database.go            # 数据库连接
│   ├── handler/                  # HTTP 处理器
│   │   ├── node.go               # 节点管理处理器
│   │   └── monitor.go            # 监控管理处理器
│   ├── model/node.go             # 数据模型
│   └── service/                  # 业务逻辑
│       ├── node.go               # 节点服务
│       ├── docker.go             # 部署生成服务
│       ├── docker_client.go      # Docker API 客户端
│       ├── health.go             # 健康检查服务
│       ├── scheduler.go          # 监控调度器
│       ├── config.go             # 配置解析服务
│       └── zip.go                # ZIP 处理服务
├── config.yaml                   # 配置文件
└── README.md                     # 项目说明

## 监控系统

### 健康检查功能
- **容器状态检查** - 检查 Docker 容器运行状态
- **端口连通性测试** - 测试 Prometheus 端口连接
- **响应时间监控** - 记录节点响应时间
- **Prometheus 指标** - 收集节点性能指标

### 自动化调度器
- **定时检查** - 默认30秒间隔自动检查
- **可配置间隔** - 支持动态调整检查频率
- **并发检查** - 并行检查多个节点
- **状态持久化** - 健康状态保存到数据库

### 性能指标
- **节点健康统计** - 健康/不健康节点数量
- **平均响应时间** - 所有健康节点平均响应时间
- **最大/最小响应时间** - 响应时间范围统计
- **在线率** - 节点在线百分比
```