# Quilibrium Node Monitoring Guide

## 📊 概述

本指南介绍如何使用 Prometheus + Grafana 监控 Quilibrium 节点。系统提供了 **77 个 Quilibrium 指标** 和 **16 个 Blossomsub 网络指标**，全面覆盖节点的运行状态、性能和网络健康状况。

## 🚀 快速开始

### 访问地址

- **Grafana**: http://192.168.1.6:3000
- **Prometheus**: http://192.168.1.6:9090
- **Node Metrics**: http://192.168.1.6:8080/metrics

### 可用的 Dashboard

1. **Quilibrium Node Monitoring** - 节点核心监控
2. **Quilibrium Prover Registry** - Prover 注册和状态
3. **Quilibrium Worker Manager** - Worker 工作线程管理

---

## 📈 监控指标分类

### 1. 全局共识引擎 (Global Consensus Engine)

#### 引擎状态
```
quilibrium_global_consensus_engine_state
```
当前引擎状态：
- `0` = Stopped (已停止)
- `1` = Starting (启动中)
- `2` = Loading (加载中)
- `3` = Collecting (收集中)
- `4` = Proving (证明中)
- `5` = Publishing (发布中)
- `6` = Verifying (验证中)
- `7` = Stopping (停止中)

#### 帧处理指标
```
quilibrium_global_consensus_current_frame_number           # 当前帧号
quilibrium_global_consensus_current_difficulty             # 当前难度
quilibrium_global_consensus_time_since_last_proven_frame_seconds  # 距上次证明帧的时间
```

#### 性能指标 (Histograms)
```
quilibrium_global_consensus_frame_processing_duration_seconds  # 帧处理耗时
quilibrium_global_consensus_frame_proving_duration_seconds     # 帧证明耗时
quilibrium_global_consensus_frame_publishing_duration_seconds  # 帧发布耗时
quilibrium_global_consensus_frame_validation_duration_seconds  # 帧验证耗时
```

#### 全局协调
```
quilibrium_global_consensus_global_coordination_total            # 协调周期总数
quilibrium_global_consensus_global_coordination_duration_seconds # 协调耗时
quilibrium_global_consensus_state_summaries_aggregated          # 聚合的状态摘要数
quilibrium_global_consensus_shard_commitments_collected         # 收集的分片承诺数
quilibrium_global_consensus_shard_commitment_collection_duration_seconds
```

#### 提案、投票、活跃性检查
```
# 全局提案
quilibrium_global_consensus_proposal_processing_duration_seconds
quilibrium_global_consensus_proposal_validation_duration_seconds

# 全局投票
quilibrium_global_consensus_vote_processing_duration_seconds
quilibrium_global_consensus_vote_validation_duration_seconds

# 全局活跃性检查
quilibrium_global_consensus_liveness_check_processing_duration_seconds
quilibrium_global_consensus_liveness_check_validation_duration_seconds

# 全局确认
quilibrium_global_consensus_confirmation_processing_duration_seconds
quilibrium_global_consensus_confirmation_validation_duration_seconds
```

#### 分片操作
```
# 分片帧
quilibrium_global_consensus_shard_frame_processing_duration_seconds
quilibrium_global_consensus_shard_frame_validation_duration_seconds

# 分片提案
quilibrium_global_consensus_shard_proposal_processing_duration_seconds
quilibrium_global_consensus_shard_proposal_validation_duration_seconds

# 分片投票
quilibrium_global_consensus_shard_vote_processing_duration_seconds
quilibrium_global_consensus_shard_vote_validation_duration_seconds

# 分片活跃性检查
quilibrium_global_consensus_shard_liveness_check_processing_duration_seconds
quilibrium_global_consensus_shard_liveness_check_validation_duration_seconds

# 分片确认
quilibrium_global_consensus_shard_confirmation_processing_duration_seconds
quilibrium_global_consensus_shard_confirmation_validation_duration_seconds
```

#### 执行器管理
```
quilibrium_global_consensus_executors_registered           # 已注册的执行器数量
quilibrium_global_consensus_executor_registration_total    # 执行器注册总数 {action="register"}
```

---

### 2. Prover Registry（证明者注册表）

#### 核心指标
```
quilibrium_prover_registry_allocations_total        # Prover 分配总数
quilibrium_prover_registry_cached_provers_total     # 缓存的 Prover 总数
quilibrium_prover_registry_shard_tries_total        # Shard trie 总数
```

#### Prover 状态分布
```
quilibrium_prover_registry_provers_total{status="active"}    # 活跃的 Prover
quilibrium_prover_registry_provers_total{status="joining"}   # 加入中
quilibrium_prover_registry_provers_total{status="leaving"}   # 离开中
quilibrium_prover_registry_provers_total{status="paused"}    # 暂停
quilibrium_prover_registry_provers_total{status="left"}      # 已离开
```

