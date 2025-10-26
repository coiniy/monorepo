'use client';

import { Bell, Refresh } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { Badge } from '@/components/ui/Badge';

interface HeaderProps {
  title: string;
  subtitle?: string;
  actions?: React.ReactNode;
  onRefresh?: () => void;
  refreshing?: boolean;
}

export function Header({ title, subtitle, actions, onRefresh, refreshing }: HeaderProps) {
  return (
    <header className="bg-white border-b border-gray-200 px-6 py-4">
      <div className="flex items-center justify-between">
        <div className="flex-1">
          <h1 className="text-2xl font-bold text-gray-900">{title}</h1>
          {subtitle && (
            <p className="mt-1 text-sm text-gray-600">{subtitle}</p>
          )}
        </div>
        
        <div className="flex items-center space-x-4">
          {onRefresh && (
            <Button
              variant="outline"
              size="sm"
              onClick={onRefresh}
              loading={refreshing}
              className="flex items-center space-x-2"
            >
              <Refresh className="h-4 w-4" />
              <span>Refresh</span>
            </Button>
          )}
          
          {actions}
          
          {/* Notification Bell */}
          <div className="relative">
            <Button variant="ghost" size="sm" className="p-2">
              <Bell className="h-5 w-5" />
            </Button>
            <Badge 
              variant="danger" 
              className="absolute -top-1 -right-1 h-5 w-5 rounded-full p-0 flex items-center justify-center text-xs"
            >
              3
            </Badge>
          </div>
        </div>
      </div>
    </header>
  );
}