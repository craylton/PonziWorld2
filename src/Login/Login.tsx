import { useState, useEffect } from 'react';
import './Login.css';
import { useNavigate } from 'react-router-dom';
import { removeAuthToken, setAuthToken } from '../auth';
import PageHeader from '../components/PageHeader';

interface LoginProps {
  onLogin: () => Promise<void>;
}

export default function Login({ onLogin }: LoginProps) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    // Clear any existing auth token when visiting login page
    removeAuthToken();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    if (username.trim() && password.trim()) {
      setLoading(true);
      try {
        const res = await fetch('/api/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            username: username.trim(),
            password: password.trim()
          }),
        });
        if (!res.ok) {
          let data;
          try {
            data = await res.json();
          } catch {
            data = {};
          }
          setError(`Login failed: ${data.error || 'Unknown error'}`);
        } else {
          const data = await res.json();
          // Store the JWT
          setAuthToken(data.token);
          // Navigate to dashboard
          await onLogin();
          navigate('/');
        }
      } catch (error) {
        console.error('Login error:', error);
        setError('Network error');
      } finally {
        setLoading(false);
      }
    }
  };

  return (
    <div className="login-container">
      <PageHeader title="Login" />
      <form className="login-form" onSubmit={handleSubmit}>
        <input
          type="text"
          placeholder="Enter username"
          value={username}
          onChange={e => setUsername(e.target.value)}
          required
        />
        <input
          type="password"
          placeholder="Enter password"
          value={password}
          onChange={e => setPassword(e.target.value)}
          required
        />
        {error && <div className="error-msg">{error}</div>}
        <button
          type="submit"
          disabled={loading}
          className="login-btn"
        >
          {loading ? 'Logging in...' : 'Login'}
        </button>
        <button
          type="button"
          className="new-bank-btn"
          onClick={() => navigate('/new')}
        >
          New player? Click here to start a new bank
        </button>
      </form>
    </div>
  );
}
