import { type ClassValue, clsx } from 'clsx';

export function cn(...inputs: ClassValue[]) {
  return clsx(inputs);
}

export function formatBytes(bytes: number, decimals: number = 2): string {
  if (bytes === 0) return '0 Bytes';

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

export function formatDuration(ms: number): string {
  if (ms < 1000) return `${Math.round(ms)}ms`;
  if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`;
  if (ms < 3600000) return `${Math.round(ms / 60000)}m`;
  return `${Math.round(ms / 3600000)}h`;
}

export function formatRelativeTime(date: string | Date): string {
  const now = new Date();
  const target = new Date(date);
  const diffMs = now.getTime() - target.getTime();
  
  if (diffMs < 60000) return 'just now';
  if (diffMs < 3600000) return `${Math.round(diffMs / 60000)}m ago`;
  if (diffMs < 86400000) return `${Math.round(diffMs / 3600000)}h ago`;
  return `${Math.round(diffMs / 86400000)}d ago`;
}

export function getStatusColor(status: string): string {
  switch (status.toLowerCase()) {
    case 'healthy':
    case 'running':
      return 'text-green-600 bg-green-50';
    case 'unhealthy':
    case 'stopped':
    case 'exited':
      return 'text-red-600 bg-red-50';
    case 'checking':
    case 'starting':
      return 'text-yellow-600 bg-yellow-50';
    case 'unknown':
    default:
      return 'text-gray-600 bg-gray-50';
  }
}

export function getStatusIcon(status: string): string {
  switch (status.toLowerCase()) {
    case 'healthy':
    case 'running':
      return '🟢';
    case 'unhealthy':
    case 'stopped':
    case 'exited':
      return '🔴';
    case 'checking':
    case 'starting':
      return '🟡';
    case 'unknown':
    default:
      return '⚪';
  }
}

export function debounce<T extends (...args: any[]) => any>(
  func: T,
  wait: number
): (...args: Parameters<T>) => void {
  let timeout: NodeJS.Timeout;
  return (...args: Parameters<T>) => {
    clearTimeout(timeout);
    timeout = setTimeout(() => func(...args), wait);
  };
}

export function throttle<T extends (...args: any[]) => any>(
  func: T,
  limit: number
): (...args: Parameters<T>) => void {
  let inThrottle: boolean;
  return (...args: Parameters<T>) => {
    if (!inThrottle) {
      func(...args);
      inThrottle = true;
      setTimeout(() => (inThrottle = false), limit);
    }
  };
}

export function copyToClipboard(text: string): Promise<void> {
  if (navigator.clipboard && window.isSecureContext) {
    return navigator.clipboard.writeText(text);
  } else {
    // Fallback for older browsers
    const textArea = document.createElement('textarea');
    textArea.value = text;
    textArea.style.position = 'absolute';
    textArea.style.opacity = '0';
    textArea.style.left = '-999999px';
    textArea.style.top = '-999999px';
    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();
    return new Promise((resolve, reject) => {
      document.execCommand('copy') ? resolve() : reject();
      textArea.remove();
    });
  }
}

export function generateNodeName(): string {
  const adjectives = ['swift', 'bright', 'clever', 'brave', 'quick', 'smart', 'bold', 'fast'];
  const nouns = ['falcon', 'eagle', 'wolf', 'lion', 'tiger', 'shark', 'hawk', 'bear'];
  
  const adjective = adjectives[Math.floor(Math.random() * adjectives.length)];
  const noun = nouns[Math.floor(Math.random() * nouns.length)];
  const number = Math.floor(Math.random() * 1000);
  
  return `${adjective}-${noun}-${number}`;
}

export function validatePort(port: number): boolean {
  return port >= 1024 && port <= 65535;
}

export function getPortRange(basePort: number): string {
  return `${basePort}-${basePort + 4}`;
}