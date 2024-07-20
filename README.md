# OAuth2 Server Implementation

This project is a minimal implementation of an OAuth2 server using Go and PostgreSQL, designed to follow the OAuth2 RFC specifications. It provides endpoints for client registration and will include additional OAuth2 functionality such as authorization, token management, and user information retrieval. The application is containerized using Docker Compose for easy setup and deployment.



## Features

- OAuth2 Registration Endpoint (More endpoints based on  OAuth2 RFC will be implemented)
- Built with Go
- PostgreSQL for data storage
- Docker Compose for easy setup

## Project Structure

- `main.go`: Entry point for the application.
- `store/entities`: Contains GORM models for database interaction.
- `handlers/`: Contains HTTP handlers for various endpoints.
- `docker-compose.yml`: Docker Compose configuration for running and PostgreSQL instance.

## Prerequisites

- Docker and Docker Compose installed on your machine.

## Getting Started

### Clone the Repository

```bash
git clone https://github.com/manuelrojas19/oauth2-golang-server
cd oauth2-golang-server
```

### Build and Run the Application

Build and start the database using Docker Compose

```bash
docker-compose up -d
```

### Run the Go application

```bash
go run main.go
``` 

### API Endpoints (Implemented)

#### Register: `/register`

- **Description**: This endpoint allows for the registration of OAuth2 clients.
- **Method**: `POST`
- **Request**:
    - **Headers**:
        - `Content-Type: application/json`
        - `Authorization: Bearer YOUR_ACCESS_TOKEN`
    - **Body**:
      ```json
      {
        "client_name": "Example Client",
        "grant_types": ["authorization_code", "refresh_token"],
        "response_types": ["code", "token"],
        "token_endpoint_auth_method": "client_secret_basic",
        "redirect_uris": ["https://example.com/callback"]
      }
      ```
- **Response**:
    - **Status**: `201 Created` (on successful registration)
    - **Body**:
      ```json
      {
        "client_id": "string",
        "client_secret": "string",
        "client_name": "Example Client",
        "grant_types": ["authorization_code", "refresh_token"],
        "response_types": ["code", "token"],
        "token_endpoint_auth_method": "client_secret_basic",
        "redirect_uris": ["https://example.com/callback"]
      }
      ```

#### Example `curl` Command

To register a new OAuth2 client, you can use the following `curl` command:

```bash
curl -X POST http://localhost:8080/register \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
    -d '{
        "client_name": "Example Client",
        "grant_types": ["authorization_code", "refresh_token"],
        "response_types": ["code", "token"],
        "token_endpoint_auth_method": "client_secret_basic",
        "redirect_uris": ["https://example.com/callback"]
    }'
```