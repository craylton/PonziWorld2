
## General Development Guidelines

1. This project uses the following technologies:
   - Frontend: Vite, React, TypeScript, CSS. Backend: Go, MongoDB
   - The frontend is in the `src` directory and the backend is in the `backend` directory.
2. The app is designed for mobile first (less than 900px wide), but should also work on desktop.
  - The biggest impact of this is on mobile the sidebars on the dashboard are hidden by default and can be toggled with buttons.
  - Also there is little need for `:hover` effects since these won't be seen on mobile.
3. After making frontend changes, make sure to run `npm run lint` to check for linting issues. There should be no warnings or errors.
4. After making backend changes, run the backend tests to ensure everything is still working.
  - You can run the tests with `go test -v ./tests/` in the `backend` directory.
  - Tip: before navigating to the backend directory, check if you're already there!
5. Before implementing a new feature, ask yourself if this is the best way to do it.
  - If you think there might be a better way, discuss it with the team first.
  - If unsure on certain details, ask for clarification or guidance.
6. If you need guidance on how the code is structured, there are a few README files around the place which contain useful information.

## Code Style

1. Style-wise, it's better to have multiple small components rather than one large component.
  - When a component grows too large, consider breaking it down into smaller components.
2. This one's important! After making any code changes, review your own changes thoroughly and amend them if necessary.
  - When reviewing your own changes you might realise you've missed something or that you can still improve/tidy your code.
3. Aim for really simple code.
  - If you notice any duplicated code, try to consolidate.
  - If you notice any unnecessarily complex code, try to simplify it.
  - If you notice any code which isn't used, remove it.
4. Generally speaking, comments aren't used unless absolutely necessary.

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
