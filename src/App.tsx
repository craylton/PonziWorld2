import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Login from './Login/Login';
import Dashboard from './Dashboard/Dashboard';
import NewBank from './NewBank/NewBank';
import ProtectedRoute from './ProtectedRoute';
import { useState } from 'react'
import './App.css'

interface User {
  id: string;
  username: string;
  bankName: string;
  claimedCapital: number;
  actualCapital: number;
}

function App() {
  const [user, setUser] = useState<User | null>(null);
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
          <Route index element={<Dashboard user={user!} />} />
        </Route>
        <Route path="*" element={<Navigate to={user ? "/" : "/login"} replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
