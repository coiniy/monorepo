# 🎯 Quilibrium 运维监控体系

## 📊 Dashboard 架构

基于运维最佳实践，采用**分层监控**架构，快速定位问题：

```
Level 1: Overview (总览)      ← 30秒了解全局
    ↓
Level 2: Node Monitoring (详细监控)  ← 单节点深入分析
    ↓
Level 3: Specialized (专项分析)      ← Prover/Worker 专项
```

---

## 🚀 Dashboard 列表

### 1. 📍 Overview (总览) - **首选 Dashboard**
**文件**: `01-overview.json`
**用途**: 快速了解所有节点的健康状况
**访问**: http://192.168.1.6:3000/d/quilibrium-overview

#### 关键功能：
- **🚦 健康状态卡片** (6个)
  - Engine State (引擎状态) - 绿色=Proving
  - Active Provers (活跃 Prover) - 应该 > 0
  - Time Since Last Frame (距上次帧时间) - 应该 < 600s
  - Current Frame (当前帧号) - 持续增长
  - Active Workers (活跃 Worker) - 应该 = CPU 核心数
  - Network Peers (网络 Peer) - 应该 > 5

- **📋 节点对比表**
  - 所有节点的关键指标并排显示
  - 自动颜色编码：红色=问题，绿色=正常
  - 便于快速发现异常节点

- **📈 趋势图**
  - Frame 增长趋势
  - Time Since Last Frame 趋势
  - Prover 状态分布
  - 网络 Peer 数量变化

#### 使用场景：
✅ **每日巡检** - 查看这个 Dashboard 就够了
✅ **故障快速定位** - 红色/黄色指标表示需要关注
✅ **多节点对比** - 立即发现哪个节点有问题

---

### 2. 🔬 Node Monitoring (节点详细监控)
**文件**: `02-node-monitoring.json` (原 node-monitoring.json)
**用途**: 单个节点的详细性能分析

#### 包含指标：
- Engine State 详情
- Frame Processing/Proving Duration (P50/P95/P99)
- Network Message Delivery Rate
- Mesh Peer Counts
- Hypergraph Import Tree Size
- Event Distributor Subscribers

#### 使用场景：
🔍 **性能分析** - 查看帧处理时间分布
🔍 **网络诊断** - 检查消息投递率和 Peer 连接
🔍 **趋势监控** - 长期性能趋势

---

### 3. 🎖️ Prover Registry (Prover 专项)
**文件**: `03-prover-registry.json` (原 prover-registry.json)
**用途**: Prover 注册和状态监控

#### 包含指标：
- Total Allocations/Cached Provers/Shard Tries
- Provers by Status (active/joining/leaving/paused/left)
- Registry Refresh Rate & Duration
- Prover Operations Rate (add/remove/update)

#### 使用场景：
🎖️ **Prover 状态检查** - 确认 Active Prover 数量
🎖️ **注册问题诊断** - 查看是否卡在 joining 状态
🎖️ **性能优化** - 监控 Registry 刷新时间

---

### 4. ⚙️ Worker Manager (Worker 专项)
**文件**: `04-worker-manager.json` (原 worker/worker-manager.json)
**用途**: Worker 线程管理

#### 包含指标：
- Active Workers (应该 = CPU 核心数)
- Allocated Workers (当前分配的 Worker)
- Total Storage
- Worker Operation Duration

#### 使用场景：
⚙️ **资源检查** - 确认 Worker 数量正常
⚙️ **性能监控** - Worker 操作耗时

---

## 🎯 运维工作流程

### 场景 1: 日常巡检 (每天 1 次)
```
1. 打开 Overview Dashboard
2. 检查顶部 6 个状态卡片是否全绿
3. 查看节点对比表，确认所有节点正常
4. 完成 ✅ (耗时 < 1 分钟)
```

### 场景 2: 发现问题
```
1. Overview 显示某节点异常（红色/黄色）
2. 查看节点对比表，确认具体哪个指标异常
3. 点击相应的专项 Dashboard 深入分析
   - 引擎问题 → Node Monitoring
   - Prover 问题 → Prover Registry
   - Worker 问题 → Worker Manager
4. 根据指标变化定位根因
```

### 场景 3: 性能分析
```
1. Node Monitoring → Frame Processing Duration
2. 查看 P99/P95/P50 分布
3. 如果 P99 > 5s，检查：
   - Worker Manager → Active Workers 是否足够
   - Overview → Network Peers 是否正常
   - 系统资源（CPU/内存/磁盘）
```

---

## 📊 关键阈值参考

| 指标 | 正常范围 | 警告 | 严重 |
|------|---------|------|------|
| Engine State | 4 (Proving) | 2-3 | 0 (Stopped) |
| Active Provers | ≥ 1 | 0 | - |
| Time Since Last Frame | < 300s | 300-600s | > 600s |
| Active Workers | = CPU 核心数 | < 核心数 | 0 |
| Network Peers | ≥ 10 | 5-10 | < 5 |
| Frame Processing P99 | < 2s | 2-5s | > 5s |

