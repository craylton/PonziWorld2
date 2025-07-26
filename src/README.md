# Frontend Architecture

This document provides a high-level overview of the frontend architecture for the PonziWorld2 application.

## Tech Stack

- **Framework**: [React](https://reactjs.org/)
- **Build Tool**: [Vite](https://vitejs.dev/)
- **Language**: [TypeScript](https://www.typescriptlang.org/)
- **Styling**: CSS

## Folder Structure

The frontend code is located in the `src` directory. Here is a breakdown of the key folders and their purposes:

-   **`/components`**: Contains small, reusable UI components that are shared across different parts of the application (e.g., `Popup.tsx`).

-   **`/contexts`**: Manages global state using React's Context API. This is our primary method for state management, avoiding the need for a larger library like Redux.
    -   Each context has a provider (e.g., `BankContext.tsx`) and a custom hook (e.g., `useBankContext.ts`) for easy consumption in components.

-   **`/models`**: Defines TypeScript interfaces and types for all the data structures used in the application (e.g., `Player.ts`, `Asset.ts`). This helps ensure type safety and consistency.

-   **`/utils`**: Holds utility functions that can be used throughout the application.

-   **Feature-based folders (e.g., `/Dashboard`, `/Login`)**: Most of the application is structured by features. Each feature folder contains the components, styles, and logic specific to that feature. For example, the `/Dashboard` directory contains all components that make up the main player dashboard.

## Key Files

-   **`main.tsx`**: The entry point of the application. It renders the root `App` component.

-   **`App.tsx`**: The main component of the application. It sets up the routing and wraps the application with necessary context providers.

-   **`ProtectedRoute.tsx`**: A higher-order component used to protect routes that require user authentication. It ensures that only logged-in players can access certain pages, like the Dashboard.

-   **`auth.ts`**: Contains authentication-related functions, such as handling user login, logout, and managing authentication tokens.

## How It All Works

1.  **Initialization**: The app starts with `main.tsx`, which renders `App.tsx`.
2.  **State Management**: `App.tsx` wraps the application in context providers (from `/contexts`). These providers make global data, like bank and asset information, available to all child components.
3.  **Routing**: The application uses `react-router-dom` (configured in `App.tsx`) to navigate between different screens like Login, New Bank, and the Dashboard.
4.  **Authentication**: When a user tries to access a protected route, `ProtectedRoute.tsx` checks their authentication status using logic from `auth.ts`. If not logged in, they are redirected to the `/login` page.
5.  **Components**: Screens are built by composing components. Feature-specific components are located within their respective feature folders, and shared components are in `/components`. Components use the custom context hooks (e.g., `useBankContext`) to access and manipulate global state.
6.  **Data Models**: When fetching data from the backend, the application uses the TypeScript models from the `/models` directory to ensure the data is correctly typed.
