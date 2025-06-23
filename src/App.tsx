import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Login from './Login/Login';
import Dashboard from './Dashboard/Dashboard';
import NewBank from './NewBank/NewBank';
import ProtectedRoute from './ProtectedRoute';
import { useState, useEffect } from 'react'
import './App.css'
import type { User } from './User';
import { isAuthenticated, makeAuthenticatedRequest, removeAuthToken } from './auth';

function App() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Check if user is already authenticated on app start
    const initializeAuth = async () => {
      if (isAuthenticated()) {
        try {
          const response = await makeAuthenticatedRequest('/api/user');
          if (response.ok) {
            const userData = await response.json();
            setUser(userData);
          } else {
            // Token is invalid, remove it
            removeAuthToken();
          }
        } catch (error) {
          console.error('Failed to fetch user data:', error);
          removeAuthToken();
        }
      }
      setLoading(false);
    };

    initializeAuth();
  }, []);

  const handleLogout = () => {
    removeAuthToken();
    setUser(null);
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={
          user ? <Navigate to="/" replace /> : <Login onLogin={setUser} />
        } />
        <Route path="/new" element={
          user ? <Navigate to="/" replace /> : <NewBank />
        } />
        <Route path="/" element={<ProtectedRoute isAuthenticated={!!user} />}>
          <Route index element={<Dashboard user={user!} onLogout={handleLogout} />} />
        </Route>
        <Route path="*" element={<Navigate to={user ? "/" : "/login"} replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
