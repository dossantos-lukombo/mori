
## Project Description
This project was part of my school's JavaScript course. We developed the project using Gitea (our school's chosen VCS), which is why the commit history is relatively sparse. My school colleagues and I built a Facebook-like social network with various social features.

## Features
- **Followers:** Follow and be followed by other users.
- **Profiles:** Create and manage personal profiles.
- **Posts:** Share posts with your network.
- **Groups:** Form and manage groups.
- **Notifications:** Stay updated with real-time notifications.
- **Chats:** Engage in private and group chats.
- **Dockerized Environment:** Easily run the project in a containerized setup.
- **Test Suite:** Comprehensive tests for server routes and functionalities.
- **Security:** Built-in SQL injection prevention using parameterized queries.

## Run the Project

### Prerequisites
- [NodeJS](https://nodejs.org/en/) installed.
- [Go](https://golang.org/) installed.
- [Docker](https://docs.docker.com/get-docker/) (optional, for containerized deployment).

### Development / Test Phase
#### Frontend
1. Navigate to the **/frontend** directory:
   cd frontend

    Install dependencies and start the development server:

    npm install
    npm run serve

    The Vue.js app will typically be available at http://localhost:8080.

Backend

    Navigate to the /backend directory:

cd backend

Start the backend server:

    go run server.go

    The Go server will run on port 8081.

LLM Server

    Navigate to the /backend/logicllm directory:

cd backend/logicllm

Start the LLM logic server with Uvicorn:

    uvicorn server:app --host 127.0.0.1 --port 8000

    This starts the LLM server on http://127.0.0.1:8000.

Docker (Optional)

You can also run the entire project using Docker.

    Building Docker Images:

        Backend:

cd backend
docker build -t mori-backend .

Frontend:

cd frontend
docker build -t mori-frontend .

LLM Server:

    cd backend/logicllm
    docker build -t mori-logicllm .

Running with Docker Compose:
Create a docker-compose.yml file in your project root (if not already present) and run:

    docker-compose up --build

    This will build and start all containers together.

What I Did

I designed the project using Figma, which served as the basis for the HTML and CSS. I developed the notifications and chat system, while the remaining frontend components were built in collaboration with my front-end partner, Vic.
Stack

Frontend:

    Vue.js

    HTML & CSS

    Figma

Backend:

    Go

    PostgreSQL

Security

Our backend code prevents SQL injection by using parameterized queries. All user inputs are passed as parameters (using placeholders like $1, $2, etc.) to the SQL queries. This approach ensures that inputs are safely escaped by the database driver, protecting the application from malicious SQL injection attacks.
Testing

A comprehensive suite of unit tests for server routes is included in server_test.go. To run the tests:

    Navigate to the /backend directory.

    Execute:

    go test -v

The test output includes a visual ASCII table showing the endpoint, expected output, and actual output for easy verification.
Authors

Backend Team:

    Quentin Boiteux

    Dos Santos

Frontend Team:

    Daryl Parisi

Additional Documentation

For more detailed information, refer to:

    docker.md: Instructions on building and running the project with Docker.

    security.md: Details on SQL injection prevention and overall security practices.

    TEST_PHASE.md: Guidelines on running the project in test mode.

    CONTRIBUTING.md: (If applicable) Guidelines for contributing to the project.


This README provides an overview of the project, instructions for running and testing in different environments, and additional details on security and the technology stack. Feel free to adjust paths and specifics to fit your project's setup.

