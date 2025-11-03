import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import "./NewBank.css";
import PageHeader from "../components/PageHeader";

const NewBank: React.FC = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [bankName, setBankName] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setSuccess(false);
    try {
      const res = await fetch("/api/newPlayer", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password, bankName }),
      });
      if (!res.ok) {
        const data = await res.json();
        setError(data.error || "Failed to create player");
      } else {
        setSuccess(true);
        setUsername("");
        setPassword("");
        setBankName("");
      }
    } catch {
      setError("Network error");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="new-bank-container">
      <PageHeader title="Create a New Bank" />
      <form className="new-bank-form" onSubmit={handleSubmit}>
        <label>
          Username
          <input
            type="text"
            value={username}
            onChange={e => setUsername(e.target.value)}
            placeholder="Enter your username"
            required
          />
        </label>
        <label>
          Password
          <input
            type="password"
            value={password}
            onChange={e => setPassword(e.target.value)}
            placeholder="Enter your password"
            required
          />
        </label>
        <label>
          Bank Name
          <input
            type="text"
            value={bankName}
            onChange={e => setBankName(e.target.value)}
            placeholder="Enter bank name"
            required
          />
        </label>
        <button
          type="submit"
          className="confirm-btn"
          disabled={!username || !password || !bankName || loading}
        >
          {loading ? "Creating..." : "Confirm"}
        </button>
        {error && <div className="error-msg">{error}</div>}
        {success && <div className="success-msg">Bank created! You can now log in.</div>}
        <button
          type="button"
          className="back-login-btn"
          onClick={() => navigate('/login')}
        >
          Back to Login
        </button>
      </form>
    </div>
  );
};

export default NewBank;
