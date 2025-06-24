import React, { useState } from 'react';

export default function AskQuestion() {
  const [title, setTitle] = useState('');
  const [body, setBody] = useState('');

  const submit = async (e) => {
    e.preventDefault();
    await fetch('/api/questions', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title, body }),
    });
    // redirect to home or question page
  };

  return (
    <form onSubmit={submit}>
      <input value={title} onChange={e => setTitle(e.target.value)} placeholder="Title" />
      <textarea value={body} onChange={e => setBody(e.target.value)} placeholder="Body" />
      <button type="submit">Ask</button>
    </form>
  );
}
