package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
)

// GeneratePeerPrivKey 生成私钥用于 P2P
// 注意: Quilibrium 使用 Ed448,但为了简化我们使用 Ed25519 生成随机密钥
// 如果需要完全兼容,需要添加 github.com/cloudflare/circl 依赖
func GeneratePeerPrivKey() (string, error) {
	// 生成 114 字节的随机密钥 (Ed448 私钥长度)
	hostKey := make([]byte, 114)
	if _, err := rand.Read(hostKey); err != nil {
		return "", errors.Wrap(err, "failed to generate host key")
	}

	// 转换为十六进制字符串
	return hex.EncodeToString(hostKey), nil
}

// GenerateEncryptionKey 生成 32 字节的加密密钥用于 keystore
func GenerateEncryptionKey() (string, error) {
	keystoreKey := make([]byte, 32)
	if _, err := rand.Read(keystoreKey); err != nil {
		return "", errors.Wrap(err, "failed to generate keystore key")
	}

	return hex.EncodeToString(keystoreKey), nil
}

// GenerateProvingKey 生成 proving key (可选)
func GenerateProvingKey(keystoreKey []byte) (string, string, error) {
	// 生成 Ed25519 proving key
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to generate proving key")
	}

	provingKey := privateKey.Seed()

	// 加密 private key
	iv := make([]byte, 12)
	if _, err := rand.Read(iv); err != nil {
		return "", "", errors.Wrap(err, "failed to generate IV")
	}

	aesCipher, err := aes.NewCipher(keystoreKey)
	if err != nil {
		return "", "", errors.Wrap(err, "could not construct cipher")
	}

	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return "", "", errors.Wrap(err, "could not construct block")
	}

	ciphertext := gcm.Seal(nil, iv, provingKey, nil)
	ciphertext = append(iv, ciphertext...)

	return hex.EncodeToString(ciphertext), hex.EncodeToString(publicKey), nil
}

// GenerateDefaultKeysYml 生成默认的 keys.yml 内容
func GenerateDefaultKeysYml() string {
	return "null:\n"
}

// GenerateDefaultConfigYml 生成默认的 config.yml 内容
func GenerateDefaultConfigYml(peerPrivKey, encryptionKey string) string {
	return fmt.Sprintf(`alias: null
key:
  keyManagerType: file
  keyManagerFile:
    path: .config/keys.yml
    createIfMissing: false
    encryptionKey: %s
p2p:
  d: 0
  dLo: 0
  dHi: 0
  dScore: 0
  dOut: 0
  historyLength: 0
  historyGossip: 0
  dLazy: 0
  gossipFactor: 0
  gossipRetransmission: 0
  heartbeatInitialDelay: 0s
  heartbeatInterval: 0s
  fanoutTTL: 0s
  prunePeers: 0
  pruneBackoff: 0s
  unsubscribeBackoff: 0s
  connectors: 0
  maxPendingConnections: 0
  connectionTimeout: 0s
  directConnectTicks: 0
  directConnectInitialDelay: 0s
  opportunisticGraftTicks: 0
  opportunisticGraftPeers: 0
  graftFloodThreshold: 0s
  maxIHaveLength: 0
  maxIHaveMessages: 0
  maxIDontWantMessages: 0
  iWantFollowupTime: 0s
  iDontWantMessageThreshold: 0
  iDontWantMessageTTL: 0
  bootstrapPeers:
  - /dns/bootstrap.quilibrium.com/udp/8336/quic-v1/p2p/Qme3g6rJWuz8HVXxpDb7aV2hiFq8bZJNqxMmwzmASzfq1M
  - /dns/quecifer.quilibrium.com/udp/8336/quic-v1/p2p/QmdWF9bGTH5mwJXkxrG859HA5r34MxXtMSTuEikSMDSESv
  - /dns/quagmire.quilibrium.com/udp/8336/quic-v1/p2p/QmaQ9KAaKtqXhYSQ5ARQNnn8B8474cWGvvD6PgJ4gAtMrx
  - /ip4/65.109.17.13/udp/8336/quic-v1/p2p/Qmc35n99eojSvW3PkbfBczJoSX92WmnnKh3Fg114ok3oo4
  - /ip4/65.108.194.84/udp/8336/quic-v1/p2p/QmP8C7g9ZRiWzhqN2AgFu5onS6HwHzR6Vv1TCHxAhnCSnq
  - /ip4/15.204.100.222/udp/8336/quic-v1/p2p/Qmef3Z3RvGg49ZpDPcf2shWtJNgPJNpXrowjUcfz23YQ3V
  listenMultiaddr: /ip4/0.0.0.0/udp/8336/quic-v1
  streamListenMultiaddr: ""
  peerPrivKey: %s
  traceLogFile: ""
  traceLogStdout: false
  network: 0
  lowWatermarkConnections: 0
  highWatermarkConnections: 0
  directPeers: []
  grpcServerRateLimit: 0
  minBootstrapPeers: 0
  bootstrapParallelism: 0
  discoveryParallelism: 0
  discoveryPeerLookupLimit: 0
  pingTimeout: 0s
  pingPeriod: 0s
  pingAttempts: 0
  validateQueueSize: 0
  validateWorkers: 0
  subscriptionQueueSize: 0
  peerOutboundQueueSize: 0
engine:
  provingKeyId: default-proving-key
  filter: ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
  genesisSeed: 726573697374206d7563682c206f626579206c6974746c657c083fb0a4274b1f70e9aa2b3f
  maxFrames: -1
  pendingCommitWorkers: 4
  minimumPeersRequired: 0
  statsMultiaddr: ""
  dataWorkerBaseListenMultiaddr: ""
  dataWorkerBaseP2PPort: 0
  dataWorkerBaseStreamPort: 0
  dataWorkerMemoryLimit: 0
  dataWorkerP2PMultiaddrs: []
  dataWorkerStreamMultiaddrs: []
  dataWorkerCount: 0
  dataWorkerFilters: []
  multisigProverEnrollmentPaths: []
  syncTimeout: 0s
  syncCandidates: 0
  syncMessageLimits:
    maxRecvMsgSize: 0
    maxSendMsgSize: 0
  enableMasterProxy: true
  rewardStrategy: ""
  archiveMode: false
  delegateAddress: ""
  allowExcessiveGOMAXPROCS: false
  blacklist: []
  alertKey: ""
  difficulty: 0
  rebuildstart: ""
  rebuildend: ""
  framePublish:
    mode: ""
    threshold: 0
    fragmentation:
      algorithm: ""
      reedSolomon:
        dataShards: 0
        parityShards: 0
    ballastSize: 0
db:
  path: .config/store
  workerPathPrefix: .config/worker-store/%%d
  workerPaths: []
  noticePercentage: 0
  warnPercentage: 0
  terminatePercentage: 0
  inmemorydonotuse: false

listenGrpcMultiaddr: ""
listenRESTMultiaddr: ""
logFile: ""
logger:
  path: .logs
  maxSize: 50
  maxBackups: 5
  maxAge: 10
  compress: true
`, encryptionKey, peerPrivKey)
}

// StasisSeed 是创世区块种子 (与 Quilibrium 节点相同)
const StasisSeed = "5GD3xAkCUJ8NREkiHJ6FUVFvFi8WJpDvFqVyBqZqRwJqYiW9GVzN2R1D1VY9KNtF8Q3q2qUZvWPHzFfAqYf2zR9P5hJc9MZZMhKvLX9Y3NHPnRpMQ8zg1vBt2dBhKRvP"
