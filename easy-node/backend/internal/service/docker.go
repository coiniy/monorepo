package service

import (
	"easy-node-backend/internal/config"
	"easy-node-backend/internal/model"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"gorm.io/gorm"
)

type DockerService struct {
	db     *gorm.DB
	config *config.Config
}

func NewDockerService(db *gorm.DB, config *config.Config) *DockerService {
	return &DockerService{
		db:     db,
		config: config,
	}
}

func (s *DockerService) GenerateDockerCompose() error {
	var nodes []model.Node
	if err := s.db.Where("enabled = ?", true).Order("sort_no ASC").Find(&nodes).Error; err != nil {
		return fmt.Errorf("failed to get nodes: %w", err)
	}

	// 生成节点配置文件
	if err := s.generateNodeConfigs(nodes); err != nil {
		return fmt.Errorf("failed to generate node configs: %w", err)
	}

	tmpl := `version: '3.8'

services:
{{- range .Nodes}}
  node{{.SortNo}}:
    build: .
    container_name: quilibrium-node-{{.SortNo}}
    volumes:
      - ./node{{.SortNo}}-config:/root/.config
    ports:
      - "{{.BasePort}}:8340"
      - "{{add .BasePort 1}}:8336"
      - "{{add .BasePort 2}}:8337"
      - "{{add .BasePort 3}}:8338"
      - "{{add .BasePort 4}}:8080"
    restart: always
    deploy:
      resources:
        limits:
          cpus: '4'
          memory: 8G
        reservations:
          cpus: '2'
          memory: 4G
{{end}}
  envoy:
    image: envoyproxy/envoy:v1.31-latest
    container_name: quilibrium-envoy
    ports:
      - "8340:8340"
      - "8336:8336/udp"
      - "9901:9901"
    volumes:
      - ./envoy.yaml:/etc/envoy/envoy.yaml:ro
    restart: always
    depends_on:
{{- range .Nodes}}
      - node{{.SortNo}}
{{- end}}

networks:
  default:
    driver: bridge`

	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}

	t, err := template.New("docker-compose").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	composeFile := filepath.Join(s.config.Docker.DeploymentPath, "docker-compose.yml")
	file, err := os.Create(composeFile)
	if err != nil {
		return fmt.Errorf("failed to create docker-compose file: %w", err)
	}
	defer file.Close()

	data := struct {
		Nodes []model.Node
	}{
		Nodes: nodes,
	}

	if err := t.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// 生成 Envoy 配置文件
	if err := s.generateEnvoyConfig(nodes); err != nil {
		return fmt.Errorf("failed to generate Envoy config: %w", err)
	}

	// 生成 Prometheus 配置文件
	if err := s.generatePrometheusConfig(nodes); err != nil {
		return fmt.Errorf("failed to generate Prometheus config: %w", err)
	}

	return nil
}

func (s *DockerService) generateEnvoyConfig(nodes []model.Node) error {
	envoyTmpl := `admin:
  address:
    socket_address: { address: 0.0.0.0, port_value: 9901 }

static_resources:
  listeners:
  # TCP listener for port 8340
  - name: tcp_8340_listener
    address:
      socket_address: { address: 0.0.0.0, port_value: 8340 }
    filter_chains:
    - filters:
      - name: envoy.filters.network.tcp_proxy
        typed_config:
          "@type": type.googleapis.com/envoy.extensions.filters.network.tcp_proxy.v3.TcpProxy
          stat_prefix: tcp_8340
          cluster: nodes_8340_cluster

  # UDP listener for port 8336 (QUIC)
  - name: udp_8336_listener
    address:
      socket_address: { address: 0.0.0.0, port_value: 8336, protocol: UDP }
    listener_filters:
    - name: envoy.filters.udp_listener.udp_proxy
      typed_config:
        "@type": type.googleapis.com/envoy.extensions.filters.udp.udp_proxy.v3.UdpProxyConfig
        stat_prefix: udp_8336
        cluster: nodes_8336_cluster
        use_original_src_ip: false

  clusters:
  # TCP cluster for port 8340
  - name: nodes_8340_cluster
    connect_timeout: 5s
    type: ROUND_ROBIN
    lb_policy: ROUND_ROBIN
    load_assignment:
      cluster_name: nodes_8340_cluster
      endpoints:
      - lb_endpoints:
{{- range .Nodes}}
        - endpoint:
            address:
              socket_address: { address: node{{.SortNo}}, port_value: 8340 }
{{- end}}

  # UDP cluster for port 8336 (QUIC)
  - name: nodes_8336_cluster
    connect_timeout: 5s
    type: ROUND_ROBIN
    lb_policy: ROUND_ROBIN
    load_assignment:
      cluster_name: nodes_8336_cluster
      endpoints:
      - lb_endpoints:
{{- range .Nodes}}
        - endpoint:
            address:
              socket_address: { address: node{{.SortNo}}, port_value: 8336 }
{{- end}}`

	t, err := template.New("envoy").Parse(envoyTmpl)
	if err != nil {
		return fmt.Errorf("failed to parse Envoy template: %w", err)
	}

	// 获取部署目录
	envoyConfigPath := filepath.Join(s.config.Docker.DeploymentPath, "envoy.yaml")

	file, err := os.Create(envoyConfigPath)
	if err != nil {
		return fmt.Errorf("failed to create Envoy config file: %w", err)
	}
	defer file.Close()

	data := struct {
		Nodes []model.Node
	}{
		Nodes: nodes,
	}

	if err := t.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute Envoy template: %w", err)
	}

	return nil
}

