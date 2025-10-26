# Easy Node Frontend

Modern web interface for managing Quilibrium nodes, built with Next.js 14, TypeScript, and Tailwind CSS.

## Features

- **Dashboard** - Real-time overview of node cluster health and performance
- **Node Management** - Complete CRUD operations for nodes
- **Status Control** - Start, stop, restart nodes with real-time status updates
- **Health Monitoring** - Comprehensive health checks and performance metrics
- **Modern UI** - Clean, responsive design with intuitive navigation
- **Real-time Updates** - Automatic refresh and live status monitoring

## Tech Stack

- **Next.js 14** - React framework with App Router
- **TypeScript** - Type-safe development
- **Tailwind CSS** - Utility-first CSS framework
- **Lucide React** - Beautiful icon library
- **Custom Hooks** - Reusable state management
- **API Integration** - Full backend API integration

## Getting Started

### Prerequisites

- Node.js 18.17 or later
- Backend server running on `localhost:8080`

### Installation

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Start production server
npm start
```

The application will be available at [http://localhost:3000](http://localhost:3000).

## Project Structure

```
frontend/
├── src/
│   ├── app/                    # Next.js App Router pages
│   │   ├── page.tsx           # Dashboard
│   │   ├── nodes/             # Node management
│   │   ├── status/            # Status control
│   │   └── layout.tsx         # Root layout
│   ├── components/            # Reusable components
│   │   ├── ui/                # UI components
│   │   └── Layout/            # Layout components
│   ├── hooks/                 # Custom React hooks
│   ├── lib/                   # Utility functions and API
│   └── types/                 # TypeScript definitions
├── public/                    # Static assets
├── tailwind.config.ts         # Tailwind configuration
├── next.config.js            # Next.js configuration
└── package.json              # Dependencies and scripts
```

## API Integration

The frontend integrates with the backend API through:

- **Node Management** - CRUD operations for nodes
- **Status Control** - Container lifecycle management
- **Health Monitoring** - Real-time health checks
- **Performance Metrics** - System performance data

API calls are handled through custom hooks that provide loading states, error handling, and automatic refresh capabilities.

## Components

### Layout Components
- **Sidebar** - Navigation with active state indication
- **Header** - Page title, actions, and refresh functionality
- **MainLayout** - Wrapper combining sidebar and header

### UI Components
- **Button** - Multiple variants and loading states
- **Card** - Content containers with header/footer
- **Badge** - Status indicators with color variants

### Custom Hooks
- **useNodes** - Node management operations
- **useNodeStatus** - Container status control
- **useHealth** - Health monitoring data
- **useMonitor** - Performance metrics

## Styling

The application uses Tailwind CSS with:

- Custom color palette for status indicators
- Responsive grid layouts
- Smooth animations and transitions
- Dark mode support (configurable)
- Custom scrollbar styling

## Development

### Code Style
- TypeScript strict mode enabled
- ESLint with Next.js configuration
- Consistent naming conventions
- Component composition patterns

### Performance
- Client-side state management
- Optimized re-rendering
- Lazy loading where appropriate
- Efficient API calls with caching

## Deployment

The frontend can be deployed to any platform supporting Next.js:

```bash
# Build the application
npm run build

# Start production server
npm start
```

Make sure to configure the backend API URL in `next.config.js` for production deployment.

## Environment Variables

Create a `.env.local` file for environment-specific configuration:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Browser Support

- Chrome/Edge 88+
- Firefox 89+
- Safari 14+
- Mobile browsers (iOS Safari, Chrome Mobile)

## Contributing

1. Follow the existing code style and patterns
2. Add TypeScript types for new data structures
3. Include error handling in API calls
4. Test responsive design on mobile devices
5. Update documentation for new features