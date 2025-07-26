# PonziWorld2 Backend Architecture

This document provides an overview of the backend architecture for the PonziWorld2 application. The backend is written in Go and follows a layered architecture pattern to separate concerns and improve maintainability.

## Project Structure

The backend code is organized into the following directories:

-   `main.go`: The entry point of the application. It initializes the database connection, sets up the dependency injection container, configures the routes, and starts the web server.
-   `auth/`: Contains the logic for handling JWT (JSON Web Token) based authentication. `jwt.go` is responsible for generating and validating tokens.
-   `config/`: This directory is responsible for setting up the dependency injection (DI) container.
    -   `container.go`: The main container that holds all the application's dependencies.
    -   `repository_container.go`: Initializes and configures all the repositories.
    -   `service_container.go`: Initializes and configures all the services, injecting repository dependencies.
-   `database/`: Manages the MongoDB database connection.
    -   `database.go`: Establishes the connection to the database.
    -   `indexes.go`: Creates necessary indexes in the database for performance.
-   `handlers/`: These are the controllers in an MVC pattern. They handle incoming HTTP requests, parse request bodies and parameters, call the appropriate services, and return HTTP responses. Each model has its own handler file.
-   `middleware/`: Contains HTTP middleware.
    -   `auth.go`: An authentication middleware that protects routes by validating JWTs.
    -   `cors.go`: Handles Cross-Origin Resource Sharing (CORS) configuration.
-   `models/`: Defines the core data structures (structs) used throughout the application, representing entities like `Player`, `Bank`, `Asset`, etc.
-   `repositories/`: The data access layer. Repositories are responsible for all communication with the database (queries, inserts, updates, deletes). They abstract the database implementation from the rest of the application. `interfaces.go` defines the interfaces that the repositories implement.
-   `requestcontext/`: Defines a custom context that can be passed through the application, for example to carry user information from the authentication middleware.
-   `routes/`: Defines all the API routes for the application. `routes.go` maps the HTTP endpoints to their corresponding handlers and applies middleware where necessary.
-   `services/`: This layer contains the core business logic of the application. Services are called by handlers, and they coordinate the application's response to a request. They use repositories to interact with the database.
-   `tests/`: Contains all the integration and unit tests for the backend.

## Request Lifecycle

A typical request flows through the application as follows:

1.  An HTTP request hits an endpoint defined in `routes/routes.go`.
2.  The request may pass through one or more `middleware` functions (e.g., for logging, CORS, or authentication). The authentication middleware, if the route is protected, will validate the JWT and may add user information to the request context.
3.  The route dispatches the request to a specific `handler` function.
4.  The `handler` parses the request (e.g., decoding a JSON body into a `model` struct).
5.  The `handler` calls one or more `service` methods to perform the required business logic.
6.  The `service` uses `repository` methods to interact with the database.
7.  The `repository` executes queries against the MongoDB database and returns data, which is mapped to the `models`.
8.  The `service` processes the data from the repository and returns it to the `handler`.
9.  The `handler` creates an HTTP response (e.g., a JSON response with a status code) and sends it back to the client.

## Dependency Injection

The application uses a dependency injection container (`config/container.go`) to manage the lifecycle and dependencies of services and repositories. This promotes loose coupling and makes the code easier to test and maintain. Services are initialized with their required repository dependencies.

## Running Tests

To ensure the application is working correctly, you can run the test suite from the `backend` directory:

```shell
go test -v ./tests/
```
