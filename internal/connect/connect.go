package connect

import (
	"fmt"
	"time"

	"subspace/internal/browser"
	"subspace/internal/config"
	"subspace/internal/logger"
	"subspace/internal/stealth"
	"subspace/internal/storage"
)

/*
CONNECTION MODULE - EDUCATIONAL IMPLEMENTATION

Demonstrates a state machine for connection request management.
Shows how to enforce rate limits and track profile states.

STATE MACHINE:
discovered → requested → accepted → cooled_down
         ↓
      rejected

FEATURES:
- Daily/hourly connection limits
- Profile state tracking
- Cooldown period enforcement
- Personalized note support (see messaging module)
*/

// Connector handles connection request operations
type Connector struct {
	browser browser.Controller
	stealth *stealth.Stealth
	storage *storage.Storage
	limits  config.LimitsConfig
	log     *logger.ContextLogger
}

// New creates a new connector
func New(b browser.Controller, s *stealth.Stealth, storage *storage.Storage, limits config.LimitsConfig) *Connector {
	return &Connector{
		browser: b,
		stealth: s,
		storage: storage,
		limits:  limits,
		log:     logger.NewContext("connect"),
	}
}

// ProcessDailyConnections processes pending connection requests
func (c *Connector) ProcessDailyConnections() error {
	c.log.Info("Starting daily connection processing")
	start := time.Now()

	// Check daily and hourly limits
	connectionsToday := c.storage.GetActionCountToday("connection")
	connectionsLastHour := c.storage.GetActionCountLastHour("connection")

	c.log.Info("Current connection counts",
		"today", connectionsToday,
		"last_hour", connectionsLastHour,
		"limit_daily", c.limits.ConnectionsPerDay,
		"limit_hourly", c.limits.ConnectionsPerHour)

	// Check if we've hit daily limit
	if connectionsToday >= c.limits.ConnectionsPerDay {
		c.log.Warn("Daily connection limit reached, entering cooldown",
			"count", connectionsToday,
			"limit", c.limits.ConnectionsPerDay)
		
		// Log cooldown start
		cooldownUntil := time.Now().Add(time.Duration(c.limits.CooldownMinutes) * time.Minute)
		c.log.Info("Cooldown until", "time", cooldownUntil.Format(time.RFC3339))
		
		return nil
	}

	// Check if we've hit hourly limit
	if connectionsLastHour >= c.limits.ConnectionsPerHour {
		c.log.Warn("Hourly connection limit reached, waiting",
			"count", connectionsLastHour,
			"limit", c.limits.ConnectionsPerHour)
		return nil
	}

	// Get profiles in "discovered" state
	candidates := c.storage.GetProfilesByState(storage.StateDiscovered)
	c.log.Info("Found candidate profiles", "count", len(candidates))

	if len(candidates) == 0 {
		c.log.Info("No candidates to process")
		return nil
	}

	// Calculate how many we can send
	remainingDaily := c.limits.ConnectionsPerDay - connectionsToday
	remainingHourly := c.limits.ConnectionsPerHour - connectionsLastHour
	
	maxToSend := remainingDaily
	if remainingHourly < maxToSend {
		maxToSend = remainingHourly
	}

	c.log.Info("Planning to send connections", "max", maxToSend)

	// Process profiles
	sent := 0
	for i, profile := range candidates {
		if sent >= maxToSend {
			c.log.Info("Reached send limit for this batch", "sent", sent)
			break
		}

		c.log.Info("Processing profile",
			"index", i+1,
			"total", len(candidates),
			"name", profile.Name)

		// Send connection request
		if err := c.SendConnectionRequest(profile); err != nil {
			c.log.Error("Failed to send connection request",
				"profile", profile.Name,
				"error", err)
			
			// Log failed action
			c.storage.LogAction("connection", profile.ID, false, err)
			
			// Don't stop on error, continue with next
			continue
		}

		sent++
		
		// Enforce cooldown between requests (stealth)
		c.stealth.EnforceCooldown("connection", 30) // 30 seconds minimum
	}

	logger.Timing("connect", "process_daily", start, nil)
	c.log.Info("Daily connection processing complete",
		"sent", sent,
		"remaining_daily", remainingDaily-sent)

	return nil
}

