import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

function AskQuestion() {
  const [title, setTitle] = useState('');
  const [body, setBody] = useState('');
  const [tags, setTags] = useState('');
  const navigate = useNavigate();

  const submit = async (e) => {
    e.preventDefault();
    const token = localStorage.getItem('token');
    try {
      const res = await fetch('/api/questions', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': token ? `Bearer ${token}` : '',
        },
        body: JSON.stringify({ title, body, tags: tags.split(' ').filter(Boolean) }),
      });
      if (res.ok) {
        navigate('/');
      } else {
        const data = await res.json();
        alert('Failed to ask question: ' + (data.error || 'Unknown error'));
      }
    } catch (err) {
      alert('Failed to ask question: Network error');
    }
  };

  return (
    <div style={{ maxWidth: 700, margin: '0 auto' }}>
      <h2 style={{ fontWeight: 400, fontSize: 28, marginBottom: 24 }}>Ask a public question</h2>
      <form onSubmit={submit} style={{ background: '#fff', padding: 24, borderRadius: 4, border: '1px solid #e4e6e8' }}>
        <div style={{ marginBottom: 16 }}>
          <label style={{ fontWeight: 600 }}>Title</label>
          <input
            value={title}
            onChange={e => setTitle(e.target.value)}
            required
            placeholder="Be specific and imagine you’re asking another person"
            style={{ width: '100%', padding: 8, borderRadius: 4, border: '1px solid #ccc', marginTop: 4 }}
          />
        </div>
        <div style={{ marginBottom: 16 }}>
          <label style={{ fontWeight: 600 }}>Body</label>
          <textarea
            value={body}
            onChange={e => setBody(e.target.value)}
            required
            placeholder="Include all the information someone would need to answer your question"
            rows={8}
            style={{ width: '100%', padding: 8, borderRadius: 4, border: '1px solid #ccc', marginTop: 4 }}
          />
        </div>
        <div style={{ marginBottom: 16 }}>
          <label style={{ fontWeight: 600 }}>Tags</label>
          <input
            value={tags}
            onChange={e => setTags(e.target.value)}
            placeholder="e.g. javascript react css"
            style={{ width: '100%', padding: 8, borderRadius: 4, border: '1px solid #ccc', marginTop: 4 }}
          />
        </div>
        <button type="submit" style={{ background: '#f48024', color: '#fff', padding: '0.5rem 1.5rem', border: 'none', borderRadius: 4, fontWeight: 600, fontSize: 16 }}>
          Post your question
        </button>
      </form>
    </div>
  );
}

export default AskQuestion;
