package browser

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
	
	"subspace/internal/config"
	"subspace/internal/logger"
)

// Browser wraps Rod browser functionality with a clean interface
// This abstraction prevents business logic from directly calling Rod APIs
type Browser struct {
	browser *rod.Browser
	Page    *rod.Page
	config  config.AppConfig
	log     *logger.ContextLogger
}

// New creates a new browser instance with stealth configuration
func New(cfg config.AppConfig) (*Browser, error) {
	log := logger.NewContext("browser")
	
	log.Info("Initializing browser", "headless", cfg.Headless)
	
	// Launch browser with configured options
	l := launcher.New().
		Headless(cfg.Headless).
		UserDataDir("") // Don't persist user data by default

	// Start the launcher
	url, err := l.Launch()
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	// Connect to browser
	browser := rod.New().ControlURL(url)
	if err := browser.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to browser: %w", err)
	}

	// Create a new page
	page, err := stealth.Page(browser)
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}

	// Set user agent
	if cfg.UserAgent != "" {
		if err := page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
			UserAgent: cfg.UserAgent,
		}); err != nil {
			log.Warn("Failed to set user agent", "error", err)
		}
	}

	b := &Browser{
		browser: browser,
		Page:    page,
		config:  cfg,
		log:     log,
	}

	log.Info("Browser initialized successfully")
	return b, nil
}

// Navigate navigates to a URL with error handling
func (b *Browser) Navigate(url string) error {
	b.log.Info("Navigating to URL", "url", url)
	start := time.Now()
	
	if err := b.Page.Navigate(url); err != nil {
		logger.Timing("browser", "navigate", start, err)
		return fmt.Errorf("failed to navigate: %w", err)
	}
	
	// Wait for page load
	if err := b.Page.WaitLoad(); err != nil {
		logger.Timing("browser", "navigate", start, err)
		return fmt.Errorf("page load timeout: %w", err)
	}
	
	logger.Timing("browser", "navigate", start, nil)
	return nil
}

// WaitForElement waits for an element to be visible (mock implementation)
// In production, this would use real selectors
func (b *Browser) WaitForElement(selector string, timeout time.Duration) error {
	b.log.Debug("Waiting for element", "selector", selector, "timeout", timeout)
	
	// EDUCATIONAL NOTE: In a real implementation, this would use:
	// element, err := b.Page.Timeout(timeout).Element(selector)
	// For this PoC, we simulate the wait
	
	time.Sleep(500 * time.Millisecond) // Simulate wait
	
	// Return success for demo purposes
	return nil
}

// Click performs a click action (mock implementation)
// In production, this would find and click real elements
func (b *Browser) Click(selector string) error {
	b.log.Debug("Clicking element", "selector", selector)
	
	// EDUCATIONAL NOTE: Real implementation would be:
	// element, err := b.Page.Element(selector)
	// if err != nil { return err }
	// return element.Click(proto.InputMouseButtonLeft)
	
	// For PoC, we just log the action
	b.log.Info("Mock click executed", "selector", selector)
	return nil
}

// Type simulates typing text (mock implementation)
// Actual typing with human-like behavior is handled by stealth package
func (b *Browser) Type(selector, text string) error {
	b.log.Debug("Typing into element", "selector", selector, "text_length", len(text))
	
	// EDUCATIONAL NOTE: Real implementation would be:
	// element, err := b.Page.Element(selector)
	// if err != nil { return err }
	// return element.Input(text)
	
	// For PoC, we just log the action
	b.log.Info("Mock type executed", "selector", selector, "text_length", len(text))
	return nil
}

// GetText retrieves text from an element (mock implementation)
func (b *Browser) GetText(selector string) (string, error) {
	b.log.Debug("Getting text from element", "selector", selector)
	
	// EDUCATIONAL NOTE: Real implementation would be:
	// element, err := b.Page.Element(selector)
	// if err != nil { return "", err }
	// return element.Text()
	
	// Return mock data for demo
	return "Mock text content", nil
}

