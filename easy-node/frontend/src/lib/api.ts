import { Node, NodeStatus, NodeHealthInfo, HealthStats, PerformanceMetrics, CreateNodeRequest } from '@/types/node';

const API_BASE = '/api/v1';

// Helper function for API calls
async function apiCall<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE}${endpoint}`, {
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
    ...options,
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Network error' }));
    throw new Error(error.error || `HTTP ${response.status}`);
  }

  return response.json();
}

// Node Management APIs
export const nodeApi = {
  // Get all nodes
  getNodes: (): Promise<Node[]> => 
    apiCall('/nodes'),

  // Create node
  createNode: (data: CreateNodeRequest): Promise<Node> =>
    apiCall('/nodes', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  // Update node
  updateNode: (id: number, data: Partial<Node>): Promise<Node> =>
    apiCall(`/nodes/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  // Delete node
  deleteNode: (id: number): Promise<{ message: string }> =>
    apiCall(`/nodes/${id}`, {
      method: 'DELETE',
    }),

  // Upload node config
  uploadNodeConfig: (formData: FormData): Promise<Node> =>
    fetch(`${API_BASE}/nodes/upload`, {
      method: 'POST',
      body: formData,
    }).then(async (response) => {
      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Upload failed');
      }
      return response.json();
    }),

  // Batch upload nodes
  batchUploadNodes: (formData: FormData): Promise<{
    message: string;
    total: number;
    success_count: number;
    error_count: number;
    results: any[];
  }> =>
    fetch(`${API_BASE}/nodes/batch-upload`, {
      method: 'POST',
      body: formData,
    }).then(async (response) => {
      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Batch upload failed');
      }
      return response.json();
    }),
};

// Node Status Management APIs
export const nodeStatusApi = {
  // Get all node status
  getNodesStatus: (): Promise<{ nodes: NodeStatus[] }> =>
    apiCall('/nodes/status'),

  // Start node
  startNode: (id: number): Promise<{ message: string }> =>
    apiCall(`/nodes/${id}/start`, {
      method: 'POST',
    }),

  // Stop node
  stopNode: (id: number): Promise<{ message: string }> =>
    apiCall(`/nodes/${id}/stop`, {
      method: 'POST',
    }),

  // Restart node
  restartNode: (id: number): Promise<{ message: string }> =>
    apiCall(`/nodes/${id}/restart`, {
      method: 'POST',
    }),

  // Get node logs
  getNodeLogs: (id: number, lines?: number): Promise<{
    node_name: string;
    lines: number;
    logs: string;
  }> => {
    const params = lines ? `?lines=${lines}` : '';
    return apiCall(`/nodes/${id}/logs${params}`);
  },
};

// Health Monitoring APIs
export const healthApi = {
  // Get all nodes health
  getNodesHealth: (): Promise<{ nodes: NodeHealthInfo[] }> =>
    apiCall('/nodes/health'),

  // Get health stats
  getHealthStats: (): Promise<HealthStats> =>
    apiCall('/nodes/health/stats'),

  // Check single node health
  checkNodeHealth: (id: number): Promise<NodeHealthInfo> =>
    apiCall(`/nodes/${id}/health`),

  // Check all nodes health
  checkAllNodesHealth: (): Promise<{ nodes: NodeHealthInfo[]; total: number }> =>
    apiCall('/nodes/health/check', {
      method: 'POST',
    }),
};

// Monitor Management APIs
export const monitorApi = {
  // Get scheduler status
  getSchedulerStatus: (): Promise<{
    running: boolean;
    interval: string;
    next_check: string;
  }> =>
    apiCall('/monitor/scheduler/status'),

  // Start scheduler
  startScheduler: (): Promise<{ message: string; status: any }> =>
    apiCall('/monitor/scheduler/start', {
      method: 'POST',
    }),

  // Stop scheduler
  stopScheduler: (): Promise<{ message: string }> =>
    apiCall('/monitor/scheduler/stop', {
      method: 'POST',
    }),

  // Update interval
  updateInterval: (interval: number): Promise<{ message: string; new_interval: number }> =>
    apiCall('/monitor/scheduler/interval', {
      method: 'PUT',
      body: JSON.stringify({ interval }),
    }),

  // Trigger health check
  triggerHealthCheck: (): Promise<{ message: string }> =>
    apiCall('/monitor/health/trigger', {
      method: 'POST',
    }),

  // Get performance metrics
  getPerformanceMetrics: (): Promise<PerformanceMetrics> =>
    apiCall('/monitor/metrics'),

  // Get trend data
  getTrendData: (hours?: number): Promise<{
    time_range: number;
    data_points: any[];
    message: string;
  }> => {
    const params = hours ? `?hours=${hours}` : '';
    return apiCall(`/monitor/trends${params}`);
  },
};

// Docker Management APIs
export const dockerApi = {
  // Regenerate docker compose
  regenerateDockerCompose: (): Promise<{ message: string }> =>
    apiCall('/docker/regenerate', {
      method: 'POST',
    }),

  // Start all nodes
  startAllNodes: (): Promise<{ message: string }> =>
    apiCall('/docker/start-all', {
      method: 'POST',
    }),

  // Stop all nodes
  stopAllNodes: (): Promise<{ message: string }> =>
    apiCall('/docker/stop-all', {
      method: 'POST',
    }),
};