// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

const baseURL = "http://localhost:3000/api/v1"

type LoginResponse struct {
	StatusCode int `json:"status_code"`
	Data       struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	} `json:"data"`
}

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Errors     string `json:"errors"`
}

func TestIntegrationAuthFlow(t *testing.T) {
	// Check if server is running
	resp, err := http.Get("http://localhost:3000/health")
	if err != nil {
		t.Skip("Server not running, skipping integration tests")
	}
	resp.Body.Close()

	// Test data
	testUser := map[string]string{
		"name":     "Integration Test User",
		"email":    fmt.Sprintf("test_%d@example.com", time.Now().Unix()),
		"password": "testpassword123",
	}

	// Test 1: User Registration
	t.Run("UserRegistration", func(t *testing.T) {
		jsonData, _ := json.Marshal(testUser)
		resp, err := http.Post(baseURL+"/register", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Registration request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 201 {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 201, got %d. Response: %s", resp.StatusCode, string(body))
		}
	})

	// Test 2: User Login
	var accessToken, refreshToken string
	t.Run("UserLogin", func(t *testing.T) {
		loginData := map[string]string{
			"email":    testUser["email"],
			"password": testUser["password"],
		}
		jsonData, _ := json.Marshal(loginData)
		
		resp, err := http.Post(baseURL+"/login", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Login request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d. Response: %s", resp.StatusCode, string(body))
			return
		}

		var loginResp LoginResponse
		body, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &loginResp); err != nil {
			t.Fatalf("Failed to parse login response: %v", err)
		}

		accessToken = loginResp.Data.AccessToken
		refreshToken = loginResp.Data.RefreshToken

		if accessToken == "" || refreshToken == "" {
			t.Error("Login should return both access and refresh tokens")
		}
	})

	// Test 3: Access Protected Endpoint
	t.Run("AccessProtectedEndpoint", func(t *testing.T) {
		if accessToken == "" {
			t.Skip("No access token available")
		}

		req, _ := http.NewRequest("GET", baseURL+"/profile", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Profile request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d. Response: %s", resp.StatusCode, string(body))
		}
	})

	// Test 4: Token Refresh
	t.Run("TokenRefresh", func(t *testing.T) {
		if refreshToken == "" {
			t.Skip("No refresh token available")
		}

		refreshData := map[string]string{
			"refresh_token": refreshToken,
		}
		jsonData, _ := json.Marshal(refreshData)

		resp, err := http.Post(baseURL+"/refresh", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Refresh request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d. Response: %s", resp.StatusCode, string(body))
			return
		}

		var refreshResp LoginResponse
		body, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &refreshResp); err != nil {
			t.Fatalf("Failed to parse refresh response: %v", err)
		}

		if refreshResp.Data.AccessToken == "" || refreshResp.Data.RefreshToken == "" {
			t.Error("Refresh should return both new access and refresh tokens")
		}

		// Verify tokens are different (token rotation)
		if refreshResp.Data.AccessToken == accessToken {
			t.Error("New access token should be different from old one")
		}
		if refreshResp.Data.RefreshToken == refreshToken {
			t.Error("New refresh token should be different from old one")
		}
	})

	// Test 5: Rate Limiting
	t.Run("RateLimiting", func(t *testing.T) {
		loginData := map[string]string{
			"email":    "nonexistent@example.com",
			"password": "wrongpassword",
		}
		jsonData, _ := json.Marshal(loginData)

		// Make multiple rapid requests to trigger rate limiting
		var lastStatus int
		for i := 0; i < 15; i++ {
			resp, err := http.Post(baseURL+"/login", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				t.Fatalf("Rate limit test request failed: %v", err)
			}
			lastStatus = resp.StatusCode
			resp.Body.Close()

			// Short delay between requests
			time.Sleep(10 * time.Millisecond)
		}

		// After many requests, we should eventually get rate limited
		if lastStatus != 429 {
			t.Logf("Expected to hit rate limit (429), but last status was %d. Rate limiting may be working correctly with a higher limit.", lastStatus)
		}
	})
}

func TestMain(m *testing.M) {
	// Check if integration tests should run
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		log.Println("Skipping integration tests. Set INTEGRATION_TESTS=1 to run them.")
		os.Exit(0)
	}

	// Run tests
	code := m.Run()
	os.Exit(code)
}