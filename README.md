# Organization API

This repository houses a Golang-based API application designed for managing organizations. The application includes features such as token management, CRUD operations for organizations, user invitations, and integration with MongoDB using Docker.

## Setup

To run the application locally, you need to have docker installed on your machine and make (if you don't have make check the Makefile for `compose-build`)
- `make compose-build`

## Documentation

Check out the postman documentations [here](https://documenter.getpostman.com/view/20745767/2s9YyzdJDh)

## Notes

- I tried my best to follow the given structure, even though I have some takes on it like we can merge handlers and controllers into single `service`.
- The application is very simple and it's not meant for production use.
- I will refactor the code but as I don't have time for this now, all I care about is finishing the project before time (later on I will create branches satisfying the new changes)
- As you have noticed I have mongo-express to see the database

![](./.assests/mongo-express.png)

## Project Structure

The project structure is designed to assist you in getting started quickly. You can modify it as needed for your specific requirements.

- **cmd/**: Contains the main application file.
  - **main.go**: The entry point of the application.

- **pkg/**: Core logic of the application divided into different packages.
  - **api/**: API handling components.
    - **handlers/**: API route handlers.
    - **middleware/**: Middleware functions.
    - **routes/**: Route definitions.
  - **controllers/**: Business logic for each route.
  - **database/**: Database-related code.
    - **mongodb/**
      - **models/**: Data models.
      - **repository/**: Database operations.
  - **utils/**: Utility functions.
  - **app.go**: Application initialization and setup.

- **docker/**: Docker-related files.
  - **Dockerfile**: Instructions for building the application image.

- **docker-compose.yaml**: Configuration for Docker Compose.

- **config/**: Configuration files for the application.
  - **app-config.yaml**: General application settings.
  - **database-config.yaml**: Database connection details.

- **tests/**: Directory for tests.
  - **e2e/**: End-to-End tests.
  - **unit/**: Unit tests.

- **.gitignore**: Specifies files and directories to be ignored by Git.

## Todo
- [ ] Add production branch
- [ ] Refactor code
