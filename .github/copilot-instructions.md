
This project is a monorepo with a Vite+React+TypeScript frontend and a Go backend. The backend exposes a REST API and connects to MongoDB. The frontend fetches data from the backend and displays it.

## General Development Guidelines

1. You cannot run commands in powershell with `&&` or `||` separators, you must use `;` to separate them instead.
- So for example, instead of `npm install && npm run dev`, you should run `npm install; npm run dev`.
2. The step-by-step instructions for running the project are in the README. When changing the run process, please update the README accordingly.
3. The app is designed to be run on mobile first (less than 440px wide), but should also work on desktop.
4. After making any frontend changes, make sure to run `npm run lint` to check for linting issues. There should be no warnings or errors.
5. Style-wise, it's better to have multiple small components rather than one large component. When a component grows too large, consider breaking it down into smaller components.

## The Application
The app contains the following screens:
- Login - for existing users to log in
- New Bank - for registering a new player
- Dashboard - the main screen for logged-in users, showing their bank details and transactions
The dashboard is the main screen and is divided into several sections:
  - Bank Details - shows the player's bank details at the top
  - Important Information - shows the important information in the middle
  - Investors - shows the list of investors at the left side
  - Investments - shows the list of investments at the right side
Much of this is still in development, so you may not see all features implemented yet. And there is still more to come.