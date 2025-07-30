# Manual cURL Commands for E2E Testing

This document provides individual cURL commands for manual testing of the Venturo Golang Core API.

## Prerequisites

1. Start the server: `docker-compose up --build`
2. Run migrations: `docker-compose run --rm app go run ./cmd/migrate/main.go up`
3. Server should be running on `http://localhost:3000`

## 1. Health Check

```bash
curl http://localhost:3000/health
```

## 2. User Registration

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "securepassword123"
  }' \
  http://localhost:3000/api/v1/register
```

## 3. User Login

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "securepassword123"
  }' \
  http://localhost:3000/api/v1/login
```

**Save the tokens from the response:**
- `access_token`: Use for protected endpoints
- `refresh_token`: Use for token refresh

## 4. Get User Profile (Protected)

```bash
curl -X GET \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  http://localhost:3000/api/v1/profile
```

## 5. Update User Profile (Protected)

```bash
curl -X PUT \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Test User"
  }' \
  http://localhost:3000/api/v1/profile
```

## 6. Token Refresh

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }' \
  http://localhost:3000/api/v1/refresh
```

## 7. Create Post (Protected)

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "title": "My Test Post",
    "body": "This is the content of my test post."
  }' \
  http://localhost:3000/api/v1/posts
```

**Save the `id` from the response for subsequent operations.**

## 8. Get All Posts (Public)

```bash
curl -X GET \
  "http://localhost:3000/api/v1/posts?page=1&limit=10"
```

## 9. Get Post by ID (Public)

```bash
curl -X GET \
  http://localhost:3000/api/v1/posts/YOUR_POST_ID
```

## 10. Update Post (Protected)

```bash
curl -X PUT \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "title": "Updated Test Post",
    "body": "This is the updated content of my test post."
  }' \
  http://localhost:3000/api/v1/posts/YOUR_POST_ID
```

## 11. Delete Post (Protected)

```bash
curl -X DELETE \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  http://localhost:3000/api/v1/posts/YOUR_POST_ID
```

## 12. Test Rate Limiting

Run this command multiple times rapidly (>10 times) to trigger rate limiting:

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "email": "nonexistent@example.com",
    "password": "wrongpassword"
  }' \
  http://localhost:3000/api/v1/login
```

You should eventually receive a `429 Too Many Requests` response.

## 13. Test Invalid Token

```bash
curl -X GET \
  -H "Authorization: Bearer invalid_token_12345" \
  http://localhost:3000/api/v1/profile
```

Should return `401 Unauthorized`.

## Example Complete Flow

Here's a complete example with actual commands (replace tokens as you get them):

```bash
# 1. Register
curl -X POST -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com", "password": "password123"}' \
  http://localhost:3000/api/v1/register

# 2. Login (save the tokens)
curl -X POST -H "Content-Type: application/json" \
  -d '{"email": "john@example.com", "password": "password123"}' \
  http://localhost:3000/api/v1/login

# 3. Create post (use access_token from login)
curl -X POST -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{"title": "Hello World", "body": "My first post"}' \
  http://localhost:3000/api/v1/posts

# 4. Get all posts
curl http://localhost:3000/api/v1/posts?page=1&limit=5

# 5. Update post (use post ID from create response)
curl -X PUT -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{"title": "Updated Hello World", "body": "My updated first post"}' \
  http://localhost:3000/api/v1/posts/12345678-1234-1234-1234-123456789012

# 6. Delete post
curl -X DELETE \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  http://localhost:3000/api/v1/posts/12345678-1234-1234-1234-123456789012
```

## Swagger Documentation

For interactive testing, visit the Swagger UI at:
```
http://localhost:3000/swagger/index.html
```

## Notes

- Replace `YOUR_ACCESS_TOKEN`, `YOUR_REFRESH_TOKEN`, and `YOUR_POST_ID` with actual values from responses
- Tokens expire, so you may need to refresh or re-login during testing
- The access token typically expires in 15 minutes, refresh token in 24 hours
- Rate limiting is set to 10 requests per minute for auth endpoints
- All protected endpoints require the `Authorization: Bearer <token>` header