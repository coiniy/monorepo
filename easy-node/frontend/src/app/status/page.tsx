'use client';

import { useState } from 'react';
import { MainLayout } from '@/components/Layout/MainLayout';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Badge } from '@/components/ui/Badge';
import { useNodeStatus } from '@/hooks/useNodes';
import { dockerApi } from '@/lib/api';
import { getStatusColor } from '@/lib/utils';
import { 
  Play, 
  Square, 
  RotateCcw, 
  Server,
  PlayCircle,
  StopCircle,
  FileText,
  Activity
} from 'lucide-react';

export default function StatusPage() {
  const { nodeStatus, loading, error, fetchNodeStatus, startNode, stopNode, restartNode } = useNodeStatus();
  const [refreshing, setRefreshing] = useState(false);
  const [actionLoading, setActionLoading] = useState<{ [key: number]: string }>({});
  const [bulkLoading, setBulkLoading] = useState<string | null>(null);

  const handleRefresh = async () => {
    setRefreshing(true);
    try {
      await fetchNodeStatus();
    } finally {
      setRefreshing(false);
    }
  };

  const handleNodeAction = async (nodeId: number, action: 'start' | 'stop' | 'restart') => {
    setActionLoading(prev => ({ ...prev, [nodeId]: action }));
    try {
      switch (action) {
        case 'start':
          await startNode(nodeId);
          break;
        case 'stop':
          await stopNode(nodeId);
          break;
        case 'restart':
          await restartNode(nodeId);
          break;
      }
    } catch (err) {
      alert(`Failed to ${action} node: ${err instanceof Error ? err.message : 'Unknown error'}`);
    } finally {
      setActionLoading(prev => {
        const newState = { ...prev };
        delete newState[nodeId];
        return newState;
      });
    }
  };

  const handleBulkAction = async (action: 'start' | 'stop') => {
    setBulkLoading(action);
    try {
      if (action === 'start') {
        await dockerApi.startAllNodes();
      } else {
        await dockerApi.stopAllNodes();
      }
      await fetchNodeStatus(); // Refresh status after bulk action
    } catch (err) {
      alert(`Failed to ${action} all nodes: ${err instanceof Error ? err.message : 'Unknown error'}`);
    } finally {
      setBulkLoading(null);
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status.toLowerCase()) {
      case 'running': return <PlayCircle className="h-4 w-4 text-green-500" />;
      case 'stopped': 
      case 'exited': return <StopCircle className="h-4 w-4 text-red-500" />;
      default: return <Activity className="h-4 w-4 text-gray-400" />;
    }
  };

  const getStatusVariant = (status: string) => {
    switch (status.toLowerCase()) {
      case 'running': return 'success';
      case 'stopped':
      case 'exited': return 'danger';
      case 'starting': return 'warning';
      default: return 'secondary';
    }
  };

  const runningNodes = nodeStatus.filter(n => n.container_status === 'running').length;
  const stoppedNodes = nodeStatus.filter(n => n.container_status === 'stopped' || n.container_status === 'exited').length;

  return (
    <MainLayout
      title="Node Status Control"
      subtitle="Start, stop, and restart your Quilibrium nodes"
      onRefresh={handleRefresh}
      refreshing={refreshing}
      actions={
        <div className="flex space-x-2">
          <Button 
            variant="success" 
            onClick={() => handleBulkAction('start')}
            loading={bulkLoading === 'start'}
            disabled={bulkLoading !== null}
          >
            <PlayCircle className="h-4 w-4 mr-2" />
            Start All
          </Button>
          <Button 
            variant="destructive" 
            onClick={() => handleBulkAction('stop')}
            loading={bulkLoading === 'stop'}
            disabled={bulkLoading !== null}
          >
            <StopCircle className="h-4 w-4 mr-2" />
            Stop All
          </Button>
        </div>
      }
    >
      <div className="space-y-6">
        {/* Status Overview */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center space-x-2">
                <Server className="h-5 w-5 text-blue-600" />
                <div>
                  <div className="text-2xl font-bold">{nodeStatus.length}</div>
                  <div className="text-sm text-gray-600">Total Nodes</div>
                </div>
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center space-x-2">
                <PlayCircle className="h-5 w-5 text-green-600" />
                <div>
                  <div className="text-2xl font-bold text-green-600">{runningNodes}</div>
                  <div className="text-sm text-gray-600">Running</div>
                </div>
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center space-x-2">
                <StopCircle className="h-5 w-5 text-red-600" />
                <div>
                  <div className="text-2xl font-bold text-red-600">{stoppedNodes}</div>
                  <div className="text-sm text-gray-600">Stopped</div>
                </div>
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center space-x-2">
                <Activity className="h-5 w-5 text-blue-600" />
                <div>
                  <div className="text-2xl font-bold">
                    {nodeStatus.length > 0 ? Math.round((runningNodes / nodeStatus.length) * 100) : 0}%
                  </div>
                  <div className="text-sm text-gray-600">Uptime</div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Node Status Table */}
        <Card>
          <CardHeader>
            <CardTitle>Node Status Control</CardTitle>
          </CardHeader>
          <CardContent>
            {loading ? (
              <div className="flex items-center justify-center py-8">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
              </div>
            ) : error ? (
              <div className="text-center py-8 text-red-600">
                Error: {error}
              </div>
            ) : nodeStatus.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                <Server className="mx-auto h-12 w-12 text-gray-400 mb-4" />
                <h3 className="text-lg font-medium mb-2">No nodes found</h3>
                <p className="text-sm">No nodes available for status control.</p>
              </div>
            ) : (
              <div className="space-y-4">
                {nodeStatus.map((node) => (
                  <div key={node.id} className="border border-gray-200 rounded-lg p-4 hover:bg-gray-50 transition-colors">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-4">
                        <div className="flex-shrink-0">
                          <div className="h-12 w-12 rounded-lg bg-blue-100 flex items-center justify-center">
                            <Server className="h-6 w-6 text-blue-600" />
                          </div>
                        </div>
                        <div className="flex-1 min-w-0">
                          <div className="flex items-center space-x-3">
                            <h3 className="text-lg font-medium text-gray-900">{node.name}</h3>
                            <Badge variant={getStatusVariant(node.container_status)}>
                              <div className="flex items-center space-x-1">
                                {getStatusIcon(node.container_status)}
                                <span>{node.container_status}</span>
                              </div>
                            </Badge>
                            {!node.enabled && (
                              <Badge variant="secondary">Disabled</Badge>
                            )}
                          </div>
                          <div className="mt-1 flex items-center space-x-4 text-sm text-gray-500">
                            <span>ID: {node.id}</span>
                            <span>Port: {node.base_port}</span>
                            <span>Container: {node.container_name}</span>
                          </div>
                        </div>
                      </div>
                      
                      <div className="flex items-center space-x-2">
                        <Button
                          variant="success"
                          size="sm"
                          onClick={() => handleNodeAction(node.id, 'start')}
                          loading={actionLoading[node.id] === 'start'}
                          disabled={node.container_status === 'running' || !node.enabled || actionLoading[node.id] !== undefined}
                        >
                          <Play className="h-4 w-4 mr-1" />
                          Start
                        </Button>
                        
                        <Button
                          variant="destructive"
                          size="sm"
                          onClick={() => handleNodeAction(node.id, 'stop')}
                          loading={actionLoading[node.id] === 'stop'}
                          disabled={node.container_status === 'stopped' || node.container_status === 'exited' || !node.enabled || actionLoading[node.id] !== undefined}
                        >
                          <Square className="h-4 w-4 mr-1" />
                          Stop
                        </Button>
                        
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleNodeAction(node.id, 'restart')}
                          loading={actionLoading[node.id] === 'restart'}
                          disabled={!node.enabled || actionLoading[node.id] !== undefined}
                        >
                          <RotateCcw className="h-4 w-4 mr-1" />
                          Restart
                        </Button>
                        
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => {
                            // TODO: Implement logs modal
                            alert('Logs feature coming soon!');
                          }}
                        >
                          <FileText className="h-4 w-4 mr-1" />
                          Logs
                        </Button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Quick Actions */}
        <Card>
          <CardHeader>
            <CardTitle>Quick Actions</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <Button
                variant="success"
                onClick={() => handleBulkAction('start')}
                loading={bulkLoading === 'start'}
                disabled={bulkLoading !== null}
                className="h-16 flex flex-col items-center justify-center space-y-1"
              >
                <PlayCircle className="h-6 w-6" />
                <span>Start All Nodes</span>
              </Button>
              
              <Button
                variant="destructive"
                onClick={() => handleBulkAction('stop')}
                loading={bulkLoading === 'stop'}
                disabled={bulkLoading !== null}
                className="h-16 flex flex-col items-center justify-center space-y-1"
              >
                <StopCircle className="h-6 w-6" />
                <span>Stop All Nodes</span>
              </Button>
              
              <Button
                variant="outline"
                onClick={async () => {
                  try {
                    await dockerApi.regenerateDockerCompose();
                    alert('Docker compose regenerated successfully!');
                  } catch (err) {
                    alert(`Failed to regenerate: ${err instanceof Error ? err.message : 'Unknown error'}`);
                  }
                }}
                className="h-16 flex flex-col items-center justify-center space-y-1"
              >
                <RotateCcw className="h-6 w-6" />
                <span>Regenerate Config</span>
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </MainLayout>
  );
}