**Status 说明**:
- `0` = joining (加入中) - Prover 正在注册加入网络
- `1` = leaving (离开中) - Prover 正在离开网络
- `2` = active (活跃) - Prover 正常工作
- `3` = paused (暂停) - Prover 已暂停
- `4` = left (已离开) - Prover 已退出

#### 操作指标
```
quilibrium_prover_registry_operations_total{operation="add"}      # 添加操作
quilibrium_prover_registry_operations_total{operation="remove"}   # 删除操作
quilibrium_prover_registry_operations_total{operation="update"}   # 更新操作
```

#### 性能指标
```
quilibrium_prover_registry_refresh_duration_seconds  # Registry 刷新耗时
quilibrium_prover_registry_refresh_total{status}     # 刷新次数 (success/error)
quilibrium_prover_registry_lookup_duration_seconds   # Prover 查询耗时
quilibrium_prover_registry_lookup_total{result}      # 查询次数 (found/not_found/error)
```

---

### 3. Worker Manager（工作线程管理）

```
quilibrium_worker_manager_active_workers        # 活跃的 Worker 数量（对应 CPU 核心）
quilibrium_worker_manager_allocated_workers     # 已分配的 Worker 数量
quilibrium_worker_manager_total_storage_bytes   # Worker 总存储容量
quilibrium_worker_manager_operation_duration_seconds{operation}  # 操作耗时
```

**重要说明**:
- `active_workers` 对应你的 CPU 核心数（例如 32 核 = 32 workers）
- `allocated_workers` 是当前分配给特定任务的 worker 数
- Worker 数量 ≠ Prover 数量（Worker 是计算资源，Prover 是网络角色）

---

### 4. Hypergraph（超图状态管理）

#### 基础指标
```
quilibrium_hypergraph_size_total                 # 超图总大小
quilibrium_hypergraph_import_tree_size          # 导入树大小
quilibrium_hypergraph_import_tree_total         # 树导入总数
```

#### 操作性能
```
# 顶点操作
quilibrium_hypergraph_add_vertex_duration_seconds
quilibrium_hypergraph_remove_vertex_duration_seconds
quilibrium_hypergraph_vertex_adds_shards
quilibrium_hypergraph_vertex_removes_shards

# 超边操作
quilibrium_hypergraph_add_hyperedge_duration_seconds
quilibrium_hypergraph_remove_hyperedge_duration_seconds
quilibrium_hypergraph_hyperedge_adds_shards
quilibrium_hypergraph_hyperedge_removes_shards

# 其他操作
quilibrium_hypergraph_commit_duration_seconds
quilibrium_hypergraph_transaction_duration_seconds
quilibrium_hypergraph_transaction_total
quilibrium_hypergraph_import_tree_duration_seconds
quilibrium_hypergraph_within_operation_duration_seconds
quilibrium_hypergraph_within_operation_total
```

#### 数据管理
```
quilibrium_hypergraph_vertex_data_pruning_duration_seconds  # 顶点数据修剪耗时
quilibrium_hypergraph_traversal_proof_keys_per_request      # 遍历证明每请求密钥数
```

---

### 5. Event Distributor（事件分发器）

```
quilibrium_event_distributor_starts_total           # 分发器启动总数
quilibrium_event_distributor_subscribers_count      # 当前活跃订阅者数
quilibrium_event_distributor_subscriptions_total    # 创建的订阅总数
quilibrium_event_distributor_uptime_seconds         # 分发器运行时间
```

所有指标都带有 `{distributor_type="global"}` 标签。

---

### 6. Time Reel（时间轴）

```
quilibrium_time_reel_tree_depth              # 时间轴树深度
quilibrium_time_reel_tree_node_count         # 时间轴树节点数
quilibrium_time_reel_pending_frames_count    # 等待父帧的帧数
quilibrium_time_reel_equivocators_tracked    # 正在跟踪的双重签名者数
```

所有指标都带有 `{reel_type="global"}` 标签。

---

### 7. Dynamic Fees（动态费用）

```
quilibrium_dynamic_fees_filters_tracked       # 当前跟踪的过滤器数
quilibrium_dynamic_fees_filters_pruned_total  # 因不活跃被修剪的过滤器总数
```

---

### 8. Blossomsub 网络（P2P 网络层）

#### Peer 管理
```
blossomsub_add_peer_total         # 添加到 mesh 的 peer 总数
blossomsub_remove_peer_total      # 从 mesh 移除的 peer 总数
blossomsub_throttle_peer_total    # 被限流的 peer 总数
```

#### Mesh 拓扑
```
blossomsub_join_total             # 加入 mesh 总数
blossomsub_graft_total            # Graft 总数（添加连接）
blossomsub_prune_total            # Prune 总数（移除连接）
blossomsub_mesh_peer_counts       # 每个 bitmask 的 mesh peer 数量
```

