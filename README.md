# GoCalendar
GoCalendar is a backend microservice API built with Go, designed to streamline event management. This service provides a reliable and efficient solution for handling scheduling operations, making it a perfect addition to larger applications or as a standalone service. Integrate GoCalendar to enhance your application's scheduling capabilities with high performance and reliability.

### Table of Contents
- Features
- Getting Started
- API Routes
- Technologies Used

#### Features
GoCalendar offers comprehensive functionality for managing events, including creating, updating, cancelling, and deleting events.

#### Getting Started
To get started with GoCalendar, clone the repository and follow the setup instructions below.

##### Prerequisites
- Go 1.21 or later
- PostgreSQL

##### Installation
 1. Clone the repository:
`git clone https://github.com/yourusername/GoCalendar.git`
`cd GoCalendar`

2. Install dependencies:
`go mod tidy`

3. Set up your database and configure the connection in the .env file.

4. Run the application:
`go run main.go`

#### API Routes
##### Authentication
- **POST** `/api/create_user`: Register a new user.
- POST `/api/login`: Authenticate an existing user and obtain a token.
##### Event Management
- **GET** `/api/view_events`: Retrieve a list of events for the authenticated user.
- **POST** `/api/schedule_event`: Schedule a new event.
- **POST** `/api/cancel_event`: Cancel an existing event.
- **POST** `/api/update_event`: Update details of an existing event.

#### Technologies Used
- **Go**: The primary language for the backend service.
- **PostgreSQL**: For storing event data (can be replaced with another database).
- **JWT**: For secure user authentication.
