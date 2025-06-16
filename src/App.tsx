import Login from './Login/Login';
import Dashboard from './Dashboard/Dashboard';
import { useState } from 'react'
import './App.css'
import NewBank from './NewBank/NewBank';
import { NavigationProvider, useNavigation } from './navigation';

function MainApp() {
  const [username, setUsername] = useState<string | null>(null);
  const { page } = useNavigation();

  if (page === 'newbank') {
    return <NewBank />;
  }

  if (!username) {
    return <Login onLogin={setUsername} />;
  }

  return <Dashboard username={username} />;
}

function App() {
  return (
    <NavigationProvider>
      <MainApp />
    </NavigationProvider>
  );
}

export default App;
