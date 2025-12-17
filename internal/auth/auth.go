package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-rod/rod/lib/proto"
	
	"subspace/internal/browser"
	"subspace/internal/config"
	"subspace/internal/logger"
	"subspace/internal/stealth"
	"subspace/internal/storage"
)

/*
AUTHENTICATION MODULE - EDUCATIONAL IMPLEMENTATION

This module demonstrates session management and login flow concepts.
It does NOT contain real selectors or working LinkedIn automation.

FEATURES:
- Session cookie persistence and reuse
- Mock login flow showing stealth integration
- Checkpoint detection simulation
- Exponential backoff on retries

EDUCATIONAL NOTE:
In a real system, this would:
1. Use actual form selectors
2. Handle CAPTCHA challenges
3. Implement 2FA flows
4. Monitor session validity
*/

// Authenticator handles login and session management
type Authenticator struct {
	browser browser.Controller
	stealth *stealth.Stealth
	storage *storage.Storage
	config  config.AuthConfig
	log     *logger.ContextLogger
}

// New creates a new authenticator
func New(b browser.Controller, s *stealth.Stealth, storage *storage.Storage) *Authenticator {
	// Load auth config from environment
	cfg := config.AuthConfig{
		SessionCookiePath: config.GetEnv("SESSION_COOKIE_PATH", "./data/session.json"),
		ReuseSession:      true,
		CheckpointRetries: 3,
	}

	return &Authenticator{
		browser: b,
		stealth: s,
		storage: storage,
		config:  cfg,
		log:     logger.NewContext("auth"),
	}
}

// Login performs the login flow with session reuse and stealth
func (a *Authenticator) Login() error {
	a.log.Info("Starting authentication process")
	start := time.Now()

	// Step 1: Try to reuse existing session
	if a.config.ReuseSession {
		if err := a.tryLoadSession(); err == nil {
			a.log.Info("Session restored from cookies")
			logger.Timing("auth", "login", start, nil)
			return nil
		}
		a.log.Info("No valid session found, proceeding with login")
	}

	// Step 2: Perform login with retry logic
	var lastErr error
	for attempt := 1; attempt <= a.config.CheckpointRetries; attempt++ {
		a.log.Info("Login attempt", "attempt", attempt, "max", a.config.CheckpointRetries)
		
		// Record attempt in storage
		a.storage.LogAction("login_attempt", "", false, nil)

		err := a.performLogin()
		if err == nil {
			// Success! Save session
			if err := a.saveSession(); err != nil {
				a.log.Warn("Failed to save session", "error", err)
			}
			
			// Record successful login in storage
			a.storage.LogAction("login_success", "", true, nil)
			
			logger.Timing("auth", "login", start, nil)
			return nil
		}

		lastErr = err
		a.log.Warn("Login attempt failed", "attempt", attempt, "error", err)

		// Check if it's a checkpoint (security challenge)
		if a.isCheckpoint(err) {
			a.log.Warn("Security checkpoint detected, waiting before retry")
			// Exponential backoff
			backoff := time.Duration(attempt*attempt) * time.Minute
			time.Sleep(backoff)
		} else {
			// Other error, short delay
			time.Sleep(5 * time.Second)
		}
	}

	logger.Timing("auth", "login", start, lastErr)
	return fmt.Errorf("login failed after %d attempts: %w", a.config.CheckpointRetries, lastErr)
}

// performLogin executes the mock login flow
func (a *Authenticator) performLogin() error {
	a.log.Info("Executing login flow")

	// Get credentials from environment
	email := os.Getenv("LOGIN_EMAIL")
	password := os.Getenv("LOGIN_PASSWORD")

	if email == "" || password == "" {
		return fmt.Errorf("LOGIN_EMAIL and LOGIN_PASSWORD must be set")
	}

	// EDUCATIONAL NOTE: This is a MOCK flow demonstrating stealth integration
	// Real selectors would be used in production (which we deliberately don't provide)

	// Step 1: Navigate to login page
	a.log.Info("Navigating to login page")
	// In production: a.browser.Navigate("https://www.linkedin.com/login")
	a.stealth.RandomDelay()

	// Step 2: Wait for page to load
	a.stealth.WaitForPageLoad()

	// Step 3: Random scroll to simulate reading
	a.stealth.RandomScroll()

	// Step 4: Move mouse to email field (demonstrating Bézier movement)
	a.log.Info("Moving to email field")
	a.stealth.MoveMouse(400, 300) // Mock coordinates
	a.stealth.RandomDelay()

	// Step 5: Type email with human-like behavior
	a.log.Info("Entering email", "email", maskEmail(email))
	// In production: a.browser.Type("#username-field", email)
	a.stealth.TypeHumanLike("mock-email-selector", email)

	// Step 6: Move to password field
	a.stealth.WanderMouse() // Simulate mouse wandering
	a.stealth.MoveMouse(400, 400)
	a.stealth.RandomDelay()

	// Step 7: Type password
	a.log.Info("Entering password")
	// In production: a.browser.Type("#password-field", password)
	a.stealth.TypeHumanLike("mock-password-selector", "********") // Never log real password

	// Step 8: Thinking pause before submit
	a.stealth.ThinkingPause()

	// Step 9: Click login button
	a.log.Info("Clicking login button")
	a.stealth.MoveMouse(400, 500)
	// In production: a.browser.Click("#login-submit")
	a.stealth.RandomDelay()

	// Step 10: Wait for navigation
	a.log.Info("Waiting for login to complete")
	a.stealth.WaitForNavigation()

	// Step 11: Verify login success
	// In production: Check for presence of dashboard elements or profile menu
	// For PoC, we simulate success
	currentURL := "https://www.linkedin.com/feed/" // Mock
	a.log.Info("Login flow completed", "current_url", currentURL)

	// Simulate checkpoint detection randomly (10% chance for demo)
	if a.stealth.ShouldProceed(0.1) {
		return fmt.Errorf("checkpoint_detected: security verification required")
	}

	return nil
}

