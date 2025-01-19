# Chirpy - A CLI App Mimicking Twitter

## Overview

Chirpy is a simple command-line interface (CLI) application designed to mimic the core functionalities of Twitter. Users can create, read, update, and delete short messages (chirps), manage user accounts, and handle authentication securely.

## Why Chirpy?

Why Chirpy, you ask? Well, every aspiring developer needs a pet project, and Chirpy is the ultimate exercise in building an API that _looks_ like it means business. Want to flex those API muscles? Dream of crafting database migrations that no one will ever migrate? Chirpy is here for you. It’s the perfect playground to hone your skills in creating a full-fledged, feature-rich CLI tool that mimics the big leagues—all while knowing that the only user might be... you. But hey, practice makes perfect, right? And what better way to practice than by building a tiny Twitter clone that lives and dies in your local environment? Chirpy: because why not?In the fast-paced world of social media, Chirpy offers a minimalist approach to communication. It allows users to interact with a Twitter-like platform directly from their terminal, providing a unique, distraction-free environment for sharing thoughts. Ideal for developers, CLI enthusiasts, or anyone looking for a lightweight alternative to traditional social media platforms.

## Installation and Setup

### Prerequisites

- Go (version 1.16 or higher)
- PostgreSQL (for the backend database)
- GOPATH properly configured

### Installation

1. Clone the Chirpy repository:

   ```bash
   git clone https://github.com/juaniten/chirpy.git
   cd chirpy
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Configure the environment variables:
   Create a `.env` file in the project root with the following variables:

   ```env
   DATABASE_URL=postgres://<username>:<password>@localhost:5432/chirpy_db
   JWT_SECRET=<your-jwt-secret>
   POLKA_KEY=<your-polka-key>
   ```

4. Set up the database:

   ```bash
   go run scripts/migrate.go
   ```

5. Run the application:
   ```bash
   go run main.go
   ```

## API Documentation

### Authentication

#### `POST /login`

Authenticates a user and returns a JWT token and refresh token.

**Request:**

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**

```json
{
  "id": "<uuid>",
  "created_at": "<timestamp>",
  "updated_at": "<timestamp>",
  "email": "user@example.com",
  "is_chirpy_red": false,
  "token": "<jwt-token>",
  "refresh_token": "<refresh-token>"
}
```

#### `POST /refresh-token`

Generates a new access token using a refresh token.

**Request Header:**

```
Authorization: Bearer <refresh-token>
```

**Response:**

```json
{
  "token": "<new-jwt-token>"
}
```

#### `POST /revoke-token`

Revokes a given refresh token.

**Request Header:**

```
Authorization: Bearer <refresh-token>
```

**Response:**

- HTTP 204 No Content

### User Management

#### `POST /users`

Creates a new user.

**Request:**

```json
{
  "email": "newuser@example.com",
  "password": "securepassword"
}
```

**Response:**

```json
{
  "id": "<uuid>",
  "created_at": "<timestamp>",
  "updated_at": "<timestamp>",
  "email": "newuser@example.com",
  "is_chirpy_red": false
}
```

#### `PUT /users`

Updates an existing user's email and/or password.

**Request Header:**

```
Authorization: Bearer <access-token>
```

**Request Body:**

```json
{
  "email": "updateduser@example.com",
  "password": "newpassword"
}
```

**Response:**

```json
{
  "id": "<uuid>",
  "created_at": "<timestamp>",
  "updated_at": "<timestamp>",
  "email": "updateduser@example.com",
  "is_chirpy_red": false
}
```

### Chirp Management

#### `POST /chirps`

Creates a new chirp.

**Request Header:**

```
Authorization: Bearer <access-token>
```

**Request Body:**

```json
{
  "body": "This is a new chirp!"
}
```

**Response:**

```json
{
  "id": "<uuid>",
  "created_at": "<timestamp>",
  "updated_at": "<timestamp>",
  "body": "This is a new chirp!",
  "user_id": "<uuid>"
}
```

#### `GET /chirps`

Fetches all chirps or chirps from a specific author.

**Query Parameters:**

- `author_id`: Filter by author (optional)
- `sort`: Sort order (`asc` or `desc`)

**Response:**

```json
[
  {
    "id": "<uuid>",
    "created_at": "<timestamp>",
    "updated_at": "<timestamp>",
    "body": "First chirp!",
    "user_id": "<uuid>"
  },
  {
    "id": "<uuid>",
    "created_at": "<timestamp>",
    "updated_at": "<timestamp>",
    "body": "Second chirp!",
    "user_id": "<uuid>"
  }
]
```

#### `GET /chirps/:chirpID`

Fetches a specific chirp by its ID.

**Response:**

```json
{
  "id": "<uuid>",
  "created_at": "<timestamp>",
  "updated_at": "<timestamp>",
  "body": "This is a chirp",
  "user_id": "<uuid>"
}
```

#### `DELETE /chirps/:chirpID`

Deletes a specific chirp by its ID.

**Request Header:**

```
Authorization: Bearer <access-token>
```

**Response:**

- HTTP 204 No Content

### Webhook

#### `POST /polka-webhook`

Handles incoming webhook events from the Polka service.

**Request Header:**

```
X-API-Key: <polka-key>
```

**Request Body:**

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "<uuid>"
  }
}
```

**Response:**

- HTTP 204 No Content
