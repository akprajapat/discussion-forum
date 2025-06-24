import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

function Home() {
  const [questions, setQuestions] = useState([]);
  const [search, setSearch] = useState('');

  useEffect(() => {
    const fetchQuestions = async () => {
      try {
        const res = await fetch(`/api/questions?q=${encodeURIComponent(search)}`);
        const data = await res.json();
        // Ensure data is always an array
        setQuestions(Array.isArray(data) ? data : []);
      } catch (err) {
        console.error('Home: fetch error', err);
      }
    };
    fetchQuestions();
  }, [search]);

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <h1 style={{ fontWeight: 400, fontSize: 32 }}>Top Questions</h1>
        <Link to="/questions/ask" style={{ background: '#f48024', color: '#fff', padding: '0.5rem 1rem', borderRadius: 4, textDecoration: 'none', fontWeight: 600 }}>Ask Question</Link>
      </div>
      <input
        type="text"
        placeholder="Search questions..."
        value={search}
        onChange={e => setSearch(e.target.value)}
        style={{ marginBottom: 24, width: '100%', padding: 8, borderRadius: 4, border: '1px solid #ccc' }}
      />
      <div>
        {questions.length === 0 && <div style={{ color: '#888', marginTop: 32 }}>No questions found.</div>}
        {questions.map(q => (
          <div key={q.id || q._id} style={{ display: 'flex', background: '#fff', border: '1px solid #e4e6e8', borderRadius: 4, marginBottom: 16, padding: 16 }}>
            <div style={{ width: 80, textAlign: 'center', color: '#6a737c', fontSize: 18 }}>
              <div><b>{q.votes || 0}</b><br />votes</div>
            </div>
            <div style={{ flex: 1 }}>
              <Link to={`/questions/${q.id || q._id}`} style={{ fontSize: 20, color: '#0074cc', textDecoration: 'none', fontWeight: 500 }}>{q.title}</Link>
              <div style={{ margin: '8px 0', color: '#232629' }}>{q.body?.slice(0, 120)}{q.body && q.body.length > 120 ? '...' : ''}</div>
              <div>
                {(q.tags || []).map(tag => (
                  <span key={tag} style={{ background: '#e1ecf4', color: '#39739d', borderRadius: 3, padding: '2px 8px', marginRight: 8, fontSize: 13 }}>{tag}</span>
                ))}
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

export default Home;
