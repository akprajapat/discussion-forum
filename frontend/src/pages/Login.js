import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

function Login({ onLogin }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const res = await fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
      });
      if (res.ok) {
        const data = await res.json();
        localStorage.setItem('token', data.token);
        if (data.user) {
          localStorage.setItem('user', JSON.stringify(data.user));
          if (onLogin) onLogin(data.user);
        } else {
          const userObj = { username };
          localStorage.setItem('user', JSON.stringify(userObj));
          if (onLogin) onLogin(userObj);
        }
        navigate('/');
      } else {
        const data = await res.json();
        alert('Login failed: ' + (data.error || 'Unknown error'));
      }
    } catch (err) {
      alert('Login failed: Network error');
    }
  };

  return (
    <div style={{ maxWidth: 400, margin: '40px auto', background: '#fff', padding: 32, borderRadius: 4, border: '1px solid #e4e6e8' }}>
      <h2 style={{ fontWeight: 400, fontSize: 24, marginBottom: 24 }}>Login</h2>
      <form onSubmit={handleSubmit}>
        <div style={{ marginBottom: 16 }}>
          <label>Username</label>
          <input
            type="text"
            value={username}
            onChange={e => setUsername(e.target.value)}
            required
            style={{ width: '100%', padding: 8, borderRadius: 4, border: '1px solid #ccc', marginTop: 4 }}
          />
        </div>
        <div style={{ marginBottom: 24 }}>
          <label>Password</label>
          <input
            type="password"
            value={password}
            onChange={e => setPassword(e.target.value)}
            required
            style={{ width: '100%', padding: 8, borderRadius: 4, border: '1px solid #ccc', marginTop: 4 }}
          />
        </div>
        <button type="submit" style={{ background: '#f48024', color: '#fff', padding: '0.5rem 1.5rem', border: 'none', borderRadius: 4, fontWeight: 600, fontSize: 16 }}>
          Login
        </button>
      </form>
    </div>
  );
}

export default Login;