// SendConnectionRequest sends a connection request to a profile
func (c *Connector) SendConnectionRequest(profile *storage.Profile) error {
	c.log.Info("Sending connection request", "name", profile.Name, "profile_id", profile.ID)
	start := time.Now()

	// Step 1: Navigate to profile
	c.log.Debug("Navigating to profile", "url", profile.ProfileURL)
	// In production: c.browser.Navigate(profile.ProfileURL)
	c.stealth.RandomDelay()

	// Step 2: Wait for page load and scroll around (human-like)
	c.stealth.ThinkingPause()
	c.stealth.RandomScroll()
	c.stealth.WanderMouse()

	// Step 3: Look for the "Connect" button
	c.log.Debug("Looking for Connect button")
	// EDUCATIONAL NOTE: In production:
	// connectBtn := c.browser.Page.Element("[aria-label='Invite ... to connect']")
	
	// Step 4: Move mouse to button
	c.stealth.MoveMouse(800, 400) // Mock coordinates
	c.stealth.RandomDelay()

	// Step 5: Click connect button
	c.log.Debug("Clicking Connect button")
	// In production: c.browser.Click(connectBtn selector)
	
	// Step 6: Handle "Add a note" dialog (if appears)
	c.stealth.ThinkingPause()
	
	// Check if we should add a personalized note
	// For now, send without note (can be enhanced with messaging module)
	c.log.Debug("Sending without note")
	
	// Step 7: Click "Send" button in dialog
	c.stealth.MoveMouse(700, 500)
	c.stealth.RandomDelay()
	// In production: c.browser.Click("[aria-label='Send invitation']")

	// Step 8: Wait for confirmation
	c.stealth.RandomDelay()

	// Step 9: Update profile state
	now := time.Now()
	profile.State = storage.StateRequested
	profile.RequestedAt = &now

	if err := c.storage.SaveProfile(profile); err != nil {
		logger.Timing("connect", "send_request", start, err)
		return fmt.Errorf("failed to update profile state: %w", err)
	}

	// Log successful action
	c.storage.LogAction("connection", profile.ID, true, nil)

	logger.Timing("connect", "send_request", start, nil)
	c.log.Info("Connection request sent successfully", "profile", profile.Name)

	return nil
}

// CheckAcceptedConnections checks for newly accepted connections
func (c *Connector) CheckAcceptedConnections() error {
	c.log.Info("Checking for accepted connections")

	// Get profiles in "requested" state
	requested := c.storage.GetProfilesByState(storage.StateRequested)
	c.log.Info("Profiles awaiting acceptance", "count", len(requested))

	// EDUCATIONAL NOTE: In production, this would:
	// 1. Navigate to "My Network" page
	// 2. Check connection list
	// 3. Compare with requested profiles
	// 4. Update states for accepted ones
	//
	// For PoC, we simulate random acceptances

	accepted := 0
	for _, profile := range requested {
		// Simulate 20% chance of acceptance (for demo purposes)
		if c.stealth.ShouldProceed(0.2) {
			now := time.Now()
			profile.State = storage.StateAccepted
			profile.AcceptedAt = &now

			if err := c.storage.SaveProfile(profile); err != nil {
				c.log.Error("Failed to update profile", "error", err)
				continue
			}

			c.log.Info("Connection accepted", "name", profile.Name)
			accepted++
		}
	}

	c.log.Info("Acceptance check complete", "newly_accepted", accepted)
	return nil
}

// MoveToCooldown moves a profile to cooled_down state
func (c *Connector) MoveToCooldown(profile *storage.Profile) error {
	c.log.Info("Moving profile to cooldown", "name", profile.Name)

	now := time.Now()
	profile.State = storage.StateCooledDown
	profile.CooledDownAt = &now

	if err := c.storage.SaveProfile(profile); err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}

// WithdrawConnectionRequest withdraws a pending connection request
func (c *Connector) WithdrawConnectionRequest(profile *storage.Profile) error {
	c.log.Info("Withdrawing connection request", "name", profile.Name)
	start := time.Now()

	// EDUCATIONAL NOTE: In production:
	// 1. Navigate to "Sent" invitations
	// 2. Find the profile
	// 3. Click "Withdraw" button
	// 4. Confirm withdrawal

	// Mock withdrawal
	c.stealth.RandomDelay()

	// Update state
	profile.State = storage.StateDiscovered // Reset to discovered
	profile.RequestedAt = nil

	if err := c.storage.SaveProfile(profile); err != nil {
		logger.Timing("connect", "withdraw", start, err)
		return fmt.Errorf("failed to update profile: %w", err)
	}

	logger.Timing("connect", "withdraw", start, nil)
	c.log.Info("Connection request withdrawn", "profile", profile.Name)

	return nil
}

// GetPendingRequests returns profiles awaiting acceptance
func (c *Connector) GetPendingRequests() []*storage.Profile {
	return c.storage.GetProfilesByState(storage.StateRequested)
}

// GetAcceptedConnections returns accepted connections
func (c *Connector) GetAcceptedConnections() []*storage.Profile {
	return c.storage.GetProfilesByState(storage.StateAccepted)
}

// CanSendMore checks if we can send more connections today
func (c *Connector) CanSendMore() bool {
	today := c.storage.GetActionCountToday("connection")
	hourly := c.storage.GetActionCountLastHour("connection")

	return today < c.limits.ConnectionsPerDay && hourly < c.limits.ConnectionsPerHour
}

// GetStats returns connection statistics
func (c *Connector) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"connections_today":      c.storage.GetActionCountToday("connection"),
		"connections_last_hour":  c.storage.GetActionCountLastHour("connection"),
		"pending_requests":       len(c.GetPendingRequests()),
		"accepted_connections":   len(c.GetAcceptedConnections()),
		"limit_daily":            c.limits.ConnectionsPerDay,
		"limit_hourly":           c.limits.ConnectionsPerHour,
		"can_send_more":          c.CanSendMore(),
	}
}
