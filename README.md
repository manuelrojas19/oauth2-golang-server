# Go OAuth2 Server Implementation

This project is an implementation of an OAuth2 server using Go and PostgreSQL, designed to follow the OAuth2 RFC
specifications. It provides endpoints for client registration and will include additional OAuth2 functionality such as
authorization, token management, and user information retrieval. The application is containerized using Docker Compose
for easy setup and deployment.

## Project Architecture

The server follows a layered architecture to ensure separation of concerns and maintainability:

-   **`handlers/`**: Contains HTTP handlers responsible for processing incoming requests, decoding payloads, calling
    service-layer logic, and sending HTTP responses. Handles API-specific error translation.
-   **`services/`**: Implements the core business logic of the OAuth2 server. It orchestrates operations between
    repositories, handles complex validations, and manages the overall OAuth2 flows (e.g., authorization, token issuance).
-   **`store/` and `store/repositories/`**: Manages data persistence. The `store` package defines data models (entities),
    while `store/repositories` defines interfaces and implementations for database interactions (e.g., saving clients,
    fetching scopes).
-   **`oauth/`**: Houses the core OAuth2 specification-related types and logic, such as `Client`, `Token`, `AuthCode`
    structs, and OAuth2 specific enums (grant types, response types, auth methods).
-   **`api/`**: Defines request and response structures for the public API, as well as common API error definitions.
-   **`configuration/`**: Handles application configuration, including database, Redis, and server settings.
-   **`utils/`**: Provides various utility functions, such as JSON encoding/decoding, encryption, and validation helpers.

## Features

- OAuth2 Client Registration Endpoint (`/register`)
- OAuth2 Authorization Endpoint (`/authorize`)
- OAuth2 Token Endpoint (`/token`)
- OAuth2 Revocation Endpoint (`/revoke`)
- OAuth2 Introspection Endpoint (`/introspect`)
- OAuth2 Userinfo Endpoint (`/userinfo`)
- OAuth2 JWKS Endpoint (`/.well-known/jwks.json`)
- Built with Go
- PostgreSQL for data storage
- Redis for session storage
- Google as Identity Provider (IDP) integration for authentication.
- Docker Compose for easy setup and development environment.

## Technology Stack

