This project is a monorepo with a Vite+React+TypeScript frontend and a Go backend. The backend exposes a REST API and connects to MongoDB. The frontend fetches data from the backend and displays it.

## General Development Guidelines

1. You cannot run commands in powershell with `&&` or `||` separators, you must use `;` to separate them instead.
  - So for example, instead of `npm install && npm run dev`, you should run `npm install; npm run dev`.
2. The step-by-step instructions for running the project are in the README.
3. The app is designed for mobile first (less than 900px wide), but should also work on desktop.
  - The biggest impact of this is on mobile the sidebars on the dashboard are hidden by default and can be toggled with buttons.
  - Also there is little need for `:hover` effects since this won't be seen on mobile.
4. After making frontend changes, make sure to run `npm run lint` to check for linting issues. There should be no warnings or errors.
5. Style-wise, it's better to have multiple small components rather than one large component.
  - When a component grows too large, consider breaking it down into smaller components.
6. After making any code changes, review your own changes thoroughly and amend them if necessary.
  - When reviewing your own changes you might realise you've missed something or that you can still improve/tidy your code.
7. After making backend changes, run the backend tests to ensure everything is still working.
  - You can run the tests with `go test -v ./tests/` in the `backend` directory.
8. Do not implement a proposed solution that you believe to be suboptimal or not in line with best practice.
  - Instead, make your suggestion and see what the team thinks.
  - This is to avoid introducing potential issues or technical debt.
  - If unsure on certain details, ask for clarification or guidance.

## The Application Layout

The app contains the following screens:
- Login - for existing players to log in
- New Bank - for registering a new player
- Dashboard - the main screen for logged-in players, showing their bank details and transactions
  - Header - shows the player's bank details at a glance
  - Assets - shows the player's investments in the middle of the screen
    - Clicking an asset will show a popup with more details
    - After clicking an asset, the player has the option to buy or sell that asset
  - Investors - shows the list of investors at the left side
  - Settings - shows the settings at the right side
Much of this is still in development, so you may not see all features implemented yet. And there is still more to come.

## The Game

- The game is a 'daily' game, meaning players log in once a day to manage their banks
- Players are able to invest in one another's banks - as the game is multiplayer
- Graphics are minimal, focusing on functionality and gameplay
