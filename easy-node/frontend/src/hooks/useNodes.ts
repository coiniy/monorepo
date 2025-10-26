'use client';

import { useState, useEffect, useCallback } from 'react';
import { Node, NodeStatus, NodeHealthInfo, HealthStats, PerformanceMetrics } from '@/types/node';
import { nodeApi, nodeStatusApi, healthApi, monitorApi } from '@/lib/api';

export function useNodes() {
  const [nodes, setNodes] = useState<Node[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchNodes = useCallback(async () => {
    try {
      setLoading(true);
      const data = await nodeApi.getNodes();
      setNodes(data);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch nodes');
    } finally {
      setLoading(false);
    }
  }, []);

  const createNode = useCallback(async (data: { name: string; base_port: number }) => {
    try {
      const newNode = await nodeApi.createNode(data);
      setNodes(prev => [...prev, newNode]);
      return newNode;
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Failed to create node');
    }
  }, []);

  const updateNode = useCallback(async (id: number, data: Partial<Node>) => {
    try {
      const updatedNode = await nodeApi.updateNode(id, data);
      setNodes(prev => prev.map(node => node.id === id ? updatedNode : node));
      return updatedNode;
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Failed to update node');
    }
  }, []);

  const deleteNode = useCallback(async (id: number) => {
    try {
      await nodeApi.deleteNode(id);
      setNodes(prev => prev.filter(node => node.id !== id));
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Failed to delete node');
    }
  }, []);

  useEffect(() => {
    fetchNodes();
  }, [fetchNodes]);

  return {
    nodes,
    loading,
    error,
    fetchNodes,
    createNode,
    updateNode,
    deleteNode,
  };
}

export function useNodeStatus() {
  const [nodeStatus, setNodeStatus] = useState<NodeStatus[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchNodeStatus = useCallback(async () => {
    try {
      setLoading(true);
      const data = await nodeStatusApi.getNodesStatus();
      setNodeStatus(data.nodes);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch node status');
    } finally {
      setLoading(false);
    }
  }, []);

  const startNode = useCallback(async (id: number) => {
    try {
      await nodeStatusApi.startNode(id);
      await fetchNodeStatus(); // Refresh status
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Failed to start node');
    }
  }, [fetchNodeStatus]);

  const stopNode = useCallback(async (id: number) => {
    try {
      await nodeStatusApi.stopNode(id);
      await fetchNodeStatus(); // Refresh status
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Failed to stop node');
    }
  }, [fetchNodeStatus]);

  const restartNode = useCallback(async (id: number) => {
    try {
      await nodeStatusApi.restartNode(id);
      await fetchNodeStatus(); // Refresh status
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Failed to restart node');
    }
  }, [fetchNodeStatus]);

  useEffect(() => {
    fetchNodeStatus();
  }, [fetchNodeStatus]);

  return {
    nodeStatus,
    loading,
    error,
    fetchNodeStatus,
    startNode,
    stopNode,
    restartNode,
  };
}

export function useHealth() {
  const [healthData, setHealthData] = useState<NodeHealthInfo[]>([]);
  const [healthStats, setHealthStats] = useState<HealthStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchHealth = useCallback(async () => {
    try {
      setLoading(true);
      const [healthResponse, statsResponse] = await Promise.all([
        healthApi.getNodesHealth(),
        healthApi.getHealthStats(),
      ]);
      setHealthData(healthResponse.nodes);
      setHealthStats(statsResponse);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch health data');
    } finally {
      setLoading(false);
    }
  }, []);

  const checkAllHealth = useCallback(async () => {
    try {
      await healthApi.checkAllNodesHealth();
      await fetchHealth(); // Refresh data
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Failed to check health');
    }
  }, [fetchHealth]);

  useEffect(() => {
    fetchHealth();
  }, [fetchHealth]);

  return {
    healthData,
    healthStats,
    loading,
    error,
    fetchHealth,
    checkAllHealth,
  };
}

export function useMonitor() {
  const [metrics, setMetrics] = useState<PerformanceMetrics | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchMetrics = useCallback(async () => {
    try {
      setLoading(true);
      const data = await monitorApi.getPerformanceMetrics();
      setMetrics(data);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch metrics');
    } finally {
      setLoading(false);
    }
  }, []);

  const triggerHealthCheck = useCallback(async () => {
    try {
      await monitorApi.triggerHealthCheck();
      await fetchMetrics(); // Refresh metrics
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Failed to trigger health check');
    }
  }, [fetchMetrics]);

  const updateInterval = useCallback(async (interval: number) => {
    try {
      await monitorApi.updateInterval(interval);
      await fetchMetrics(); // Refresh metrics
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Failed to update interval');
    }
  }, [fetchMetrics]);

  useEffect(() => {
    fetchMetrics();
  }, [fetchMetrics]);

  return {
    metrics,
    loading,
    error,
    fetchMetrics,
    triggerHealthCheck,
    updateInterval,
  };
}