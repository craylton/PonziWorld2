import Login from './Login/Login';
import Dashboard from './Dashboard/Dashboard';
import { useState } from 'react'
import './App.css'

function App() {
  const [username, setUsername] = useState<string | null>(null);

  if (!username) {
    return <Login onLogin={setUsername} />;
  }

  return <Dashboard username={username} />;
}

export default App
