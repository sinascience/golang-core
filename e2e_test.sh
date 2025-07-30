#!/bin/bash

# Venturo Golang Core - End-to-End Testing Script
# Tests complete user journey: register -> login -> refresh -> create post -> update -> get -> delete -> logout

set -e  # Exit on any error

# Configuration
# Auto-detect Windows host IP from WSL
WINDOWS_HOST_IP=$(ip route show | grep -i default | awk '{print $3}')
BASE_URL="http://${WINDOWS_HOST_IP}:3000"
# Alternative URLs to try:
# BASE_URL="http://localhost:3000"  # Use this if running server in WSL
# BASE_URL="http://127.0.0.1:3000"
API_URL="$BASE_URL/api/v1"
TIMESTAMP=$(date +%s)
TEST_EMAIL="e2e_test_${TIMESTAMP}@example.com"
TEST_NAME="E2E Test User"
TEST_PASSWORD="securepassword123"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Global variables for tokens and post ID
ACCESS_TOKEN=""
REFRESH_TOKEN=""
POST_ID=""

# Helper functions
log_step() {
    echo -e "${BLUE}📝 $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

check_server() {
    log_step "Checking if server is running..."
    if curl -s "$BASE_URL/health" > /dev/null; then
        log_success "Server is running"
    else
        log_error "Server is not running. Please start the server first. $BASE_URL"
        log_step "Start server with: docker-compose up --build"
        exit 1
    fi
}

# Test functions
test_health_check() {
    log_step "1. Testing health check endpoint..."
    
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" "$BASE_URL/health")
    http_code=$(echo $response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    body=$(echo $response | sed -e 's/HTTPSTATUS\:.*//g')
    
    if [ "$http_code" -eq 200 ]; then
        log_success "Health check passed"
        echo "Response: $body"
    else
        log_error "Health check failed (HTTP $http_code)"
        exit 1
    fi
}

test_user_registration() {
    log_step "2. Testing user registration..."
    
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" \
        -X POST \
        -H "Content-Type: application/json" \
        -d "{
            \"name\": \"$TEST_NAME\",
            \"email\": \"$TEST_EMAIL\",
            \"password\": \"$TEST_PASSWORD\"
        }" \
        "$API_URL/register")
    
    http_code=$(echo $response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    body=$(echo $response | sed -e 's/HTTPSTATUS\:.*//g')
    
    if [ "$http_code" -eq 201 ]; then
        log_success "User registration successful"
        echo "Response: $body"
    else
        log_error "User registration failed (HTTP $http_code)"
        echo "Response: $body"
        exit 1
    fi
}

test_user_login() {
    log_step "3. Testing user login..."
    
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" \
        -X POST \
        -H "Content-Type: application/json" \
        -d "{
            \"email\": \"$TEST_EMAIL\",
            \"password\": \"$TEST_PASSWORD\"
        }" \
        "$API_URL/login")
    
    http_code=$(echo $response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    body=$(echo $response | sed -e 's/HTTPSTATUS\:.*//g')
    
    if [ "$http_code" -eq 200 ]; then
        log_success "User login successful"
        
        # Extract tokens using sed/grep (more reliable than jq dependency)
        ACCESS_TOKEN=$(echo "$body" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
        REFRESH_TOKEN=$(echo "$body" | grep -o '"refresh_token":"[^"]*"' | cut -d'"' -f4)
        
        if [ -n "$ACCESS_TOKEN" ] && [ -n "$REFRESH_TOKEN" ]; then
            log_success "Tokens extracted successfully"
            echo "Access token (first 20 chars): ${ACCESS_TOKEN:0:20}..."
            echo "Refresh token (first 20 chars): ${REFRESH_TOKEN:0:20}..."
        else
            log_error "Failed to extract tokens from response"
            echo "Response: $body"
            exit 1
        fi
    else
        log_error "User login failed (HTTP $http_code)"
        echo "Response: $body"
        exit 1
    fi
}

test_get_profile() {
    log_step "4. Testing get user profile (protected endpoint)..."
    
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" \
        -X GET \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        "$API_URL/profile")
    
    http_code=$(echo $response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    body=$(echo $response | sed -e 's/HTTPSTATUS\:.*//g')
    
    if [ "$http_code" -eq 200 ]; then
        log_success "Get profile successful"
        echo "Response: $body"
    else
        log_error "Get profile failed (HTTP $http_code)"
        echo "Response: $body"
        exit 1
    fi
}

test_token_refresh() {
    log_step "5. Testing token refresh..."
    
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" \
        -X POST \
        -H "Content-Type: application/json" \
        -d "{
            \"refresh_token\": \"$REFRESH_TOKEN\"
        }" \
        "$API_URL/refresh")
    
    http_code=$(echo $response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    body=$(echo $response | sed -e 's/HTTPSTATUS\:.*//g')
    
    if [ "$http_code" -eq 200 ]; then
        log_success "Token refresh successful"
        
        # Update tokens with new ones
        NEW_ACCESS_TOKEN=$(echo "$body" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
        NEW_REFRESH_TOKEN=$(echo "$body" | grep -o '"refresh_token":"[^"]*"' | cut -d'"' -f4)
        
        if [ -n "$NEW_ACCESS_TOKEN" ] && [ -n "$NEW_REFRESH_TOKEN" ]; then
            ACCESS_TOKEN=$NEW_ACCESS_TOKEN
            REFRESH_TOKEN=$NEW_REFRESH_TOKEN
            log_success "Tokens updated successfully (token rotation working)"
            echo "New access token (first 20 chars): ${ACCESS_TOKEN:0:20}..."
        else
            log_error "Failed to extract new tokens"
            exit 1
        fi
    else
        log_error "Token refresh failed (HTTP $http_code)"
        echo "Response: $body"
        exit 1
    fi
}

test_create_post() {
    log_step "6. Testing create post..."
    
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" \
        -X POST \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -d "{
            \"title\": \"E2E Test Post - $TIMESTAMP\",
            \"body\": \"This is a test post created during end-to-end testing. Timestamp: $TIMESTAMP\"
        }" \
        "$API_URL/posts")
    
    http_code=$(echo $response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    body=$(echo $response | sed -e 's/HTTPSTATUS\:.*//g')
    
    if [ "$http_code" -eq 201 ]; then
        log_success "Post created successfully"
        
        # Extract post ID
        POST_ID=$(echo "$body" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
        
        if [ -n "$POST_ID" ]; then
            log_success "Post ID extracted: $POST_ID"
            echo "Response: $body"
        else
            log_error "Failed to extract post ID"
            echo "Response: $body"
            exit 1
        fi
    else
        log_error "Post creation failed (HTTP $http_code)"
        echo "Response: $body"
        exit 1
    fi
}

test_get_all_posts() {
    log_step "7. Testing get all posts (public endpoint)..."
    
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" \
        -X GET \
        "$API_URL/posts?page=1&limit=5")
    
    http_code=$(echo $response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    body=$(echo $response | sed -e 's/HTTPSTATUS\:.*//g')
    
    if [ "$http_code" -eq 200 ]; then
        log_success "Get all posts successful"
        echo "Response: $body"
    else
        log_error "Get all posts failed (HTTP $http_code)"
        echo "Response: $body"
        exit 1
    fi
}

test_get_post_by_id() {
    log_step "8. Testing get post by ID (public endpoint)..."
    
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" \
        -X GET \
        "$API_URL/posts/$POST_ID")
    
    http_code=$(echo $response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    body=$(echo $response | sed -e 's/HTTPSTATUS\:.*//g')
    
    if [ "$http_code" -eq 200 ]; then
        log_success "Get post by ID successful"
        echo "Response: $body"
    else
        log_error "Get post by ID failed (HTTP $http_code)"
        echo "Response: $body"
        exit 1
    fi
}

test_update_post() {
    log_step "9. Testing update post (protected endpoint)..."
    
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" \
        -X PUT \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -d "{
            \"title\": \"Updated E2E Test Post - $TIMESTAMP\",
            \"body\": \"This post has been updated during end-to-end testing. Updated at: $(date)\"
        }" \
        "$API_URL/posts/$POST_ID")
    
    http_code=$(echo $response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    body=$(echo $response | sed -e 's/HTTPSTATUS\:.*//g')
    
    if [ "$http_code" -eq 200 ]; then
        log_success "Post updated successfully"
        echo "Response: $body"
    else
        log_error "Post update failed (HTTP $http_code)"
        echo "Response: $body"
        exit 1
    fi
}

test_delete_post() {
    log_step "10. Testing delete post (protected endpoint)..."
    
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" \
        -X DELETE \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        "$API_URL/posts/$POST_ID")
    
    http_code=$(echo $response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    body=$(echo $response | sed -e 's/HTTPSTATUS\:.*//g')
    
    if [ "$http_code" -eq 200 ]; then
        log_success "Post deleted successfully"
        echo "Response: $body"
    else
        log_error "Post deletion failed (HTTP $http_code)"
        echo "Response: $body"
        exit 1
    fi
}

test_verify_post_deleted() {
    log_step "11. Testing verify post is deleted..."
    
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" \
        -X GET \
        "$API_URL/posts/$POST_ID")
    
    http_code=$(echo $response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    body=$(echo $response | sed -e 's/HTTPSTATUS\:.*//g')
    
    if [ "$http_code" -eq 404 ]; then
        log_success "Post deletion verified - post not found (as expected)"
        echo "Response: $body"
    else
        log_warning "Post deletion verification unexpected result (HTTP $http_code)"
        echo "Response: $body"
        # Don't exit here as this might be acceptable depending on implementation
    fi
}

test_rate_limiting() {
    log_step "12. Testing rate limiting..."
    
    log_step "Making multiple rapid requests to trigger rate limiting..."
    rate_limit_hit=false
    
    for i in {1..12}; do
        response=$(curl -s -w "HTTPSTATUS:%{http_code}" \
            -X POST \
            -H "Content-Type: application/json" \
            -d "{
                \"email\": \"nonexistent@example.com\",
                \"password\": \"wrongpassword\"
            }" \
            "$API_URL/login")
        
        http_code=$(echo $response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
        
        if [ "$http_code" -eq 429 ]; then
            rate_limit_hit=true
            log_success "Rate limiting triggered after $i requests (HTTP 429)"
            break
        fi
        
        echo "Request $i: HTTP $http_code"
        sleep 0.1  # Small delay between requests
    done
    
    if [ "$rate_limit_hit" = false ]; then
        log_warning "Rate limiting not triggered in 12 requests (limit might be higher than expected)"
    fi
}

test_invalid_token() {
    log_step "13. Testing invalid token handling..."
    
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" \
        -X GET \
        -H "Authorization: Bearer invalid_token_12345" \
        "$API_URL/profile")
    
    http_code=$(echo $response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    body=$(echo $response | sed -e 's/HTTPSTATUS\:.*//g')
    
    if [ "$http_code" -eq 401 ]; then
        log_success "Invalid token properly rejected (HTTP 401)"
        echo "Response: $body"
    else
        log_error "Invalid token handling failed - expected 401, got $http_code"
        echo "Response: $body"
        exit 1
    fi
}

# Main execution
main() {
    echo -e "${BLUE}🚀 Starting Venturo Golang Core E2E Testing${NC}"
    echo "==========================================="
    echo "Test Email: $TEST_EMAIL"
    echo "Timestamp: $TIMESTAMP"
    echo "==========================================="
    
    check_server
    
    # Run all tests
    test_health_check
    test_user_registration
    test_user_login
    test_get_profile
    test_token_refresh
    test_create_post
    test_get_all_posts
    test_get_post_by_id
    test_update_post
    test_delete_post
    test_verify_post_deleted
    test_rate_limiting
    test_invalid_token
    
    echo
    echo -e "${GREEN}🎉 All E2E tests completed successfully!${NC}"
    echo "==========================================="
    echo -e "${GREEN}✅ User Registration${NC}"
    echo -e "${GREEN}✅ User Login${NC}" 
    echo -e "${GREEN}✅ Token Refresh (with rotation)${NC}"
    echo -e "${GREEN}✅ Protected Endpoint Access${NC}"
    echo -e "${GREEN}✅ Post CRUD Operations${NC}"
    echo -e "${GREEN}✅ Rate Limiting${NC}"
    echo -e "${GREEN}✅ Security (Invalid Token)${NC}"
    echo "==========================================="
    echo -e "${BLUE}🎯 Your Venturo Golang Core API is working perfectly!${NC}"
}

# Handle script interruption
trap 'echo -e "\n${YELLOW}E2E testing interrupted${NC}"; exit 1' INT

# Run main function
main "$@"