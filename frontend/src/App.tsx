import { WebSocketProvider } from '@/context/websocket-context'
import { Dashboard } from '@/components/dashboard/dashboard'

function App() {
  return (
    <WebSocketProvider>
      <Dashboard />
    </WebSocketProvider>
  )
}

export default App
