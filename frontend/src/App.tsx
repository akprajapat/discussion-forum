import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import Navbar from './components/Navbar';
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import QuestionDetail from './pages/QuestionDetail';
import AskQuestion from './pages/AskQuestion';

function App() {
  return (
    <Router>
      <div className="App">
        <Navbar />
        <div className="container mx-auto px-4">
          <Switch>
            <Route exact path="/" component={Home} />
            <Route path="/login" component={Login} />
            <Route path="/register" component={Register} />
            <Route path="/questions/ask" component={AskQuestion} />
            <Route path="/questions/:id" component={QuestionDetail} />
          </Switch>
        </div>
      </div>
    </Router>
  );
}

export default App;