func (s *DockerService) generateNodeConfigs(nodes []model.Node) error {
	for _, node := range nodes {
		// 创建节点配置目录 (使用 sortNo)
		nodeConfigDir := filepath.Join(s.config.Docker.DeploymentPath, fmt.Sprintf("node%d-config", node.SortNo))
		if err := os.MkdirAll(nodeConfigDir, 0755); err != nil {
			return fmt.Errorf("failed to create node config directory: %w", err)
		}

		// 生成 config.yml
		configPath := filepath.Join(nodeConfigDir, "config.yml")
		if err := s.writeFile(configPath, node.ConfigYml); err != nil {
			return fmt.Errorf("failed to write config.yml for node %d: %w", node.SortNo, err)
		}

		// 生成 keys.yml
		keysPath := filepath.Join(nodeConfigDir, "keys.yml")
		if err := s.writeFile(keysPath, node.KeysYml); err != nil {
			return fmt.Errorf("failed to write keys.yml for node %d: %w", node.SortNo, err)
		}
	}
	return nil
}

func (s *DockerService) writeFile(filePath, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

func (s *DockerService) generatePrometheusConfig(nodes []model.Node) error {
	promTmpl := `global:
  scrape_interval: 15s
  scrape_timeout: 10s
  evaluation_interval: 15s

alerting:
  alertmanagers:
  - scheme: http
    timeout: 10s
    api_version: v2
    static_configs:
    - targets: []

scrape_configs:
- job_name: prometheus
  honor_timestamps: true
  scrape_interval: 15s
  scrape_timeout: 10s
  metrics_path: /metrics
  scheme: http
  static_configs:
  - targets:
    - localhost:9090

- job_name: quilibrium_nodes
  honor_timestamps: true
  scrape_interval: 15s
  scrape_timeout: 10s
  metrics_path: /metrics
  scheme: http
  static_configs:
{{- range .Nodes}}
  - targets:
    - host.docker.internal:{{add .BasePort 4}}
    labels:
      node_name: "node-{{.SortNo}}"
      container_name: "quilibrium-node-{{.SortNo}}"
      sort_no: "{{.SortNo}}"
{{- end}}

- job_name: envoy_proxy
  honor_timestamps: true
  scrape_interval: 15s
  scrape_timeout: 10s
  metrics_path: /stats/prometheus
  scheme: http
  static_configs:
  - targets:
    - host.docker.internal:9901
    labels:
      service: "envoy-load-balancer"
`

	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}

	t, err := template.New("prometheus").Funcs(funcMap).Parse(promTmpl)
	if err != nil {
		return fmt.Errorf("failed to parse Prometheus template: %w", err)
	}

	// 获取 dashboards 目录
	prometheusConfigPath := filepath.Join(filepath.Dir(s.config.Docker.DeploymentPath), "dashboards", "prometheus.yml")

	file, err := os.Create(prometheusConfigPath)
	if err != nil {
		return fmt.Errorf("failed to create Prometheus config file: %w", err)
	}
	defer file.Close()

	data := struct {
		Nodes []model.Node
	}{
		Nodes: nodes,
	}

	if err := t.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute Prometheus template: %w", err)
	}

	return nil
}