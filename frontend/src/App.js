import React, { useEffect, useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import AskQuestion from './pages/AskQuestion';
import QuestionDetail from './pages/QuestionDetail';

function UserMenu({ user, onLogout }) {
  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
      <span title={user?.username || 'User'}>
        <svg width="32" height="32" viewBox="0 0 32 32" style={{ borderRadius: '50%', background: '#e1ecf4', marginRight: 4 }}>
          <circle cx="16" cy="16" r="16" fill="#e1ecf4" />
          <text x="16" y="21" textAnchor="middle" fontSize="16" fill="#39739d" fontFamily="Arial, sans-serif">
            {user?.username ? user.username[0].toUpperCase() : 'U'}
          </text>
        </svg>
      </span>
      <button onClick={onLogout} style={{ background: '#f48024', color: '#fff', border: 'none', borderRadius: 4, padding: '0.3rem 1rem', fontWeight: 600, cursor: 'pointer' }}>
        Logout
      </button>
    </div>
  );
}

function App() {
  const [user, setUser] = useState(null);

  // Check for token and user info on mount
  useEffect(() => {
    const token = localStorage.getItem('token');
    const userInfo = localStorage.getItem('user');
    if (token && userInfo) {
      setUser(JSON.parse(userInfo));
    }
  }, []);

  // Listen for login/logout events from child components
  const handleLogin = (userObj) => {
    setUser(userObj);
    localStorage.setItem('user', JSON.stringify(userObj));
  };
  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setUser(null);
  };

  return (
    <Router>
      <div style={{ display: 'flex', minHeight: '100vh', background: '#f6f6f6' }}>
        <aside style={{ width: 220, background: '#232629', color: '#fff', padding: '2rem 1rem 0 1rem' }}>
          <h2 style={{ color: '#f48024', marginBottom: '2rem' }}>Discussion Forum</h2>
          <nav>
            <ul style={{ listStyle: 'none', padding: 0 }}>
              <li><Link to="/" style={{ color: '#fff', textDecoration: 'none', display: 'block', padding: '0.5rem 0' }}>Home</Link></li>
              <li><Link to="/questions/ask" style={{ color: '#fff', textDecoration: 'none', display: 'block', padding: '0.5rem 0' }}>Ask Question</Link></li>
              {!user && (
                <>
                  <li><Link to="/login" style={{ color: '#fff', textDecoration: 'none', display: 'block', padding: '0.5rem 0' }}>Login</Link></li>
                  <li><Link to="/register" style={{ color: '#fff', textDecoration: 'none', display: 'block', padding: '0.5rem 0' }}>Register</Link></li>
                </>
              )}
            </ul>
          </nav>
        </aside>
        <div style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
          <header style={{ background: '#fff', borderBottom: '1px solid #e4e6e8', padding: '1rem 2rem', fontWeight: 'bold', fontSize: '1.2rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <span></span>
            {user && <UserMenu user={user} onLogout={handleLogout} />}
          </header>
          <main style={{ flex: 1, padding: '2rem', maxWidth: 900, margin: '0 auto', width: '100%' }}>
            <Routes>
              <Route path="/" element={<Home />} />
              <Route path="/login" element={<Login onLogin={handleLogin} />} />
              <Route path="/register" element={<Register />} />
              <Route path="/questions/ask" element={<AskQuestion />} />
              <Route path="/questions/:id" element={<QuestionDetail />} />
            </Routes>
          </main>
        </div>
      </div>
    </Router>
  );
}

export default App;
