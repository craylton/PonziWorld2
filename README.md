# Ponzi World

This is a multiplayer game where players are encouraged to create their own Ponzi schemes and scam each other out of as much money as possible!

Currently it's in very early development. 

## Prerequisites
- Go 1.18+
- MongoDB (running locally on port 27017 by default)
  - You can run MongoDB easily with Docker:
    ```powershell
    docker start ponzi-mongo
    ```

## Running the Backend

0. If running locally, ensure MongoDB is running (see above for Docker command)
1. Open a terminal in the `backend` directory.
2. Run the backend:
   ```powershell
   go run .
   ```

- The backend will listen on http://localhost:8080
- The MongoDB connection string can be set with the `MONGODB_URI` environment variable (defaults to `mongodb://localhost:27017`).

## Running the Frontend

To start the React UI locally:

1. Install dependencies (if you haven't already):
   ```powershell
   npm install
   ```
2. Start the development server:
   ```powershell
   npm run dev
   ```
3. Open your browser and go to the URL shown in the terminal (usually http://localhost:5173).
