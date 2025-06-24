import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';

export default function QuestionDetail() {
  const { id } = useParams();
  const [question, setQuestion] = useState(null);

  useEffect(() => {
    fetch(`/api/questions/${id}`)
      .then(res => res.json())
      .then(setQuestion);
  }, [id]);

  if (!question) return <div>Loading...</div>;

  return (
    <div>
      <h1>{question.title}</h1>
      <p>{question.body}</p>
      <button>Upvote</button>
      <button>Downvote</button>
      <h2>Answers</h2>
      {/* Render answers and comments */}
    </div>
  );
}
