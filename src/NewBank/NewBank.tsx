import React, { useState } from "react";
import "./NewBank.css";

const NewBank: React.FC = () => {
  const [username, setUsername] = useState("");
  const [bankName, setBankName] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setSuccess(false);
    try {
      const res = await fetch("/api/user", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, bankName }),
      });
      if (!res.ok) {
        const data = await res.json();
        setError(data.error || "Failed to create user");
      } else {
        setSuccess(true);
        setUsername("");
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
      <h2>Create a New Bank</h2>
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
          Bank Name
          <input
            type="text"
            value={bankName}
            onChange={e => setBankName(e.target.value)}
            placeholder="Enter bank name"
            required
          />
        </label>
        <button type="submit" className="confirm-btn" disabled={!username || !bankName || loading}>
          {loading ? "Creating..." : "Confirm"}
        </button>
        {error && <div className="error-msg">{error}</div>}
        {success && <div className="success-msg">Bank created! You can now log in.</div>}
      </form>
    </div>
  );
};

export default NewBank;
