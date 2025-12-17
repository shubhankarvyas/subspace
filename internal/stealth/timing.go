package stealth

import (
	"time"
)

// WaitForNavigation waits for page navigation to complete with human-like timing
// This wraps sleep with proper abstraction and adds variable timing
func (s *Stealth) WaitForNavigation() {
	// Variable wait time for navigation (2-4 seconds)
	delay := s.randomInt(2000, 4000)
	s.log.Debug("Waiting for navigation", "ms", delay)
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

// WaitForPageLoad waits for page to fully load with jitter
func (s *Stealth) WaitForPageLoad() {
	delay := s.randomInt(1500, 3000)
	s.log.Debug("Waiting for page load", "ms", delay)
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

// ShortPause adds a brief, randomized pause
func (s *Stealth) ShortPause() {
	delay := s.randomInt(200, 600)
	time.Sleep(time.Duration(delay) * time.Millisecond)
}