// GetAttribute retrieves an attribute from an element (mock implementation)
func (b *Browser) GetAttribute(selector, attribute string) (string, error) {
	b.log.Debug("Getting attribute", "selector", selector, "attribute", attribute)
	
	// EDUCATIONAL NOTE: Real implementation would use element.Attribute()
	
	return "mock-value", nil
}

// Screenshot captures a screenshot of the current page
func (b *Browser) Screenshot(path string) error {
	b.log.Info("Taking screenshot", "path", path)
	
	data, err := b.Page.Screenshot(false, nil)
	if err != nil {
		return fmt.Errorf("failed to capture screenshot: %w", err)
	}
	
	// In a real implementation, save to disk
	_ = data
	b.log.Info("Screenshot captured", "size_bytes", len(data))
	return nil
}

// GetCookies retrieves current cookies
func (b *Browser) GetCookies() ([]*proto.NetworkCookie, error) {
	b.log.Debug("Retrieving cookies")
	
	cookies, err := b.Page.Cookies([]string{})
	if err != nil {
		return nil, fmt.Errorf("failed to get cookies: %w", err)
	}
	
	b.log.Info("Retrieved cookies", "count", len(cookies))
	return cookies, nil
}

// SetCookies sets cookies for the page
func (b *Browser) SetCookies(cookies []*proto.NetworkCookie) error {
	b.log.Info("Setting cookies", "count", len(cookies))
	
	// Convert NetworkCookie to NetworkCookieParam
	params := make([]*proto.NetworkCookieParam, len(cookies))
	for i, cookie := range cookies {
		params[i] = &proto.NetworkCookieParam{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Expires:  cookie.Expires,
			HTTPOnly: cookie.HTTPOnly,
			Secure:   cookie.Secure,
			SameSite: cookie.SameSite,
		}
	}
	
	if err := b.Page.SetCookies(params); err != nil {
		return fmt.Errorf("failed to set cookies: %w", err)
	}
	
	return nil
}

// ExecuteScript runs JavaScript in the page context (mock)
func (b *Browser) ExecuteScript(script string) (interface{}, error) {
	b.log.Debug("Executing script")
	
	// EDUCATIONAL NOTE: Real implementation:
	// return b.Page.Eval(script)
	
	b.log.Info("Mock script executed")
	return nil, nil
}

// GetCurrentURL returns the current page URL
func (b *Browser) GetCurrentURL() string {
	info := b.Page.MustInfo()
	return info.URL
}

// Close gracefully closes the browser
func (b *Browser) Close() error {
	b.log.Info("Closing browser")
	
	if b.Page != nil {
		if err := b.Page.Close(); err != nil {
			b.log.Warn("Error closing page", "error", err)
		}
	}
	
	if b.browser != nil {
		if err := b.browser.Close(); err != nil {
			return fmt.Errorf("failed to close browser: %w", err)
		}
	}
	
	b.log.Info("Browser closed successfully")
	return nil
}

// WaitVisible waits for an element to become visible
func (b *Browser) WaitVisible(selector string) error {
	return b.WaitForElement(selector, 10*time.Second)
}

// IsElementPresent checks if an element exists (mock)
func (b *Browser) IsElementPresent(selector string) bool {
	b.log.Debug("Checking element presence", "selector", selector)
	
	// EDUCATIONAL NOTE: Real implementation:
	// _, err := b.Page.Element(selector)
	// return err == nil
	
	// For demo, randomly return true/false
	return true
}

// HasValidSession checks if browser has a valid authenticated session
func (b *Browser) HasValidSession() bool {
	b.log.Debug("Checking session validity")
	
	cookies, err := b.GetCookies()
	if err != nil {
		return false
	}
	
	// Look for session cookies (mock - platform specific in production)
	for _, cookie := range cookies {
		if cookie.Name == "li_at" || cookie.Name == "JSESSIONID" {
			if cookie.Expires > 0 {
				expiry := time.Unix(int64(cookie.Expires), 0)
				if expiry.After(time.Now()) {
					return true
				}
			}
		}
	}
	
	return false
}
