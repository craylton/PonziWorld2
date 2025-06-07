# Backend for PonziWorld2

This is the Go backend for the PonziWorld2 finance simulation/game. It exposes a REST API and connects to MongoDB.

## Prerequisites
- Go 1.18+
- MongoDB (running locally on port 27017 by default)

## Running the Backend

1. Open a terminal in this directory (`backend`).
2. Run the backend:
   ```powershell
   go run main.go
   ```

- The backend will listen on http://localhost:8080
- The MongoDB connection string can be set with the `MONGODB_URI` environment variable (defaults to `mongodb://localhost:27017`).

## API Endpoints
- `GET /api/hello` — Returns a JSON greeting string.

---

## Development Notes
- See `main.go` for implementation details.