#### 消息传递
```
blossomsub_deliver_message_total   # 投递的消息总数
blossomsub_duplicate_message_total # 重复消息总数
blossomsub_reject_message_total    # 拒绝的消息总数
blossomsub_validate_message_total  # 验证的消息总数
```

#### RPC 通信
```
blossomsub_recv_rpc_total          # 接收的 RPC 总数
blossomsub_send_rpc_total          # 发送的 RPC 总数
blossomsub_drop_rpc_total          # 丢弃的 RPC 总数
```

#### Gossipsub 协议消息
```
blossomsub_ihave_messages          # IHave 消息中的消息数分布（Histogram）
blossomsub_iwant_messages          # IWant 消息中的消息数分布（Histogram）
blossomsub_idontwant_messages      # IDontWant 消息中的消息数分布（Histogram）
```

所有 Blossomsub 指标都带有 `{protocol}` 标签，例如：
- `{protocol="/blossomsub/2.0.0"}`
- `{protocol="/blossomsub/2.1.0"}`

---

## 🎯 关键监控场景

### 场景 1: 节点健康检查

**关键指标**:
```promql
# 引擎状态（应该是 4=Proving）
quilibrium_global_consensus_engine_state

# 当前帧号（应该持续增长）
quilibrium_global_consensus_current_frame_number

# 距上次证明帧时间（不应过长）
quilibrium_global_consensus_time_since_last_proven_frame_seconds < 600

# 活跃 Worker 数量（应该等于 CPU 核心数）
quilibrium_worker_manager_active_workers == 32

# 活跃 Prover 数量（应该 > 0）
quilibrium_prover_registry_provers_total{status="active"} > 0
```

### 场景 2: 性能监控

**帧处理性能**:
```promql
# P99 帧处理时间
histogram_quantile(0.99, rate(quilibrium_global_consensus_frame_processing_duration_seconds_bucket[5m]))

# P99 帧证明时间
histogram_quantile(0.99, rate(quilibrium_global_consensus_frame_proving_duration_seconds_bucket[5m]))

# 平均帧处理速率
rate(quilibrium_global_consensus_frame_processing_duration_seconds_count[5m])
```

**Prover Registry 性能**:
```promql
# Registry 刷新成功率
rate(quilibrium_prover_registry_refresh_total{status="success"}[5m])
/
rate(quilibrium_prover_registry_refresh_total[5m])

# P95 刷新时间
histogram_quantile(0.95, rate(quilibrium_prover_registry_refresh_duration_seconds_bucket[5m]))
```

### 场景 3: 网络健康

**网络连接**:
```promql
# Mesh peer 数量
blossomsub_mesh_peer_counts

# 消息投递速率
rate(blossomsub_deliver_message_total[5m])

# 消息丢弃率（应该很低）
rate(blossomsub_reject_message_total[5m])
rate(blossomsub_duplicate_message_total[5m])
```

**网络吞吐**:
```promql
# RPC 接收速率
rate(blossomsub_recv_rpc_total[5m])

# RPC 发送速率
rate(blossomsub_send_rpc_total[5m])
```

### 场景 4: Prover 状态监控

**Prover 分布**:
```promql
# 各状态 Prover 数量
quilibrium_prover_registry_provers_total{status="active"}
quilibrium_prover_registry_provers_total{status="joining"}
quilibrium_prover_registry_provers_total{status="leaving"}

# 总分配数
quilibrium_prover_registry_allocations_total

# Prover 操作速率
rate(quilibrium_prover_registry_operations_total{operation="add"}[5m])
```

### 场景 5: 执行器监控

```promql
# 已注册执行器数量
quilibrium_global_consensus_executors_registered

# 执行器注册速率
rate(quilibrium_global_consensus_executor_registration_total[5m])
```

---

## 🔍 故障排查

### 问题 1: 节点不生成帧

**检查指标**:
```promql
quilibrium_global_consensus_engine_state != 4  # 引擎不在 Proving 状态
quilibrium_global_consensus_time_since_last_proven_frame_seconds > 600  # 超过 10 分钟未证明
```

**可能原因**:
- 引擎未进入 Proving 状态
- CPU 资源不足
- 网络连接问题

### 问题 2: Prover 注册失败

**检查指标**:
```promql
quilibrium_prover_registry_provers_total{status="active"} == 0  # 没有活跃 Prover
quilibrium_prover_registry_provers_total{status="joining"} > 0  # 卡在 joining 状态
rate(quilibrium_prover_registry_refresh_total{status="error"}[5m]) > 0  # Registry 刷新错误
```

**可能原因**:
- ProverJoin 交易未成功
- Registry 刷新失败
- 网络连接问题

### 问题 3: 网络分区

