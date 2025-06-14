
This project is a monorepo with a Vite+React+TypeScript frontend and a Go backend. The backend exposes a REST API and connects to MongoDB. The frontend fetches data from the backend and displays it.

1. You cannot run commands in powershell with `&&` or `||` separators, you must use `;` to separate them instead.
So for example, instead of `npm install && npm run dev`, you should run `npm install; npm run dev`.

2. The step-by-step instructions for running the project are in the README. When changing the run process, please update the README accordingly.

3. The app is designed to be run on mobile first (less than 1080px wide), but should also work on desktop.

4. After making any frontend changes, make sure to run `npm run lint` to check for linting errors.