// tryLoadSession attempts to restore a previous session
func (a *Authenticator) tryLoadSession() error {
	a.log.Info("Attempting to load saved session", "path", a.config.SessionCookiePath)

	// Check if cookie file exists
	if _, err := os.Stat(a.config.SessionCookiePath); os.IsNotExist(err) {
		return fmt.Errorf("no session file found")
	}

	// Read cookies from file
	data, err := os.ReadFile(a.config.SessionCookiePath)
	if err != nil {
		return fmt.Errorf("failed to read session file: %w", err)
	}

	var cookies []*proto.NetworkCookie
	if err := json.Unmarshal(data, &cookies); err != nil {
		return fmt.Errorf("failed to parse session file: %w", err)
	}

	if len(cookies) == 0 {
		return fmt.Errorf("no cookies found in session file")
	}

	// Check if cookies are expired
	now := time.Now()
	validCookies := make([]*proto.NetworkCookie, 0)
	for _, cookie := range cookies {
		if cookie.Expires > 0 && time.Unix(int64(cookie.Expires), 0).Before(now) {
			a.log.Debug("Cookie expired", "name", cookie.Name)
			continue
		}
		validCookies = append(validCookies, cookie)
	}

	if len(validCookies) == 0 {
		return fmt.Errorf("all cookies expired")
	}

	// Set cookies in browser
	if len(validCookies) > 0 {
		if err := a.browser.SetCookies(validCookies); err != nil {
			return fmt.Errorf("failed to set cookies: %w", err)
		}
	}

	// Navigate to verify session
	// In production: a.browser.Navigate("https://www.linkedin.com/feed/")
	a.stealth.WaitForPageLoad()

	a.log.Info("Session loaded successfully", "cookies", len(validCookies))
	return nil
}

// saveSession saves the current session cookies
func (a *Authenticator) saveSession() error {
	a.log.Info("Saving session cookies")

	cookies, err := a.browser.GetCookies()
	if err != nil {
		return fmt.Errorf("failed to get cookies: %w", err)
	}

	// Serialize cookies to JSON
	data, err := json.MarshalIndent(cookies, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cookies: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(a.config.SessionCookiePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write to file
	if err := os.WriteFile(a.config.SessionCookiePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	a.log.Info("Session saved", "path", a.config.SessionCookiePath, "cookies", len(cookies))
	return nil
}

// isCheckpoint checks if an error indicates a security checkpoint
func (a *Authenticator) isCheckpoint(err error) bool {
	if err == nil {
		return false
	}
	
	// Check error message for checkpoint indicators
	msg := strings.ToLower(err.Error())
	checkpointIndicators := []string{
		"checkpoint",
		"verification",
		"security check",
		"captcha",
		"suspicious activity",
		"verify your identity",
	}

	for _, indicator := range checkpointIndicators {
		if strings.Contains(msg, indicator) {
			return true
		}
	}

	return false
}

// IsAuthenticated checks if currently logged in
// Delegates to browser controller which encapsulates session logic
func (a *Authenticator) IsAuthenticated() bool {
	return a.browser.HasValidSession()
}

// Logout clears the session (mock implementation)
func (a *Authenticator) Logout() error {
	a.log.Info("Logging out")

	// In production: Navigate to logout URL and wait
	// For PoC, just clear cookies
	if err := os.Remove(a.config.SessionCookiePath); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove session file: %w", err)
		}
	}

	a.log.Info("Logout complete")
	return nil
}

// maskEmail masks email for logging (privacy)
func maskEmail(email string) string {
	if len(email) < 3 {
		return "***"
	}
	
	atIndex := -1
	for i, c := range email {
		if c == '@' {
			atIndex = i
			break
		}
	}
	
	if atIndex <= 0 {
		return email[:2] + "***"
	}
	
	return email[:2] + "***" + email[atIndex:]
}
