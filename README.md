# OAuth2 Server Implementation

This project is a minimal implementation of an OAuth2 server using Go and PostgreSQL, designed to follow the OAuth2 RFC
specifications. It provides endpoints for client registration and will include additional OAuth2 functionality such as
authorization, token management, and user information retrieval. The application is containerized using Docker Compose
for easy setup and deployment.

## Features

- OAuth2 Registration Endpoint (More endpoints based on OAuth2 RFC will be implemented)
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

#### Authorize: `/Authorize`

- **Description**: This endpoint handles the authorization of registered OAuth2 clients.
- **Method**: `POST`
- **Request**:
    - **Headers**:
        - `Content-Type: application/json`
    - **Body**:
      ```json
      {
        "client_id": "string",
        "redirect_uri": "https://example.com/callback",
        "response_type": "code",
        "scope": "read write",
        "state": "string"
      }
      ```
- **Response**:
    - **Status**: `302 Found` (on successful authorization)
    - **Headers**
        - `Location: https://example.com/callback?code=authorization_code&state=string`
    
#### Token: `/token`

- **Description**: This endpoint exchanges an authorization code or refresh token for an access token, or handles other
  token-related requests.
- **Method**: `POST`
- **Request**:
    - **Headers**:
        - `Content-Type: application/x-www-form-urlencoded`

    - **Query Params:**
        - For `client_credentials` grant type:
          ```
          grant_type=authorization_code&code=AUTH_CODE&redirect_uri=https://your-app.com/callback&client_id=CLIENT_ID&client_secret=CLIENT_SECRET
          ```
        - For `client_credentials` grant type:
          ```
          grant_type=client_credentials&client_id=CLIENT_ID&client_secret=CLIENT_SECRET
          ```
        - For `refresh_token` grant type:
          ```
          grant_type=refresh_token&client_id=CLIENT_ID&client_secret=CLIENT_SECRET
          ```
- **Response**:
    - **Status**: `200 Okey` (on successful token issuance)
    - **Body**:
      ```json
      {
        "access_token": "string",
        "token_type": "bearer",
        "expires_in": 3600,
        "refresh_token": "string",
        "scope": "string" 
      }
      ```

#### Example `curl` Command

To register a new OAuth2 client, you can use the following `curl` command:

```bash
curl -X POST http://localhost:8080/register \
    -H "Content-Type: application/json" \
    -d '{
        "client_name": "Example Client",
        "granttype": ["authorization_code", "refresh_token"],
        "responsetype": ["code", "token"],
        "token_endpoint_auth_method": "client_secret_basic",
        "redirect_uris": ["https://example.com/callback"]
    }'
```

To get an access token, you can use the following `curl` command:

- `authorization_code` Grant Type

```bash
curl -X POST http://localhost:8080/oauth/token \
  -d "code=NDU2ZWFLNMITN2M5NS0ZODBJLWFHNZITY2UWYZMWZWRMZDK2" \
  -d "client_id=d1734396-476f-4a40-95b3-4349524be2e4" \
  -d "client_secret=61a7f6ac-e922-46ad-9b27-f49713d9a4a2" \
  -d "redirect_uri=https://your-app.com/callback" \
  -d "grant_type=authorization_code"
```

- `client_credentials` Grant Type

```bash
curl -X POST http://localhost:8080/token \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "grant_type=client_credentials&client_id=CLIENT_ID&client_secret=CLIENT_SECRET"
```

- `refresh_token` Grant Type

```bash
curl -X POST http://localhost:8080/token \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "grant_type=refresh_token&refresh_token=REFRESH_TOKEN&client_id=CLIENT_ID&client_secret=CLIENT_SECRET"
```