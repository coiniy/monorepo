'use client';

import { useState } from 'react';
import { MainLayout } from '@/components/Layout/MainLayout';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Badge } from '@/components/ui/Badge';
import { useNodes } from '@/hooks/useNodes';
import { formatRelativeTime, getPortRange } from '@/lib/utils';
import { 
  Plus, 
  Edit, 
  Trash2, 
  Server,
  Eye,
  EyeOff,
  Copy,
  CheckCircle2
} from 'lucide-react';
import { Node } from '@/types/node';

export default function NodesPage() {
  const { nodes, loading, error, fetchNodes, deleteNode } = useNodes();
  const [refreshing, setRefreshing] = useState(false);
  const [selectedNode, setSelectedNode] = useState<Node | null>(null);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [copiedPeerKey, setCopiedPeerKey] = useState<string | null>(null);

  const handleRefresh = async () => {
    setRefreshing(true);
    try {
      await fetchNodes();
    } finally {
      setRefreshing(false);
    }
  };

  const handleDelete = async (id: number, name: string) => {
    if (confirm(`Are you sure you want to delete node "${name}"?`)) {
      try {
        await deleteNode(id);
      } catch (err) {
        alert(`Failed to delete node: ${err instanceof Error ? err.message : 'Unknown error'}`);
      }
    }
  };

  const copyPeerKey = async (peerKey: string, nodeId: number) => {
    try {
      await navigator.clipboard.writeText(peerKey);
      setCopiedPeerKey(peerKey);
      setTimeout(() => setCopiedPeerKey(null), 2000);
    } catch (err) {
      console.error('Failed to copy peer key:', err);
    }
  };

  const getHealthStatusColor = (status: string) => {
    switch (status) {
      case 'healthy': return 'success';
      case 'unhealthy': return 'danger';
      case 'checking': return 'warning';
      default: return 'secondary';
    }
  };

  return (
    <MainLayout
      title="Node Management"
      subtitle="Manage your Quilibrium node cluster"
      onRefresh={handleRefresh}
      refreshing={refreshing}
      actions={
        <Button onClick={() => setShowCreateModal(true)}>
          <Plus className="h-4 w-4 mr-2" />
          Add Node
        </Button>
      }
    >
      <div className="space-y-6">
        {/* Summary Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center space-x-2">
                <Server className="h-5 w-5 text-blue-600" />
                <div>
                  <div className="text-2xl font-bold">{nodes.length}</div>
                  <div className="text-sm text-gray-600">Total Nodes</div>
                </div>
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center space-x-2">
                <CheckCircle2 className="h-5 w-5 text-green-600" />
                <div>
                  <div className="text-2xl font-bold text-green-600">
                    {nodes.filter(n => n.health_status === 'healthy').length}
                  </div>
                  <div className="text-sm text-gray-600">Healthy</div>
                </div>
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center space-x-2">
                <div className="h-5 w-5 rounded-full bg-red-500" />
                <div>
                  <div className="text-2xl font-bold text-red-600">
                    {nodes.filter(n => n.health_status === 'unhealthy').length}
                  </div>
                  <div className="text-sm text-gray-600">Unhealthy</div>
                </div>
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center space-x-2">
                <Eye className="h-5 w-5 text-blue-600" />
                <div>
                  <div className="text-2xl font-bold">
                    {nodes.filter(n => n.enabled).length}
                  </div>
                  <div className="text-sm text-gray-600">Enabled</div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Nodes Table */}
        <Card>
          <CardHeader>
            <CardTitle>Node List</CardTitle>
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
            ) : nodes.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                <Server className="mx-auto h-12 w-12 text-gray-400 mb-4" />
                <h3 className="text-lg font-medium mb-2">No nodes found</h3>
                <p className="text-sm mb-4">Get started by creating your first node.</p>
                <Button onClick={() => setShowCreateModal(true)}>
                  <Plus className="h-4 w-4 mr-2" />
                  Create Node
                </Button>
              </div>
            ) : (
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Node
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Health Status
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Ports
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Peer Key
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Status
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Last Check
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Actions
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {nodes.map((node) => (
                      <tr key={node.id} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="flex items-center">
                            <div className="flex-shrink-0 h-10 w-10">
                              <div className="h-10 w-10 rounded-lg bg-blue-100 flex items-center justify-center">
                                <Server className="h-5 w-5 text-blue-600" />
                              </div>
                            </div>
                            <div className="ml-4">
                              <div className="text-sm font-medium text-gray-900">
                                {node.name}
                              </div>
                              <div className="text-sm text-gray-500">
                                ID: {node.id}
                              </div>
                            </div>
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="flex items-center space-x-2">
                            <div className={`w-2 h-2 rounded-full ${
                              node.health_status === 'healthy' ? 'bg-green-500' :
                              node.health_status === 'unhealthy' ? 'bg-red-500' :
                              node.health_status === 'checking' ? 'bg-yellow-500 animate-pulse' :
                              'bg-gray-400'
                            }`} />
                            <Badge variant={getHealthStatusColor(node.health_status)}>
                              {node.health_status}
                            </Badge>
                          </div>
                          {node.response_time > 0 && (
                            <div className="text-xs text-gray-500 mt-1">
                              {node.response_time.toFixed(0)}ms
                            </div>
                          )}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          <code className="bg-gray-100 px-2 py-1 rounded text-xs">
                            {getPortRange(node.base_port)}
                          </code>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="flex items-center space-x-2">
                            <code className="text-xs text-gray-600 max-w-32 truncate block">
                              {node.peer_key}
                            </code>
                            <button
                              onClick={() => copyPeerKey(node.peer_key, node.id)}
                              className="text-gray-400 hover:text-gray-600"
                            >
                              {copiedPeerKey === node.peer_key ? (
                                <CheckCircle2 className="h-4 w-4 text-green-500" />
                              ) : (
                                <Copy className="h-4 w-4" />
                              )}
                            </button>
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="flex items-center space-x-2">
                            {node.enabled ? (
                              <Eye className="h-4 w-4 text-green-500" />
                            ) : (
                              <EyeOff className="h-4 w-4 text-gray-400" />
                            )}
                            <span className="text-sm text-gray-900">
                              {node.enabled ? 'Enabled' : 'Disabled'}
                            </span>
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {node.last_health_check ? formatRelativeTime(node.last_health_check) : 'Never'}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium space-x-2">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => setSelectedNode(node)}
                          >
                            <Edit className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleDelete(node.id, node.name)}
                            className="text-red-600 hover:text-red-700"
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </MainLayout>
  );
}