**检查指标**:
```promql
blossomsub_mesh_peer_counts < 5  # Mesh peer 数量过低
rate(blossomsub_deliver_message_total[5m]) == 0  # 无消息投递
rate(blossomsub_drop_rpc_total[5m]) > 0  # RPC 丢弃增加
```

**可能原因**:
- 防火墙阻止连接
- 网络延迟过高
- Peer 发现问题

### 问题 4: 性能下降

**检查指标**:
```promql
histogram_quantile(0.99, rate(quilibrium_global_consensus_frame_processing_duration_seconds_bucket[5m])) > 5  # P99 超过 5 秒
quilibrium_worker_manager_active_workers < 32  # Worker 数量减少
```

**可能原因**:
- CPU 过载
- 磁盘 I/O 瓶颈
- 内存不足

---

## 📊 告警规则建议

### 关键告警

```yaml
groups:
  - name: quilibrium_critical
    rules:
      # 引擎停止
      - alert: ConsensusEngineStopped
        expr: quilibrium_global_consensus_engine_state == 0
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Consensus engine is stopped"

      # 长时间未证明帧
      - alert: NoFrameProvenRecently
        expr: quilibrium_global_consensus_time_since_last_proven_frame_seconds > 600
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "No frame proven in last 10 minutes"

      # 没有活跃 Prover
      - alert: NoActiveProvers
        expr: quilibrium_prover_registry_provers_total{status="active"} == 0
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "No active provers registered"

      # Worker 数量异常
      - alert: WorkerCountLow
        expr: quilibrium_worker_manager_active_workers < 32
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Active worker count is below expected ({{$value}} < 32)"
```

### 性能告警

```yaml
  - name: quilibrium_performance
    rules:
      # 帧处理耗时过长
      - alert: SlowFrameProcessing
        expr: histogram_quantile(0.99, rate(quilibrium_global_consensus_frame_processing_duration_seconds_bucket[5m])) > 10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "P99 frame processing duration is high: {{$value}}s"

      # Registry 刷新失败
      - alert: RegistryRefreshFailing
        expr: rate(quilibrium_prover_registry_refresh_total{status="error"}[5m]) > 0
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Prover registry refresh errors detected"
```

### 网络告警

```yaml
  - name: quilibrium_network
    rules:
      # Mesh peer 数量低
      - alert: LowMeshPeerCount
        expr: blossomsub_mesh_peer_counts < 5
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Mesh peer count is low: {{$value}}"

      # 消息拒绝率高
      - alert: HighMessageRejectRate
        expr: rate(blossomsub_reject_message_total[5m]) > 10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High message reject rate: {{$value}}/s"
```

---

## 📝 最佳实践

### 1. 监控 Dashboard 使用

- **Node Monitoring**: 每天检查，关注引擎状态和帧号增长
- **Prover Registry**: 重点监控 Active Prover 数量和 Registry 刷新成功率
- **Worker Manager**: 确保 Active Workers 等于 CPU 核心数

### 2. 数据保留策略

推荐 Prometheus 保留策略：
```yaml
storage.tsdb.retention.time: 15d  # 保留 15 天
storage.tsdb.retention.size: 50GB # 或 50GB 上限
```

### 3. 多节点监控

在 `prometheus.yml` 中配置多个节点：
```yaml
scrape_configs:
  - job_name: quilibrium_node
    static_configs:
      - targets: ['192.168.1.6:8080']
        labels:
          node_name: "node-1"
      - targets: ['192.168.1.149:8080']
        labels:
          node_name: "node-2"
      - targets: ['192.168.1.36:8080']
        labels:
          node_name: "node-3"
```

### 4. 查询优化

- 使用 `rate()` 而不是 `irate()` 用于告警
- 使用适当的时间窗口（通常 5m）
- 使用 `histogram_quantile()` 分析延迟分布
- 按 label 聚合时使用 `by` 或 `without`

---

## 🔧 配置文件参考

### Prometheus 配置
位置: `/root/git/monorepo/node/dashboards/prometheus.yml`

### Grafana Datasource
位置: `/root/git/monorepo/node/dashboards/datasources.yml`

### Dashboard 配置
位置: `/root/git/monorepo/node/dashboards/*.json`

---

## 📚 参考资源

- **Prometheus 文档**: https://prometheus.io/docs/
- **Grafana 文档**: https://grafana.com/docs/
- **PromQL 查询**: https://prometheus.io/docs/prometheus/latest/querying/basics/

---

## 🆘 支持

如遇问题，请检查：
1. Prometheus targets 是否都是 UP 状态 (http://localhost:9090/targets)
2. Node metrics 端点是否可访问 (http://localhost:8080/metrics)
3. Grafana datasource 是否配置正确
4. 防火墙是否开放 8080 端口

---

**最后更新**: 2025-10-25
**版本**: 2.1.0
**指标总数**: 93 个（77 个 Quilibrium + 16 个 Blossomsub）
