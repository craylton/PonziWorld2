import { useState } from 'react';
import './Login.css';
import { useNavigate } from 'react-router-dom';
import type { User } from '../User';

interface LoginProps {
  onLogin: (user: User) => void;
}

export default function Login({ onLogin }: LoginProps) {
  const [username, setUsername] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    if (username.trim()) {
      setLoading(true);
      try {
        const res = await fetch('/api/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ username: username.trim() }),
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
          // Now fetch the user object
          const userRes = await fetch(`/api/user?username=${encodeURIComponent(username.trim())}`);
          if (!userRes.ok) {
            let data;
            try {
              data = await userRes.json();
            } catch {
              data = {};
            }
            setError(`Failed to fetch user: ${data.error || 'Unknown error'}`);
          } else {
            const userData = await userRes.json();
            onLogin(userData);
          }
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
      <form className="login-form" onSubmit={handleSubmit}>
        <h2>Login</h2>
        <input
          type="text"
          placeholder="Enter username"
          value={username}
          onChange={e => setUsername(e.target.value)}
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
          New user? Click here to start a new bank
        </button>
      </form>
    </div>
  );
}
