import React from 'react';
import { useParams } from 'react-router-dom';

function QuestionDetail() {
  const { id } = useParams();
  return (
    <div>
      <h2>Question Detail</h2>
      <p>Question ID: {id}</p>
      {/* Add question/answers/comments display here */}
    </div>
  );
}

export default QuestionDetail;
