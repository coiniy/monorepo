export interface Node {
  id: number;
  name: string;
  status: string;
  base_port: number;
  enabled: boolean;
  peer_key: string;
  config_yml?: string;
  keys_yml?: string;
  health_status: 'unknown' | 'healthy' | 'unhealthy' | 'checking';
  last_health_check?: string;
  health_message?: string;
  response_time: number;
  created_at: string;
  updated_at: string;
}

export interface NodeStatus {
  id: number;
  name: string;
  container_name: string;
  container_status: string;
  base_port: number;
  enabled: boolean;
}

export interface NodeHealthInfo {
  node_id: number;
  node_name: string;
  container_status: string;
  health_status: 'unknown' | 'healthy' | 'unhealthy' | 'checking';
  last_check: string;
  response_time: number;
  health_message: string;
  metrics?: Record<string, any>;
}

export interface HealthStats {
  total_nodes: number;
  healthy_nodes: number;
  unhealthy_nodes: number;
  unknown_nodes: number;
  health_percentage: number;
  avg_response_time: number;
  last_check: string;
}

export interface PerformanceMetrics {
  overview: HealthStats;
  performance: {
    total_nodes: number;
    online_nodes: number;
    avg_response_time: number;
    max_response_time: number;
    min_response_time: number;
    uptime_percentage: number;
  };
  scheduler: {
    running: boolean;
    interval: string;
    next_check: string;
  };
  timestamp: string;
}

export interface CreateNodeRequest {
  name: string;
  base_port: number;
}

export interface UploadNodeRequest extends FormData {
  name: string;
  base_port: number;
  config: File;
  keys: File;
}