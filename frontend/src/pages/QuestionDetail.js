import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';

function QuestionDetail() {
  const { id } = useParams();
  const [question, setQuestion] = useState(null);
  const [answers, setAnswers] = useState([]);
  const [answerBody, setAnswerBody] = useState('');
  const [loading, setLoading] = useState(true);
  const [comments, setComments] = useState({}); // { answerId: [comments] }
  const [commentInputs, setCommentInputs] = useState({}); // { answerId: commentText }

  useEffect(() => {
    const fetchQuestion = async () => {
      setLoading(true);
      try {
        const res = await fetch(`/api/questions/${id}`);
        const data = await res.json();
        setQuestion(data.question);
        setAnswers(data.answers || []);
      } catch (err) {
        // ...existing code...
      }
      setLoading(false);
    };
    fetchQuestion();
  }, [id]);

  useEffect(() => {
    // Fetch comments for each answer
    const fetchComments = async () => {
      let allComments = {};
      for (const a of answers) {
        try {
          const res = await fetch(`/api/answers/${a.id || a._id}/comments`);
          const data = await res.json();
          allComments[a.id || a._id] = Array.isArray(data) ? data : [];
        } catch (err) {
          allComments[a.id || a._id] = [];
        }
      }
      setComments(allComments);
    };
    if (answers.length > 0) fetchComments();
  }, [answers]);

  const submitAnswer = async (e) => {
    e.preventDefault();
    const token = localStorage.getItem('token');
    try {
      const res = await fetch(`/api/questions/${id}/answers`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': token ? `Bearer ${token}` : '',
        },
        body: JSON.stringify({ body: answerBody }),
      });
      if (res.ok) {
        setAnswerBody('');
        // reload answers
        const data = await res.json();
        setAnswers(prev => [...prev, data]);
      } else {
        const data = await res.json();
        alert('Failed to post answer: ' + (data.error || 'Unknown error'));
      }
    } catch (err) {
      alert('Failed to post answer: Network error');
    }
  };

  const voteAnswer = async (answerId, up) => {
    const token = localStorage.getItem('token');
    try {
      const res = await fetch(`/api/answers/${answerId}/vote`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': token ? `Bearer ${token}` : '',
        },
        body: JSON.stringify({ up }),
      });
      if (res.ok) {
        // Get updated answer from backend
        const updated = await res.json();
        setAnswers(answers =>
          answers.map(a =>
            (a.id || a._id) === answerId
              ? { ...a, votes: updated.votes }
              : a
          )
        );
      }
    } catch (err) {
      alert('Failed to vote: Network error');
    }
  };

  const handleCommentInput = (answerId, value) => {
    setCommentInputs(inputs => ({ ...inputs, [answerId]: value }));
  };

  const submitComment = async (answerId) => {
    const token = localStorage.getItem('token');
    const commentText = commentInputs[answerId];
    if (!commentText) return;
    try {
      const res = await fetch(`/api/answers/${answerId}/comments`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': token ? `Bearer ${token}` : '',
        },
        body: JSON.stringify({ body: commentText }),
      });
      if (res.ok) {
        const newComment = await res.json();
        setComments(comments => ({
          ...comments,
          [answerId]: [...(comments[answerId] || []), newComment],
        }));
        setCommentInputs(inputs => ({ ...inputs, [answerId]: '' }));
      } else {
        alert('Failed to add comment');
      }
    } catch (err) {
      alert('Failed to add comment: Network error');
    }
  };

  if (loading) return <div>Loading...</div>;
  if (!question) return <div>Question not found.</div>;

  return (
    <div>
      <div style={{ background: '#fff', border: '1px solid #e4e6e8', borderRadius: 4, padding: 24, marginBottom: 24 }}>
        <h2 style={{ fontWeight: 500, fontSize: 24 }}>{question.title}</h2>
        <div style={{ color: '#232629', margin: '16px 0' }}>{question.body}</div>
        <div>
          {(question.tags || []).map(tag => (
            <span key={tag} style={{ background: '#e1ecf4', color: '#39739d', borderRadius: 3, padding: '2px 8px', marginRight: 8, fontSize: 13 }}>{tag}</span>
          ))}
        </div>
      </div>
      <h3 style={{ fontWeight: 400, fontSize: 20, marginBottom: 16 }}>{answers.length} Answers</h3>
      <div>
        {answers.map(a => (
          <div key={a.id || a._id} style={{ background: '#fff', border: '1px solid #e4e6e8', borderRadius: 4, padding: 16, marginBottom: 16 }}>
            <div style={{ color: '#232629', marginBottom: 8 }}>{a.body}</div>
            <div style={{ color: '#6a737c', fontSize: 13, display: 'flex', alignItems: 'center', gap: 8 }}>
              Votes: {a.votes || 0}
              <button onClick={() => voteAnswer(a.id || a._id, true)} style={{ marginLeft: 8, background: '#e1ecf4', border: 'none', borderRadius: 3, cursor: 'pointer' }}>Like</button>
              <button onClick={() => voteAnswer(a.id || a._id, false)} style={{ background: '#fde3e1', border: 'none', borderRadius: 3, cursor: 'pointer' }}>Dislike</button>
            </div>
            <div style={{ marginTop: 12 }}>
              <b>Comments:</b>
              <ul style={{ paddingLeft: 20 }}>
                {(comments[a.id || a._id] || []).map(c => (
                  <li key={c.id || c._id} style={{ fontSize: 14, color: '#444' }}>{c.body}</li>
                ))}
              </ul>
              <form
                onSubmit={e => {
                  e.preventDefault();
                  submitComment(a.id || a._id);
                }}
                style={{ marginTop: 8, display: 'flex', gap: 8 }}
              >
                <input
                  type="text"
                  value={commentInputs[a.id || a._id] || ''}
                  onChange={e => handleCommentInput(a.id || a._id, e.target.value)}
                  placeholder="Add a comment..."
                  style={{ flex: 1, padding: 4, borderRadius: 3, border: '1px solid #ccc' }}
                />
                <button type="submit" style={{ background: '#f48024', color: '#fff', border: 'none', borderRadius: 3, padding: '0 12px', fontWeight: 600 }}>Add</button>
              </form>
            </div>
          </div>
        ))}
      </div>
      <form onSubmit={submitAnswer} style={{ background: '#fff', border: '1px solid #e4e6e8', borderRadius: 4, padding: 16, marginTop: 32 }}>
        <h4 style={{ fontWeight: 400, fontSize: 18, marginBottom: 12 }}>Your Answer</h4>
        <textarea
          value={answerBody}
          onChange={e => setAnswerBody(e.target.value)}
          required
          rows={5}
          style={{ width: '100%', padding: 8, borderRadius: 4, border: '1px solid #ccc', marginBottom: 12 }}
        />
        <button type="submit" style={{ background: '#f48024', color: '#fff', padding: '0.5rem 1.5rem', border: 'none', borderRadius: 4, fontWeight: 600, fontSize: 16 }}>
          Post Your Answer
        </button>
      </form>
    </div>
  );
}

export default QuestionDetail;
