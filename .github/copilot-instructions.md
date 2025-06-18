
This project is a monorepo with a Vite+React+TypeScript frontend and a Go backend. The backend exposes a REST API and connects to MongoDB. The frontend fetches data from the backend and displays it.

## General Development Guidelines

1. You cannot run commands in powershell with `&&` or `||` separators, you must use `;` to separate them instead.
  - So for example, instead of `npm install && npm run dev`, you should run `npm install; npm run dev`.
2. The step-by-step instructions for running the project are in the README. When changing the run process, please update the README accordingly.
3. The app is designed for mobile first (less than 440px wide), but should also work on desktop.
  - The biggest impact of this is on mobile the sidebars on the dashboard are hidden by default and can be toggled with buttons.
4. After making frontend changes, make sure to run `npm run lint` to check for linting issues. There should be no warnings or errors.
5. Style-wise, it's better to have multiple small components rather than one large component.
  - When a component grows too large, consider breaking it down into smaller components.

## The Application Layout

The app contains the following screens:
- Login - for existing users to log in
- New Bank - for registering a new player
- Dashboard - the main screen for logged-in users, showing their bank details and transactions
The dashboard is the main screen and is divided into several sections:
  - Bank Details - shows the player's bank details at the top
  - Bank information - shows the user's investments and stats in the middle
  - Investors - shows the list of investors at the left side
  - Settings - shows the settings at the right side
Much of this is still in development, so you may not see all features implemented yet. And there is still more to come.

## The Game

- The idea of the game is that a player logs in once a day to check their bank details, make investments, and manage their bank. 
- The player can also invest in other players' banks. The goal is to grow your bank and become the richest player.
- Each turn (day), the player has the opportunity to lie about how much money their bank has made. Ultimately this will lead to Ponzi schemes.
- If a player lies too much, it will increase their chances of being caught and fined, or even shut down completely.
- If a player lies too little, other players will not invest in their bank due to a low ROI.