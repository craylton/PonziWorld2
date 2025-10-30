import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Login from './Login/Login';
import Dashboard from './Dashboard/Dashboard';
import Home from './Home/Home';
import NewBank from './NewBank/NewBank';
import ProtectedRoute from './ProtectedRoute';
import { useState, useEffect } from 'react'
import './App.css'
import { isAuthenticated, removeAuthToken } from './auth';

function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [authLoading, setAuthLoading] = useState(true);

  useEffect(() => {
    // Check if auth token exists on app start
    const initializeAuth = () => {
      if (isAuthenticated()) {
        setIsLoggedIn(true);
      }
      setAuthLoading(false);
    };

    initializeAuth();
  }, []);

  // Handler after successful login
  const handleLogin = async () => {
    setIsLoggedIn(true);
  };

  const handleLogout = () => {
    removeAuthToken();
    setIsLoggedIn(false);
  };

  if (authLoading) {
    return <div>Loading...</div>;
  }

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={
          isLoggedIn ? <Navigate to="/" replace /> : <Login onLogin={handleLogin} />
        } />
        <Route path="/new" element={
          isLoggedIn ? <Navigate to="/" replace /> : <NewBank />
        } />
        <Route path="/" element={<ProtectedRoute isAuthenticated={isLoggedIn} />}>
          <Route index element={<Home onLogout={handleLogout} />} />
          <Route path="bank/:bankId" element={<Dashboard onLogout={handleLogout} />} />
        </Route>
        <Route path="*" element={<Navigate to={isLoggedIn ? "/" : "/login"} replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
