# Frontend TypeScript Code Style Guide

Readability comes first. Aim for simple, boring code that others can scan quickly and understand without surprises.

## Core principles

- Simplicity: Choose the simplest thing that works. Avoid cleverness and unnecessary abstraction.
- Modernism: Follow current TypeScript/React practices, using efficient and up to date approaches.
- Clean-code principles: Favor small functions, clear names, low nesting, and straightforward control flow.
- Minimalism: If it is not needed, leave it out. Remove dead code and keep files focused.
- Consistency: The entire application should look like it was designed and written by one person.

## Files

- Each file should focus on a single component, hook, or core concept (utility, type, or service). If a file contains multiple concepts with their own responsibilities, split them into separate files.
- Files should be small. If a file is longer than approximately 100 lines (excluding imports and type/interface definitions), it should be split into two or more files.
- Files should be simple and focused. If a file seems to have multiple concerns, split it.
- Types and interfaces should be defined near where they are used. Export them only when they are shared across modules.
- Ensure each file is in the correct feature/domain folder. If there are many files in the same directory, consider creating a new subdirectory.

## Naming and APIs

- Use standard TypeScript/JavaScript naming conventions by default.
- Use standard initialisms (ID, URL, JSON) consistently.
- Prefer full words over abbreviations.
- Boolean variable names should start with a prefix like `isLoaded`, `hasData`, `shouldRender`, or `canSubmit`.
- Function names should describe what the function does. Variable names should describe what the variable represents.
- React components use PascalCase (e.g., `BankHeader`). Component files should also be PascalCase.
- React hooks must start with `use` and be camelCase (e.g., `useBankDetails`).
- Event handlers and callbacks should be named with a clear verb (e.g., `handleSubmit`, `onCloseRequested`).
- Do not prefix interfaces with `I`. Prefer descriptive names like `Player`, `Bank`, `PlayerProps`.

## Functions

- Write small functions with one job (single responsibility principle).
- Use early returns to keep code flat and readable.
- If a function is not straightforward, or it is not small (~30 lines as a rough guide), split it.
- Keep an eye out for logic that is unrelated to the enclosing component/hook. It probably belongs in a separate utility module.
- Avoid "pass-through" functions or components which simply forward to another function without adding value.
- Prefer pure functions for domain/business logic; keep React components thin and focused on rendering and orchestrating.
- Prefer `async/await` for async flows; always handle errors and avoid floating promises.

## Testability

- Each file should be easily unit-testable. Isolate logic in pure functions and small hooks.
- Use dependency injection via function parameters/props rather than importing singletons where practical, to make mocking easier.
- Move side effects to the edges (e.g., inside `useEffect`) and keep them small and well-contained.
- Keep components controlled via props when possible; avoid hidden state unless needed.

## General

- Error handling: Use typed errors or discriminated unions when useful. Provide safe, user-friendly messages and never log secrets.
- Logging: Use a shared logger utility if present. Remove stray `console.log`/`console.error` in production code.
- Async and cancellation: When calling `fetch` or long-running operations, accept and honor an `AbortSignal` and clean up effects on unmount.
- Data immutability: Prefer immutable updates for state. Use spread syntax or helper utilities; never mutate React state directly.
- Iteration: Favor array methods (`map`, `filter`, `reduce`) or `for...of` over index-based `for` loops unless there is a clear performance reason.
- Types: Prefer explicit types on public module boundaries (exports). Rely on inference inside function bodies. Avoid `any` and excessive type assertions; prefer `unknown` over `any`.
- Nullish values: Be consistent with `undefined` vs `null`. Prefer `undefined` unless an API requires `null`. Use optional chaining (`?.`) and nullish coalescing (`??`) thoughtfully.
- Imports: Keep imports ordered and tidy: stdlib/browser, third-party, then internal modules. Avoid circular dependencies.

## Comments and tooling

- Generally speaking, write comments only if the "why" cannot be made clearer in code.
- If the intent of the code is not obvious, refactor toward clearer code rather than adding comments.
- Use Prettier for formatting and ESLint for linting. Prefer fixing code over suppressing lints.
- Keep `tsconfig` strict settings enabled where possible. Avoid disabling type checks unless absolutely necessary and localized.
- Use TSDoc/JSDoc for public utilities and shared types where it materially improves understanding.

## A quick pre-PR check

- Names are clear and friendly; no unnecessary abbreviations or stutter.
- Functions and components are small with early returns and minimal branching.
- No stray logs; no secrets in code or logs.
- Props, state, and public module boundaries are well-typed; no accidental `any`.
- Async code handles errors and supports cancellation where appropriate.
- Dead code is removed; files stay focused on one purpose.
- ESLint passes with no warnings; formatting is consistent.
