import { useState } from 'react';
import './Login.css';
import { useNavigation } from '../navigation';

interface LoginProps {
  onLogin: (username: string) => void;
}

export default function Login({ onLogin }: LoginProps) {
  const [username, setUsername] = useState('');
  const { setPage } = useNavigation();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (username.trim()) {
      onLogin(username.trim());
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
        <button type="submit">Login</button>
        <button
          type="button"
          className="new-bank-btn"
          onClick={() => setPage('newbank')}
          style={{ marginTop: '1rem', background: '#2d7ef7', color: '#fff', border: 'none', borderRadius: 6, padding: '0.75rem', fontSize: '1.1rem' }}
        >
          New Bank
        </button>
      </form>
    </div>
  );
}
