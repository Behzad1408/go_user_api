# Go REST API with MongoDB

A simple REST API for user management built with Go, Gin Framework, and MongoDB.

## 🛠 Tech Stack

- Go 1.21+
- Gin Framework
- MongoDB
- Docker & Docker Compose
- bcrypt for password hashing

## 📁 Project Structure
```
my-go-project/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── database/            # MongoDB connection
│   ├── models/              # Data models
│   ├── handlers/            # Request handlers
│   └── routes/              # API routes
├── docker-compose.yml       # Docker services configuration
├── .env.example
└── README.md
```

## 🚀 Quick Start

1. **Clone and setup:**
```bash
   git clone https://github.com/Behzad1408/go_user_api.git
   cd Go_Project
   cp .env.example .env
```

2. **Start MongoDB:**
```bash
   docker-compose up -d
```

3. **Run application:**
```bash
   go run cmd/main.go
```

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/health` | Health check |
| POST | `/api/v1/signup` | User registration |

### Example: User Registration
```bash
curl -X POST http://localhost:8080/api/v1/signup \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"123456"}'
```

## 🔧 Environment Variables

Copy `.env.example` to `.env` and configure:
```env
MONGO_USERNAME=your_username
MONGO_PASSWORD=your_password
MONGO_HOST=localhost
MONGO_PORT=27017
APP_PORT=8080
```

## 🗄 MongoDB Access

- **Mongo Express UI:** http://localhost:8081
- **MongoDB:** mongodb://localhost:27017

## 📄 License

MIT

---

Made with ❤️ using Go