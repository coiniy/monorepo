'use client';

import { useEffect, useState } from 'react';
import { MainLayout } from '@/components/Layout/MainLayout';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/Card';
import { Badge } from '@/components/ui/Badge';
import { useHealth, useMonitor } from '@/hooks/useNodes';
import { formatRelativeTime, getStatusColor } from '@/lib/utils';
import { 
  Server, 
  Activity, 
  Clock, 
  TrendingUp,
  AlertTriangle,
  CheckCircle2,
  XCircle,
  Timer
} from 'lucide-react';

export default function Dashboard() {
  const { healthStats, healthData, loading: healthLoading, fetchHealth } = useHealth();
  const { metrics, loading: metricsLoading, fetchMetrics } = useMonitor();
  const [refreshing, setRefreshing] = useState(false);

  const handleRefresh = async () => {
    setRefreshing(true);
    try {
      await Promise.all([fetchHealth(), fetchMetrics()]);
    } finally {
      setRefreshing(false);
    }
  };

  // Auto-refresh every 30 seconds
  useEffect(() => {
    const interval = setInterval(() => {
      fetchHealth();
      fetchMetrics();
    }, 30000);

    return () => clearInterval(interval);
  }, [fetchHealth, fetchMetrics]);

  const loading = healthLoading || metricsLoading;

  return (
    <MainLayout
      title="Dashboard"
      subtitle="Overview of your Quilibrium node cluster"
      onRefresh={handleRefresh}
      refreshing={refreshing}
    >
      <div className="space-y-6">
        {/* Overview Stats */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Nodes</CardTitle>
              <Server className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {loading ? '...' : healthStats?.total_nodes || 0}
              </div>
              <p className="text-xs text-muted-foreground">
                Active cluster nodes
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Healthy Nodes</CardTitle>
              <CheckCircle2 className="h-4 w-4 text-green-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-green-600">
                {loading ? '...' : healthStats?.healthy_nodes || 0}
              </div>
              <p className="text-xs text-muted-foreground">
                {loading ? '...' : `${healthStats?.health_percentage.toFixed(1)}%`} uptime
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Unhealthy Nodes</CardTitle>
              <XCircle className="h-4 w-4 text-red-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-red-600">
                {loading ? '...' : (healthStats?.unhealthy_nodes || 0) + (healthStats?.unknown_nodes || 0)}
              </div>
              <p className="text-xs text-muted-foreground">
                Requires attention
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Avg Response Time</CardTitle>
              <Timer className="h-4 w-4 text-blue-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {loading ? '...' : `${healthStats?.avg_response_time.toFixed(0)}ms`}
              </div>
              <p className="text-xs text-muted-foreground">
                Network latency
              </p>
            </CardContent>
          </Card>
        </div>

        {/* Performance Overview */}
        {metrics && (
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <Card className="lg:col-span-2">
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <TrendingUp className="h-5 w-5" />
                  <span>Performance Metrics</span>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                  <div className="text-center">
                    <div className="text-2xl font-bold text-green-600">
                      {metrics.performance.online_nodes}
                    </div>
                    <div className="text-sm text-gray-600">Online</div>
                  </div>
                  <div className="text-center">
                    <div className="text-2xl font-bold text-blue-600">
                      {metrics.performance.avg_response_time.toFixed(0)}ms
                    </div>
                    <div className="text-sm text-gray-600">Avg Response</div>
                  </div>
                  <div className="text-center">
                    <div className="text-2xl font-bold text-orange-600">
                      {metrics.performance.max_response_time.toFixed(0)}ms
                    </div>
                    <div className="text-sm text-gray-600">Max Response</div>
                  </div>
                  <div className="text-center">
                    <div className="text-2xl font-bold text-purple-600">
                      {metrics.performance.uptime_percentage.toFixed(1)}%
                    </div>
                    <div className="text-sm text-gray-600">Uptime</div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <Activity className="h-5 w-5" />
                  <span>Monitoring Status</span>
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Scheduler</span>
                  <Badge variant={metrics.scheduler.running ? 'success' : 'danger'}>
                    {metrics.scheduler.running ? 'Running' : 'Stopped'}
                  </Badge>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Check Interval</span>
                  <span className="text-sm text-gray-600">{metrics.scheduler.interval}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Next Check</span>
                  <span className="text-sm text-gray-600">
                    {formatRelativeTime(metrics.scheduler.next_check)}
                  </span>
                </div>
              </CardContent>
            </Card>
          </div>
        )}

        {/* Recent Node Status */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <Activity className="h-5 w-5" />
              <span>Recent Health Checks</span>
            </CardTitle>
          </CardHeader>
          <CardContent>
            {loading ? (
              <div className="flex items-center justify-center py-8">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
              </div>
            ) : healthData && healthData.length > 0 ? (
              <div className="space-y-3">
                {healthData.slice(0, 5).map((node) => (
                  <div key={node.node_id} className="flex items-center justify-between py-2 border-b border-gray-100 last:border-0">
                    <div className="flex items-center space-x-3">
                      <div className={`w-3 h-3 rounded-full ${
                        node.health_status === 'healthy' ? 'bg-green-500' :
                        node.health_status === 'unhealthy' ? 'bg-red-500' :
                        node.health_status === 'checking' ? 'bg-yellow-500 animate-pulse' :
                        'bg-gray-400'
                      }`} />
                      <div>
                        <div className="font-medium">{node.node_name}</div>
                        <div className="text-sm text-gray-600">{node.health_message}</div>
                      </div>
                    </div>
                    <div className="text-right">
                      <Badge 
                        variant={
                          node.health_status === 'healthy' ? 'success' :
                          node.health_status === 'unhealthy' ? 'danger' :
                          node.health_status === 'checking' ? 'warning' :
                          'secondary'
                        }
                      >
                        {node.health_status}
                      </Badge>
                      <div className="text-xs text-gray-500 mt-1">
                        {node.response_time > 0 ? `${node.response_time.toFixed(0)}ms` : '-'}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-8 text-gray-500">
                No health data available
              </div>
            )}
          </CardContent>
        </Card>

        {/* System Status */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Clock className="h-5 w-5" />
                <span>Last Updated</span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-lg font-medium">
                {healthStats ? formatRelativeTime(healthStats.last_check) : 'Never'}
              </div>
              <p className="text-sm text-gray-600 mt-1">
                Health check timestamp
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <AlertTriangle className="h-5 w-5" />
                <span>Alerts</span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                {healthStats && healthStats.unhealthy_nodes > 0 && (
                  <div className="flex items-center space-x-2 text-red-600">
                    <XCircle className="h-4 w-4" />
                    <span className="text-sm">{healthStats.unhealthy_nodes} unhealthy nodes</span>
                  </div>
                )}
                {healthStats && healthStats.unknown_nodes > 0 && (
                  <div className="flex items-center space-x-2 text-yellow-600">
                    <AlertTriangle className="h-4 w-4" />
                    <span className="text-sm">{healthStats.unknown_nodes} unknown status</span>
                  </div>
                )}
                {healthStats && healthStats.unhealthy_nodes === 0 && healthStats.unknown_nodes === 0 && (
                  <div className="flex items-center space-x-2 text-green-600">
                    <CheckCircle2 className="h-4 w-4" />
                    <span className="text-sm">All systems operational</span>
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </MainLayout>
  );
}