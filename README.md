# OAuth2 Server Implementation

This project is a basic implementation of an OAuth2 server based on the OAuth2 RFC. It is built using Go and PostgreSQL.
Currently, it only includes the registration endpoint.

## Features

- OAuth2 Registration Endpoint
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
git clone <your-repo-url>
cd <your-repo-directory>
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