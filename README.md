# Go REST API with Session-Based Authentication

A secure and simple REST API for user management built with Go, Gin, and MongoDB. This project features a complete, session-based authentication system.

## 🛠 Tech Stack

- Go (1.21+)
- Gin Framework
- MongoDB
- Docker & Docker Compose
- `bcrypt` for password hashing

## 📁 Project Structure
```
go_user_api/
├── cmd/
│ └── main.go                # Entry point
├── internal/
│ ├── database/              # MongoDB connection & indexing
│ ├── handlers/              # HTTP request handlers & middleware
│ ├── models/                # Data structures (User, Session)
│ ├── routes/                # API route definitions
│ └── user/                  # User-specific database logic
├── .env.example
├── docker-compose.yml       # Docker services configuration
└── README.md

```

## 🚀 Quick Start

1.  **Clone or Download the Repository:**
    *   **Using Git (Recommended):**
        ```bash
        git clone https://github.com/Behzad1408/go_user_api.git
        cd go_user_api
        ```
    *   **Or Download ZIP:**
        *   Download and extract the ZIP file.
        *   Rename the extracted folder `go_user_api-main` to `go_user_api`.
        *   Navigate into the folder: `cd go_user_api`

2.  **Setup Environment Variables:**
    ```bash
    cp .env.example .env
    ```
    *   **Important:** Open the `.env` file and fill in your configuration, especially `MONGO_DB_NAME`.

3.  **Start Services:**
    ```bash
    docker-compose up -d
    ```
    *   This will start the MongoDB container and Mongo Express.

4.  **Run Application:**
    ```bash
    go mod tidy
    go run cmd/main.go
    ```
    *   The server will be running on `http://localhost:8080`.

## 📡 API Endpoints

All endpoints are prefixed with `/api/v1`.

| Method | Endpoint | Protection  | Description                                 |
| :----- | :------- | :---------- | :------------------------------------------ |
| `GET`  | `/health`| Public      | Checks if the API is running.               |
| `POST` | `/signup`| Public      | Registers a new user.                       |
| `POST` | `/login` | Public      | Authenticates a user and returns a session cookie.|
| `GET`  | `/me`    | **Protected** | Retrieves the data of the logged-in user.   |

### Example: Full Authentication Flow with `curl`

1.  **Sign Up a New User:**
    ```bash
    curl -X POST http://localhost:8080/api/v1/signup \
      -H "Content-Type: application/json" \
      -d '{"username":"behzad","email":"behzad@example.com","password":"password123"}'
    ```

2.  **Login to Get the Session Cookie:**
    *   This command saves the cookie to a file named `cookies.txt`.
    ```bash
    curl -X POST http://localhost:8080/api/v1/login \
      -H "Content-Type: application/json" \
      -d '{"email":"behzad@example.com","password":"password123"}' \
      -c cookies.txt
    ```

3.  **Access the Protected Route:**
    *   This command uses the cookie from `cookies.txt` to authenticate.
    ```bash
    curl -X GET http://localhost:8080/api/v1/me -b cookies.txt
    ```

## 🔧 Environment Variables

Copy `.env.example` to `.env` and configure your settings:
```env
MONGO_USER=myuser
MONGO_PASSWORD=mypass
MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_DB_NAME=mydb
SERVER_PORT=8080
```

## 🗄 MongoDB Access

- **Mongo Express UI:** http://localhost:8081
- **MongoDB:** mongodb://localhost:27017

## 📄 License

This project is licensed under the MIT License.

---