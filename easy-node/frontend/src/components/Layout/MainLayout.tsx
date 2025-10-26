'use client';

import { Sidebar } from './Sidebar';
import { Header } from './Header';

interface MainLayoutProps {
  children: React.ReactNode;
  title: string;
  subtitle?: string;
  actions?: React.ReactNode;
  onRefresh?: () => void;
  refreshing?: boolean;
}

export function MainLayout({ 
  children, 
  title, 
  subtitle, 
  actions, 
  onRefresh, 
  refreshing 
}: MainLayoutProps) {
  return (
    <div className="flex h-screen bg-gray-50">
      <Sidebar />
      <div className="flex flex-1 flex-col overflow-hidden">
        <Header 
          title={title}
          subtitle={subtitle}
          actions={actions}
          onRefresh={onRefresh}
          refreshing={refreshing}
        />
        <main className="flex-1 overflow-auto p-6">
          {children}
        </main>
      </div>
    </div>
  );
}