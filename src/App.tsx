import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Login from './Login/Login';
import Dashboard from './Dashboard/Dashboard';
import NewBank from './NewBank/NewBank';
import ProtectedRoute from './ProtectedRoute';
import { useState } from 'react'
import './App.css'

function App() {
  const [username, setUsername] = useState<string | null>(null);
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={
          username ? <Navigate to="/" replace /> : <Login onLogin={setUsername} />
        } />
        <Route path="/new" element={
          username ? <Navigate to="/" replace /> : <NewBank />
        } />
        <Route path="/" element={<ProtectedRoute isAuthenticated={!!username} />}>
          <Route index element={<Dashboard username={username ?? ''} />} />
        </Route>
        <Route path="*" element={<Navigate to={username ? "/" : "/login"} replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
