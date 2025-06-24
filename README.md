# Stack Overflow-like Discussion Forum

A full-stack discussion built with:

- **Backend:** Golang (Gin framework), MongoDB
- **Frontend:** React

---

## Features

- User authentication (register, login, logout)
- Ask questions, post answers
- Add comments to answers
- Upvote/downvote questions and answers
- Search questions by keywords
- Responsive UI

---

## Project Structure

```
discussion-forum/
├── backend/
│   ├── main.go
│   ├── go.mod
│   ├── models/
│   ├── handlers/
│   └── ...
├── frontend/
│   ├── src/
│   ├── public/
│   ├── package.json
│   └── ...
└── README.md
```

---

## Backend Setup

### Prerequisites

- Go 1.20+
- MongoDB (local or Atlas)

### Install dependencies

```sh
cd backend
go mod tidy
```

### Environment Variables

Create a `.env` file in `backend/`:

```
MONGO_URI=mongodb://localhost:27017
```

Or export it in your shell:

```sh
export MONGO_URI="mongodb://localhost:27017"
```

### Run the server

```sh
go run main.go
```

The backend will start at [http://localhost:8080](http://localhost:8080).

---

## Frontend Setup

### Prerequisites

- Node.js 18+
- npm

### Install dependencies

```sh
cd frontend
npm install
```

### Proxy Setup

Ensure your `frontend/package.json` contains:

```json
"proxy": "http://localhost:8080"
```

### Run the frontend

```sh
npm start
```

The frontend will start at [http://localhost:3000](http://localhost:3000).

---

## Usage

1. **Register** a new user.
2. **Login** with your credentials.
3. **Ask questions**, **answer** others, **comment** on answers.
4. **Upvote/downvote** questions and answers.
5. **Search** for questions using the search bar.

---

## API Endpoints (Backend)

- `POST /api/register` — Register a new user
- `POST /api/login` — Login and receive JWT
- `GET /api/questions` — List/search questions
- `POST /api/questions` — Create a question (auth required)
- `GET /api/questions/:id` — Get question details
- `PUT /api/questions/:id/vote` — Upvote/downvote question (auth required)
- `POST /api/questions/:id/answers` — Post an answer (auth required)
- `PUT /api/answers/:id/vote` — Upvote/downvote answer (auth required)
- `GET /api/answers/:id/comments` — Get comments for an answer
- `POST /api/answers/:id/comments` — Add comment to answer (auth required)

---

## Notes

- JWT tokens are stored in `localStorage` on the frontend.
- Unique indexes are enforced for username and email in MongoDB.
- CORS is enabled by default in the backend.
- For production, configure environment variables and secure JWT secrets.

---

## License

MIT
