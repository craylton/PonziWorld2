# Ponzi World

## Game Summary

Welcome to Ponzi World! You are now an investment bank manager, and you are tasked with making as much money as possible! 🤑

But, the stock market just crashed and you've lost a fortune. Your investors aren't going to be happy. 😬 Once you tell them about this crash, they'll withdraw all their investments and your bank will go bust 📉 If  only there was some way of, uh, 'bending' the truth a bit?

This is a daily multiplayer game (very much still in early development) where players are encouraged to create their own Ponzi schemes and scam each other out of as much money as possible!

The basic idea is:
- Once per day you log in and manage your investments.
- Each day, many other players will also log in and manage their own investments.
- Some of those players will invest in your bank.
- You will be asked how much money you made yesterday. Thanks to your crooked accountant, telling the truth is optional!
- If you tell the truth, your bank might not look very promising to investors, so they won't invest much.
- If you lie too much, people might become suspicious and you could get reported, fined, and maybe even shut down completely!

## Running the game

### Prerequisites
- Go 1.18+
- MongoDB (running locally on port 27017 by default)
  - You can run MongoDB easily with Docker:
    ```powershell
    docker start ponzi-mongo
    ```

### Running the Backend

0. If running locally, ensure MongoDB is running (see above for Docker command)
1. Open a terminal in the `backend` directory.
2. Run the backend:
   ```powershell
   go run .
   ```

- The backend will listen on http://localhost:8080
- The MongoDB connection string can be set with the `MONGODB_URI` environment variable (defaults to `mongodb://localhost:27017`).

### Running the Frontend

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

## API Endpoints

### Authentication
- `POST /api/user` - Create a new user and bank
- `POST /api/login` - Login and get JWT token

### Bank Management
- `GET /api/bank` - Get user's bank details and assets
- `GET /api/performanceHistory/ownbank/{bankId}` - Get 30 days of performance history for a bank you own

### Notes
- All authenticated endpoints require a Bearer token in the Authorization header
- Game operates on a daily cycle (currently using day 0 as reference point)
