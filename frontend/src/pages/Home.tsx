import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

export default function Home() {
  const [questions, setQuestions] = useState([]);
  const [search, setSearch] = useState('');

  useEffect(() => {
    fetch(`/api/questions?q=${encodeURIComponent(search)}`)
      .then(res => res.json())
      .then(setQuestions);
  }, [search]);

  return (
    <div>
      <input
        placeholder="Search questions..."
        value={search}
        onChange={e => setSearch(e.target.value)}
      />
      <Link to="/questions/ask">Ask Question</Link>
      <ul>
        {questions.map(q => (
          <li key={q.id}>
            <Link to={`/questions/${q.id}`}>{q.title}</Link>
          </li>
        ))}
      </ul>
    </div>
  );
}
