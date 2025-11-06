## 单机环境 Prover Join 常见问题与解决方案

下面整理了在本地 `network == 99` 环境下，为了让数据工作节点成功 Join 时遇到的主要问题及对应排查方法。

### 1. JoinProof 请求超时 / `connection refused`

- **现象**：`master.log` 报错 `dial tcp 127.0.0.1:61000: connect: connection refused`。
- **原因**：数据工作线程尚未完成启动或 IPC 监听端口尚未暴露。
- **排查**：
  - 确认 `worker-*.log` 中是否出现 `data worker listening {"address": "...6100X"}`。
  - 如未启动，检查 `config.yml` 中 `engine.dataWorkerCount` 与本机资源是否匹配，必要时手动运行 `local-dev/build.sh`。
- **修复**：等待 worker 完全起来后，重新观察 Join 日志；若端口被占用，可修改 `dataWorkerBaseStreamPort` 并重新部署。

### 2. `existing seniority detected` 后依旧没有应用帧

- **现象**：Join 请求可广播，但 `message_count` 一直为 0。
- **原因**：`config/worker-store` 中已有旧的过滤器记录，导致 worker 被视为 Active。
- **排查与修复**：
  - 删除 `local-dev/config/store` 与 `local-dev/config/worker-store/*`，重新启动节点。
  - 确认新的 `worker-*.log` 有 “worker store not found, creating”。

### 3. `Join` 卡在 `status: 2`，`pending: false`

- **现象**：`event_distributor` 日志持续打印 `status: 2`、`pending=false`。
- **原因**：开发模式下 `ProverJoin` 默认将状态写成 Active。
- **修复**：修改 `global_prover_join.go` 让 dev 网络也从 `Status = 0 (Joining)` 开始，并增加调试日志，确认 `DecideWorkerJoins` 被调用。

### 4. `must wait 360 frames after join to confirm`

- **现象**：`validate message` 抛出需要等待 360 帧。
- **原因**：`ProverConfirm` 沿用了主网等待窗口。
- **修复**：在 dev 网络跳过 360 帧门槛，仅对主网帧执行等待逻辑。

### 5. `invalid allocation state for confirmation`

- **现象**：重复确认同一 Filter 时出现 “invalid allocation state”。
- **说明**：该错误表明分配已变为 Active，可忽略；确保后续不再重复发送 confirm。

### 6. 无应用帧 / 无奖励日志

- **原因**：全局层未收到应用交易。完成 Join 后，需要主动提交业务交易才会触发 “【应用帧】” 和 “【奖励】” 日志。
- **建议**：使用 mini-client（见 `local-dev/mini-client`）通过 RPC/REST 提交样例交易，观察 `message_count > 0` 与相应应用日志。

> 提示：结合 `【Join调试】` 日志可以快速确认 Join 是否已进入确认阶段；一旦确认成功，全局帧会开始接收业务负载。
