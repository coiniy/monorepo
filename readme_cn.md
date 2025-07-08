# Quilibrium Monorepo 中文介绍 (`readme_cn.md`)

本文档旨在提供对 Quilibrium Monorepo 项目的中文概览，帮助开发者快速理解其架构、功能和工作流程。

## 1. 项目概述

Quilibrium 是一个去中心化的、无需许可的协调协议。它旨在创建一个抗审查、保护隐私的通信和应用平台。此 `monorepo` 代码库包含了构建 Quilibrium 网络所需的所有核心组件，包括全功能节点、客户端工具以及底层的加密库。

项目采用混合语言开发模式：
- **Go**: 用于实现网络应用层，如 P2P 节点、客户端和网络通信逻辑。
- **Rust**: 用于实现核心加密原语，以追求极致的性能和内存安全。

## 2. 整体架构

Quilibrium 的架构可以分为以下几个层次：

- **应用层 (Go)**: 包括 `node` 和 `client`。`node` 是网络中的核心参与者，负责运行共识、处理数据和与其他节点通信。`client` 是一个与节点交互的命令行���面。
- **网络层 (Go)**: 基于 `go-libp2p` 构建，这是一个模块化的 P2P 网络堆栈。它使用 Kademlia DHT 进行节点发现，并使用 Blossomsub (一种 Gossip 协议) 进行高效的消息广播。
- **加密层 (Rust & Go)**: 核心加密算法由 Rust 实现，位于 `crates/` 目录下。这些库通过 C FFI (Foreign Function Interface) 接口暴露给 Go 代码，以便在上层应用中使用。Go 代码中也有部分加密逻辑的实现。
- **构建与部署**: 使用 `Taskfile.yaml` 管理复杂的构建和任务流程，并提供 `Dockerfile` 用于容器化部署。

## 3. 主要功能组件

### 3.1. 核心节点 (`node/`)
这是项目的主体部分，一个用 Go 编写的全功能 P2P 节点。其主要职责包括：
- **P2P 网络管理**: 初始化并管理 `libp2p` 主机，处理节点连接、发现和消息路由。
- **共识引擎**: 参与网络的共识算法，处理信标（Beacon）、时钟（Clock）和数据帧（Frames）。
- **数据存储**: 使用 PebbleDB (一个高性能的键值存储库) 持久化存储网络状态和数据。
- **RPC 服务**: 提供 gRPC 接口，允许 `client` 或其他应用与其交互。

### 3.2. 命令行客户端 (`client/`)
一个 Go 编写的命令行工具，用于和本地运行的 Quilibrium 节点通信。它提供了一系列命令，例如：
- 查询节点状态
- 获取网络信息
- 发送交易或执行特定操作

### 3.3. 加密核心库 (`crates/`)
这是项目安全和性能的基石，包含了一系列用 Rust 实现的高级密码学库：
- **VDF (`crates/vdf`)**: 可验证延迟函数 (Verifiable Delay Function)。这是共识机制的关键，提供了一个可公开验证的、需要一定时间才能完成的计算，以防止恶意行为者操纵时间或出块顺序。
- **BLS48581 (`crates/bls48581`)**: BLS 签名算法的实现，用于高效地聚合多个签名，减少共识消息的带宽占用。
- **类群 (`crates/classgroup`)**: 实现了类群算法，是 VDF 的核心数学基础。
- **其他库**: 还包括可验证加密 (`verenc`)、信道 (`channel`) 等多种密码学工具。

### 3.4. 网络层 (`go-libp2p`, `go-libp2p-blossomsub`, `go-libp2p-kad-dht`)
这些是 Go 实现的 libp2p 核心组件的分支或封装，为 Quilibrium 网络量身定制：
- **libp2p**: 提供了模块化的网络功能，包括传输、路由、加密连接等。
- **Kademlia DHT**: 一种分布式哈希表，用于在去中心化网络中高效地查找和发现其他节点。
- **Blossomsub**: 一种优化的 Gossip 协议，用于在网络中高效、可扩展地广播消息（例如交易和信标）。

## 4. 项目工作流程

一个典型的节点工作流程如下：

1.  **构建**: 开发者使用 `Taskfile.yaml` 中定义的命令 (`task build`) 来编译 Go 应用和底层的 Rust 库。
2.  **启动节点**: 运行编译好的 `node` 可执行文件。
3.  **初始化**: 节点加载配置文件，初始化 PebbleDB 数据库，并启动 gRPC 服务器。
4.  **启动网络**:
    - 节点创建一个 `libp2p` 主机实例。
    - 它配置并启动 Kademlia DHT 和 Blossomsub 服务。
    - 节点连接到代码中预设的引导节点 (Bootstrap Peers)。
5.  **节点发现与同步**:
    - 通过 DHT，节点开始发现网络中的其他对等节点，并建立连接。
    - 通过 Blossomsub，节点订阅感兴趣的消息主题（Topics），开始接收来自全网的实时消息。
6.  **参与共识**:
    - 节点根据内部时钟和接收到的网络信标，参与共识循环。
    - 它会执行 VDF 计算、验证收到的区块/数据帧、并广播自己的消息。
7.  **客户端交互**:
    - 用户在另一个终端中运行 `client` 程序。
    - `client` 通过 gRPC 连接到本地节点，发送请求（如 `status`）。
    - 节点处理请求并返回相应的数据。

## 5. 主要交互流程图

以下是使用 Mermaid 语法绘制的主要组件交互流程图。

```mermaid
graph TD
    subgraph User Interaction
        User_CLI[👤 User] -- ./client status --> Client[💻 Client CLI]
    end

    subgraph Local Machine
        Client -- gRPC Call --> Node[🚀 Node]
        Node -- Read/Write --> DB[(PebbleDB)]
    end

    subgraph P2P Network
        Node -- Initializes --> Libp2p[<-- Libp2p Host -->]
        Libp2p -- Uses --> KadDHT[🌐 Kad-DHT]
        Libp2p -- Uses --> Blossomsub[💬 Blossomsub Pub/Sub]
    end
    
    subgraph Core Logic
      Node -- Executes --> Consensus[⚙️ Consensus Engine]
      Consensus -- Calls --> Crypto[🔒 Rust Crypto Libs (VDF, BLS)]
    end

    subgraph Network Peers
        KadDHT -- Finds Peers --> OtherPeers[👥 Other Peers]
        Blossomsub -- Gossip Messages <--> OtherPeers
    end

    User_Run[👤 User] -- ./node run --> Node

    style User_CLI fill:#cde4ff
    style User_Run fill:#cde4ff
    style Client fill:#e1e1e1
    style Node fill:#d4edda
```
