import Login from './Login';
import Dashboard from './Dashboard';
import { useState } from 'react'
import './App.css'

// Types for state

function App() {
  // Only keep username state and login logic
  const [username, setUsername] = useState<string | null>(null);

  if (!username) {
    return <Login onLogin={setUsername} />;
  }

  // Show dashboard after login
  return <Dashboard username={username} />;
}

export default App