-   **Go**: The primary programming language, chosen for its performance, concurrency features (Goroutines), and strong ecosystem for building reliable network services.
-   **PostgreSQL**: Robust relational database for storing client information, scopes, authorization codes, and tokens.
-   **Redis**: Used for high-performance session management and potentially for caching frequently accessed data.
-   **Zap (Uber's Zap)**: Structured logging library for efficient and context-rich application logging.
-   **GORM**: An ORM (Object-Relational Mapper) for Go, simplifying database interactions with PostgreSQL.
-   **Docker & Docker Compose**: For containerization, providing a consistent and isolated development/production environment.

## Concurrency Model

The Go OAuth2 Server leverages **Goroutines** and **Channels** for efficient concurrency:

-   **HTTP Request Handling**: The `net/http` server automatically handles each incoming HTTP request in its own Goroutine, ensuring responsiveness for concurrent client connections.
-   **Background Tasks**: Non-blocking operations, such as sending audit logs, metrics updates, or other post-response processing, can be offloaded to separate Goroutines to avoid delaying the client's response.
-   **Parallel Processing**: Potentially long-running or independent sub-operations within a request (e.g., querying multiple external services) can be executed in parallel using Goroutines.

This approach maximizes resource utilization and ensures the server remains responsive even under heavy load.

## Security Considerations

Building an OAuth2 server requires careful attention to security. Key considerations and practices in this project include:

-   **Client Secret Encryption**: Storing client secrets securely (e.g., using bcrypt hashing) rather than plaintext.
-   **Redirect URI Validation**: Strict validation of `redirect_uri` to prevent open redirect vulnerabilities.
-   **State Parameter**: Encouraging the use of the `state` parameter in authorization requests to mitigate CSRF attacks.
-   **Token Expiration & Revocation**: Implementing proper access token expiration, refresh token mechanisms, and revocation endpoints.
-   **Secure Communication**: Assumed use of HTTPS for all communication in production environments.
-   **Input Validation**: Comprehensive validation of all incoming request parameters to prevent injection attacks and ensure compliance with OAuth2 specifications.

## Extensibility

The modular design of this project aims for high extensibility:

-   **New Grant Types/Response Types**: Easily add support for new OAuth2 grant types or response types by implementing the respective interfaces or logic within the `oauth/` and `services/` packages.
-   **Authentication Methods**: Integrate additional client authentication methods (`TokenEndpointAuthMethod`) by extending the `oauth/authmethodtype` package and corresponding validation logic.
-   **Identity Providers**: The architecture allows for flexible integration of different Identity Providers (beyond Google) by abstracting the IDP interaction logic within the authentication flow.
-   **Database/Storage**: The `store/repositories` interfaces allow for swapping out the underlying database or storage mechanism with minimal impact on the service layer.

## Prerequisites

- Docker and Docker Compose installed on your machine.
- Go (version 1.20 or later) installed locally if you plan to run without Docker.

## Getting Started

### Clone the Repository

```bash
git clone https://github.com/manuelrojas19/go-oauth2-server
cd go-oauth2-server
```

### Build and Run with Docker Compose (Recommended for Development)

This command builds the Go application Docker image, sets up PostgreSQL and Redis, and starts all services.

```bash
docker-compose up --build -d
```

### Run the Go Application Locally (without Docker for the Go app itself)

First, ensure your PostgreSQL and Redis instances are running (e.g., via `docker-compose up -d postgres redis`). Then:

```bash
go run cmd/server/main.go
```

## API Endpoints (Implemented & Documented)

### Client Registration: `/register`

- **Description**: Registers a new OAuth2 client. Returns client credentials and metadata.
- **Method**: `POST`
- **Request**:
    - **Headers**:
        - `Content-Type: application/json`
    - **Body**:
      ```json
      {
        "client_name": "MyAwesomeClient",
        "grant_types": ["authorization_code", "client_credentials", "refresh_token"],
        "response_types": ["code"],
        "token_endpoint_auth_method": "client_secret_basic",
        "redirect_uris": ["http://localhost:8080/callback", "https://other.example.com/callback"],
        "scope": "openid profile email offline_access"
      }
      ```
- **Response**:
    - **Status**: `201 Created` (on successful registration)
    - **Body**:
      ```json
      {
        "client_id": "string-uuid",
        "client_secret": "string-uuid",
        "client_id_issued_at": "unix-timestamp-string",
        "client_secret_expires_at": "unix-timestamp-string",
        "client_name": "MyAwesomeClient",
        "grant_types": ["authorization_code", "client_credentials", "refresh_token"],
        "response_types": ["code"],
        "token_endpoint_auth_method": "client_secret_basic",
        "redirect_uris": ["http://localhost:8080/callback", "https://other.example.com/callback"],
        "scopes": [{"Name": "openid", "Description": "OpenID Connect scope"}, {"Name": "profile", "Description": "User profile information"}]
      }
      ```
- **Error Responses**: (Example for invalid redirect URI)
    - **Status**: `400 Bad Request`
    - **Body**:
      ```json
      {
        "error": "invalid_request",
        "error_description": "malformed redirect_uri: http://invalid"
      }
      ```

### Authorization Endpoint: `/authorize`

- **Description**: Initiates an OAuth2 authorization flow. Users will be redirected to a login/consent page if necessary.
- **Method**: `GET`
- **Request**:
    - **Query Parameters**:
        - `response_type`: (`code` or `token`) Indicates the type of response desired.
        - `client_id`: (string) The client identifier as obtained during registration.
        - `redirect_uri`: (string) The registered redirection URI.
        - `scope`: (string, space-separated) The desired access token scopes (e.g., `openid profile email`).
        - `state`: (string, optional) An opaque value used to maintain state between the request and the callback.
    - **Example Query**:
      ```
      GET /authorize?response_type=code&client_id=YOUR_CLIENT_ID&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback&scope=openid%20profile%20email&state=xyz
      ```
- **Response**:
    - **Status**: `302 Found` (on successful authorization)
    - **Headers**:
        - `Location: YOUR_REDIRECT_URI?code=AUTHORIZATION_CODE&state=xyz`

### Token Endpoint: `/token`

- **Description**: Exchanges an authorization code for an access token, refreshes tokens, or handles client credentials grant.
- **Method**: `POST`
- **Request**:
    - **Headers**:
        - `Content-Type: application/x-www-form-urlencoded`
    - **Body (Form Data)**:
        - **Authorization Code Grant Type**:
          ```
          grant_type=authorization_code&code=YOUR_AUTH_CODE&redirect_uri=YOUR_REDIRECT_URI&client_id=YOUR_CLIENT_ID&client_secret=YOUR_CLIENT_SECRET
          ```
        - **Client Credentials Grant Type**:
          ```
          grant_type=client_credentials&client_id=YOUR_CLIENT_ID&client_secret=YOUR_CLIENT_SECRET
          ```
        - **Refresh Token Grant Type**:
          ```
          grant_type=refresh_token&refresh_token=YOUR_REFRESH_TOKEN&client_id=YOUR_CLIENT_ID&client_secret=YOUR_CLIENT_SECRET
          ```
- **Response**:
    - **Status**: `200 OK` (on successful token issuance)
    - **Body**:
      ```json
      {
        "access_token": "jwt-string",
        "token_type": "bearer",
        "expires_in": 3600, // seconds until expiration
        "refresh_token": "jwt-string",
        "scope": "openid profile email"
      }
      ```

## Example `curl` Commands

### Client Registration

To register a new OAuth2 client, use the following `curl` command:

```bash
curl -X POST http://localhost:8080/register \
    -H "Content-Type: application/json" \
    -d '{
        "client_name": "MyAwesomeClient",
        "grant_types": ["authorization_code", "client_credentials", "refresh_token"],
        "response_types": ["code"],
        "token_endpoint_auth_method": "client_secret_basic",
        "redirect_uris": ["http://localhost:8080/callback"],
        "scope": "openid profile email"
    }'
```

### Authorization Code Flow (Manual Steps)

1.  **Register a Client** (using the command above) and note the `client_id`.

2.  **Initiate Authorization** (User is redirected to login/consent):

    Open in your browser:
    `http://localhost:8080/authorize?response_type=code&client_id=YOUR_CLIENT_ID&redirect_uri=http://localhost:8080/callback&scope=openid%20profile%20email&state=random_state_string`

    (Replace `YOUR_CLIENT_ID` and `random_state_string`)

    *You will be prompted to log in and grant consent.* After successful authorization, your browser will be redirected to `http://localhost:8080/callback?code=YOUR_AUTH_CODE&state=random_state_string`.
    **Note down the `YOUR_AUTH_CODE` from the URL.**

3.  **Exchange Authorization Code for Access Token**:

    Use the `curl` command below with the `YOUR_AUTH_CODE` and your `CLIENT_SECRET` (from step 1):

    ```bash
    curl -X POST http://localhost:8080/token \
        -H "Content-Type: application/x-www-form-urlencoded" \
        -d "grant_type=authorization_code" \
        -d "code=YOUR_AUTH_CODE" \
        -d "client_id=YOUR_CLIENT_ID" \
        -d "client_secret=YOUR_CLIENT_SECRET" \
        -d "redirect_uri=http://localhost:8080/callback"
    ```

    This will return an `access_token` and potentially a `refresh_token`.

### Client Credentials Grant Flow

1.  **Register a Client** (using the client registration curl command) ensuring `client_credentials` is in `grant_types`. Note the `client_id` and `client_secret`.

2.  **Request Access Token**:

    ```bash
    curl -X POST http://localhost:8080/token \
        -H "Content-Type: application/x-www-form-urlencoded" \
        -d "grant_type=client_credentials" \
        -d "client_id=YOUR_CLIENT_ID" \
        -d "client_secret=YOUR_CLIENT_SECRET"
    ```

### Refresh Token Grant Flow

1.  **Perform Authorization Code Flow** (steps 1-3 above) to obtain a `refresh_token`.

2.  **Request New Access Token using Refresh Token**:

    ```bash
    curl -X POST http://localhost:8080/token \
        -H "Content-Type: application/x-www-form-urlencoded" \
        -d "grant_type=refresh_token" \
        -d "refresh_token=YOUR_REFRESH_TOKEN" \
        -d "client_id=YOUR_CLIENT_ID" \
        -d "client_secret=YOUR_CLIENT_SECRET"
    ```

## Development

### Running Tests

```bash
go test ./...
```

### Database Migrations

This project uses GORM for database interactions. Migrations are typically handled programmatically on application startup.

### Environment Variables

Configuration is managed through environment variables. Key variables include:

- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`: PostgreSQL connection details.
- `REDIS_ADDR`, `REDIS_PASSWORD`: Redis connection details.
- `JWT_SECRET`: Secret key for signing JWTs.
- `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_REDIRECT_URI`: Credentials for Google IDP integration.
- `SERVER_PORT`: The port on which the OAuth2 server listens.

## Contributing

Feel free to fork the repository and contribute! Please open issues for bugs or feature requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.