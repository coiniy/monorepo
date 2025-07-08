# 在 Quilibrium 上进行智能合约开发 (修正版)

本文档对之前错误的分析进行修正，并根据 Quilibrium 的真实架构，深入解读其独特的双层智能合约模型，解释**内在合约 (Intrinsic)** 与 **QCL 合约**的区别，并提供一个在 Quilibrium 上进行应用开发的指南。

## 1. Quilibrium 的双层智能合约模型

Quilibrium 采用了一种先进的混合模型来处理链上逻辑，而不是像以太坊那样仅依赖单一的虚拟机（VM）。这个模型分为两层：

### 1.1. 第一层：内在合约 (Intrinsic Contracts)

-   **定义**: 这些是协议的**核心功能**，直接用 **Go 语言**编写，并被编译进节点的核心二进制文件中。它们不是传统意义上的“智能合约”，而是协议原生的、高性能的执行引擎。
-   **代表**: `token_execution_engine.go` 是最完美的例子。它负责处理 Q 代币的所有逻辑（转账、余额、铸造等）。
-   **目的**: 为了追求极致的性能和安全性。核心的、需要高吞吐量的协议（如代币标准）被实现为原生代码，避免了虚拟机的性能开销。
-   **开发与部署**: 开发一个新的内在合约意味着直接修改 Quilibrium 节点的 Go 语言源码。部署则需要通过硬分叉（Hard Fork）进行全网的软件升级。这通常只用于协议级的核心功能开发。

### 1.2. 第二层：QCL (Quilibrium Contract Language) 合约

-   **定义**: 这才是面向普通开发者和用户的智能合约层。**QCL** 是 Quilibrium 设计的一种**自定义语言**，很可能是 Go 语言的一个安全的、确定性的子集。用户编写 QCL 代码来实现自定义的应用逻辑。
-   **执行环境**: QCL 代码不由 CPU 直接执行，而是在节点内部的一个**沙箱环境（Sandbox）或解释器（Interpreter）**中运行。这确保了用户提交的不可信代码无法破坏节点的稳定性和确定性。
-   **目的**: 提供灵活性和安全性。允许开发者在不改变核心协议的情况下，自由地部署和迭代各种去中心化应用。沙箱环境可以防止合约���行恶意操作或进入无限循环。
-   **状态交互**: 与内在合约一样，QCL 合约也通过与**超图 (Hypergraph)** 交互来读取和修改其状态。

**总结**: `token_execution_engine.go` 是一个**内在合约**，它为 QCL 合约提供了可以交互的核心功能（如代币操作）。而开发者想要创建自己的应用（例如一个 DAO、一个游戏），则需要使用 **QCL**。

## 2. 如何使用 QCL 编写和部署智能合约

基于对代码的分析，尤其是在 `node/schema/rdf.go` 中发现的线索，我们可以推断出 QCL 合约的开发和部署流程。

### 步骤 1: 使用 RDF 定义数据结构

Quilibrium 的合约开发似乎始于数据建模。开发者首先使用 **RDF (Resource Description Framework)** 词汇（可能是 Turtle 语法）来定义合约的数据结构。

-   **代码证据**: `rdf.go` 中的 `GenerateQCL` 函数接收一个文档并从中生成 QCL 代码。它会解析像 `qcl:size` 和 `qcl:order` 这样的谓词来确定数据结构。

这意味着，在编写逻辑之前，你需要先为你的应用定义一个清晰的“模式”（Schema）。

### 步骤 2: 编写 QCL 合约逻辑

虽然我们没有找到独立的 `.qcl` 文件作为直接示例，但可以推断 QCL 是一种受 Go 启发的、强类型的语言，其代码最终会被节点内的解释器执行。

下面是一个**推测性**的 QCL 示例代码，用于实现一个简单的投票合约。**请注意：此语法是基于对 Go 和内在合约模式的推断，可能与真实的 QCL 语法有出入。**

```go
// vote.qcl - A hypothetical QCL contract example

// Import necessary intrinsic contracts, like the standard library or token contract
import "intrinsics:token"
import "intrinsics:hypergraph"

// Define the contract state using structs, likely derived from an RDF schema
type VoteTopic struct {
    ID string
    Options map[string]uint64
}

type Vote struct {
    VoterID string
    Option  string
}

// @entrypoint
// Function to create a new voting topic
public func CreateTopic(topicID string, options []string) {
    // Create a hyperedge for the topic
    topicEdge := hypergraph.NewHyperedge(topicID)
    
    // Create vertices for each voting option and link them
    for _, opt := range options {
        optionVertex := hypergraph.NewVertex(opt, 0) // Store option name and initial count
        hypergraph.Link(optionVertex, topicEdge)
    }
}

// @entrypoint
// Function to cast a vote
public func CastVote(topicID string, option string, voterID string) {
    // Use an intrinsic to check if the voter holds a specific token (e.g., a governance token)
    // require(token.BalanceOf(voterID) > 0, "Voter has no voting power")

    // Find the hyperedge for the topic
    topicEdge := hypergraph.FindHyperedge(topicID)
    require(topicEdge != nil, "Topic not found")

    // Find the vertex for the chosen option
    optionVertex := hypergraph.FindVertex(topicEdge, option)
    require(optionVertex != nil, "Option not found")

    // Create a vertex representing this specific vote to prevent double-voting
    voteReceipt := hypergraph.NewVertex(voterID + "voted for" + option)
    require(hypergraph.IsLinked(voteReceipt, topicEdge) == false, "Already voted")

    // Increment the vote count on the option vertex
    currentCount := optionVertex.GetValue().(uint64)
    optionVertex.SetValue(currentCount + 1)

    // Link the receipt to the topic to mark this voter as having voted
    hypergraph.Link(voteReceipt, topicEdge)
}

// @readonly
// Function to get the results of a vote
public func GetResults(topicID string) map[string]uint64 {
    // ... logic to find the topic hyperedge and read the values from its linked option vertices ...
    return results
}
```

### 步骤 3: 部署和交互

1.  **部署**: 合约代码（可能是 RDF 结构定义或 QCL 源码）通过一笔特殊的交易被提交到网络上。节点接收到后，会将其存储在超图的特定部分。
2.  **交互**: 用户通过发送交易来调用合约的公共入口函数（如 `CreateTopic` 或 `CastVote`），并提供相应的参数。节点执行引擎会找到对应的合约代码，在沙箱中运行它，并根据执行结果更新超图状态。

这个模型兼顾了高性能（核心协议用原生 Go）和灵活性/安全性（用户应用用沙箱化的 QCL），代表了一种非常现代的区块链设计思路。