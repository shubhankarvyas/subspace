package browser

import (
	"time"

	"github.com/go-rod/rod/lib/proto"
)

// Controller defines the interface for browser operations
// This abstraction prevents business logic from depending on Rod directly
type Controller interface {
	// Navigation
	Navigate(url string) error
	WaitForElement(selector string, timeout time.Duration) error
	GetCurrentURL() string
	
	// Element Interaction
	Click(selector string) error
	Type(selector, text string) error
	GetText(selector string) (string, error)
	GetAttribute(selector, attribute string) (string, error)
	IsElementPresent(selector string) bool
	WaitVisible(selector string) error
	
	// Session Management
	GetCookies() ([]*proto.NetworkCookie, error)
	SetCookies(cookies []*proto.NetworkCookie) error
	HasValidSession() bool
	
	// Utilities
	Screenshot(path string) error
	ExecuteScript(script string) (interface{}, error)
	
	// Lifecycle
	Close() error
}
