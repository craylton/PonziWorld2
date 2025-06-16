import React, { useState } from "react";
import "./NewBank.css";

const NewBank: React.FC = () => {
  const [username, setUsername] = useState("");
  const [bankName, setBankName] = useState("");

  return (
    <div className="new-bank-container">
      <h2>Create a New Bank</h2>
      <form className="new-bank-form" onSubmit={e => e.preventDefault()}>
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
        <button type="submit" className="confirm-btn" disabled>
          Confirm
        </button>
      </form>
    </div>
  );
};

export default NewBank;