---

## 🚦 颜色编码

Dashboard 使用统一的颜色编码：

- 🟢 **绿色** = 正常，一切正常
- 🟡 **黄色** = 警告，需要关注
- 🔴 **红色** = 严重，立即处理
- 🔵 **蓝色** = 信息，中间状态

---

## 🔧 配置建议

### 1. 数据刷新间隔

推荐设置：
- **Overview**: 10s (快速发现问题)
- **Node Monitoring**: 30s (详细分析)
- **Prover Registry**: 30s
- **Worker Manager**: 1m (变化不频繁)

### 2. 时间范围

- **日常巡检**: Last 1 hour
- **问题诊断**: Last 6 hours
- **性能分析**: Last 24 hours
- **趋势分析**: Last 7 days

### 3. 多节点变量

所有 Dashboard 支持 `$node` 变量：
- 选择 `All` 查看所有节点对比
- 选择特定节点进行单节点分析
- 支持多选进行对比

---

## 💡 最佳实践

### 1. 使用 Overview 作为主 Dashboard

- 设置为 Grafana 首页（Home Dashboard）
- 每天上班第一件事：检查 Overview
- 所有指标正常 = 可以放心工作

### 2. 设置浏览器书签

```
快速访问：
- 总览：http://192.168.1.6:3000/d/quilibrium-overview
- 节点详情：http://192.168.1.6:3000/d/quilibrium-node-monitoring
```

### 3. 建立巡检 Checklist

**每日巡检** (1 分钟):
- [ ] Engine State = Proving (绿色)
- [ ] Active Provers > 0 (绿色)
- [ ] Time Since Last Frame < 600s
- [ ] Frame Number 持续增长
- [ ] Network Peers > 5

**每周检查** (5 分钟):
- [ ] Frame Processing Duration P99 < 5s
- [ ] Prover Registry Refresh 成功率 > 95%
- [ ] Network Message Delivery Rate 稳定
- [ ] 无异常 Worker 数量变化

### 4. 告警集成

如果使用 Alertmanager，建议监控：
```yaml
Critical (立即处理):
- Engine State = 0 (Stopped)
- Time Since Last Frame > 600s
- Active Provers = 0

Warning (需要关注):
- Frame Processing P99 > 5s
- Network Peers < 5
- Active Workers < 预期值
```

---

## 🔍 问题诊断指南

### 问题 1: 节点不生成帧

**症状**: Time Since Last Frame 持续增长，超过 600s

**排查步骤**:
1. Overview → 检查 Engine State 是否 = 4 (Proving)
2. Worker Manager → 检查 Active Workers 是否正常
3. Node Monitoring → 检查 Network Peers 是否 > 5
4. 系统层面 → 检查 CPU/内存/磁盘

### 问题 2: Prover 数量为 0

**症状**: Active Provers = 0

**排查步骤**:
1. Prover Registry → 检查是否有 joining 状态的 Prover
2. Prover Registry → 查看 Registry Refresh 是否成功
3. 日志检查 → `grep "prover" ~/ceremonyclient/node/.logs/*`
4. 确认 ProverJoin 交易是否成功

### 问题 3: 网络 Peer 数量低

**症状**: Network Peers < 5

**排查步骤**:
1. Node Monitoring → 查看 Mesh Peer Counts 趋势
2. 检查防火墙配置
3. 检查网络连接质量
4. 查看 Blossomsub 指标（reject/duplicate message 率）

### 问题 4: 帧处理慢

**症状**: Frame Processing Duration P99 > 5s

**排查步骤**:
1. Worker Manager → 检查 Active Workers 数量
2. 系统 → 检查 CPU 使用率
3. 系统 → 检查磁盘 I/O
4. Node Monitoring → 查看是否所有节点都慢（网络问题）

---

## 📚 参考文档

- **完整监控指南**: `MONITORING_GUIDE.md`
- **Prover 指��集成**: `consensus/provers/METRICS_INTEGRATION.md`
- **Prometheus 配置**: `prometheus.yml`
- **Grafana 数据源**: `datasources.yml`

---

## 🆘 故障联系

如果遇到监控系统本身的问题：

1. **Dashboard 不显示数据**
   - 检查 Prometheus targets: http://192.168.1.6:9090/targets
   - 确认所有 target 都是 UP 状态
   - 检查防火墙 8080 端口

2. **指标缺失**
   - 访问 http://192.168.1.6:8080/metrics
   - 确认指标是否在输出中
   - 检查 node 是否正常运行

3. **Grafana 连接问题**
   - 检查 Docker 容器状态: `docker compose ps`
   - 查看 Grafana 日志: `docker compose logs grafana`
   - 重启服务: `docker compose restart grafana`

---

**最后更新**: 2025-10-25
**Dashboard 版本**: 2.0
**设计理念**: 运维优先，快速定位，易于诊断
