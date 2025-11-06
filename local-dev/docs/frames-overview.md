# 全局帧与应用帧生成流程说明

本文档简要说明单机开发模式下 Quilibrium 节点生成区块（帧）的流程，助于定位日志并理解“【全局帧】”“【应用帧】”“【奖励】”等流程提示。

---

## 全局帧（Global Frame）

### 1. 统一收集消息
- 入口：`GlobalLivenessProvider.Collect`
- 从 Mixnet、挂起消息、奖励锁等来源汇总待打包数据。
- 校验合法性，写入内存并生成全局承诺（Commitment）。
- 日志包含：
  - `【全局帧】【收集】进入校验阶段`
  - `【全局帧】【收集】校验完成，生成承诺`
  - `【全局帧】【收集】已完成全局状态承诺`

### 2. 广播存活证明
- 入口：`GlobalLivenessProvider.SendLiveness`
- 将承诺哈希连同帧号签名后广播，表明本节点准备出块。
- 日志：`【全局帧】【存活】已广播存活签名`

### 3. 构建帧头
- 入口：`GlobalLeaderProvider.ProveNextState`
- 计算下一帧目标难度，生成帧头并打包业务消息。
- 日志包含：
  - `【全局帧】【选主】完成下一轮候选者排序`
  - `【全局帧】【构建】开始生成新帧头`
  - `【全局帧】【构建】已计算目标难度`
  - `【全局帧】【构建】将业务消息写入帧`
  - `【全局帧】【构建】新帧打包完成`

### 4. 广播提案
- 入口：`GlobalVotingProvider.SendProposal`
- 将帧序列化并广播，全局节点进入投票阶段。
- 日志：`【全局帧】【发布】准备广播帧提案`

### 5. 投票与聚合签名
- 入口：`GlobalVotingProvider.DecideAndSendVote` / `SendVote`
- 验证帧后投票，并在收集足够票数后聚合签名。
- 日志：
  - `【全局帧】【投票】已完成本地投票并广播`
  - `【全局帧】【落票】构建聚合签名完成`

### 6. 广播确认
- 入口：`GlobalVotingProvider.SendConfirmation`
- 带上聚合签名广播最终确认，并写入时间轴（Time Reel）。
- 日志：`【全局帧】【确认】已广播最终确认`

### 7. 挖矿奖励
- 入口：`ProverShardUpdate.applyReward`
- 根据帧难度、世界状态等计算奖励并写入全局奖励树。
- 日志：
  - `【奖励】计算挖矿奖励`
  - `【奖励】挖矿奖励已写入状态树`
  - 若当帧无奖励：`【奖励】本帧没有可分配奖励`

---

## 应用帧（App Frame）

应用共识流程与全局类似，但运行在各个应用（Shard）上，关键节点如下：

1. **收集消息**：`AppLivenessProvider.Collect`
   - 日志：`【应用帧】【收集】已完成消息筛选`

2. **广播存活证明**：`AppLivenessProvider.SendLiveness`
   - 日志：`【应用帧】【存活】已广播存活签名`

3. **构建帧**：`AppLeaderProvider.ProveNextState`
   - 日志：`【应用帧】【构建】开始生成新帧`
   - 含消息数、目标帧号等信息

4. **广播提案**：`AppVotingProvider.SendProposal`
   - 日志：`【应用帧】【发布】准备广播帧提案`

5. **投票与聚合**：`DecideAndSendVote`、`SendVote` 与全局类似
   - 日志：`【应用帧】【投票】已完成本地投票并广播`
   - 聚合完成：`【应用帧】【落票】聚合签名完成`

6. **广播确认**：`SendConfirmation`
   - 日志：`【应用帧】【确认】已广播最终确认`

---

## 查看日志建议
- 日志路径：`local-dev/logs/master.log`
- 推荐关键字：
  - 全局帧：`【全局帧】`
  - 应用帧：`【应用帧】`
  - 奖励：`【奖励】`
- 结合帧号可追踪完整的出块过程。

---

如需进一步排查，可配合 `state_machine` 的状态转换日志（English）与 `execution_manager` 输出一起分析。希望本文档有助于理解节点的出块流程。  
