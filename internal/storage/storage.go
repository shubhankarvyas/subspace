package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// ProfileState represents the state of a profile in the connection pipeline
type ProfileState string

const (
	StateDiscovered  ProfileState = "discovered"
	StateRequested   ProfileState = "requested"
	StateAccepted    ProfileState = "accepted"
	StateCooledDown  ProfileState = "cooled_down"
	StateRejected    ProfileState = "rejected"
)

// Profile represents a target profile
type Profile struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Title        string       `json:"title"`
	Company      string       `json:"company"`
	ProfileURL   string       `json:"profile_url"`
	State        ProfileState `json:"state"`
	DiscoveredAt time.Time    `json:"discovered_at"`
	RequestedAt  *time.Time   `json:"requested_at,omitempty"`
	AcceptedAt   *time.Time   `json:"accepted_at,omitempty"`
	CooledDownAt *time.Time   `json:"cooled_down_at,omitempty"`
	SearchQuery  string       `json:"search_query"`
	Notes        string       `json:"notes"`
}

// Message represents a message sent to a connection
type Message struct {
	ID          string    `json:"id"`
	ProfileID   string    `json:"profile_id"`
	Content     string    `json:"content"`
	SentAt      time.Time `json:"sent_at"`
	Template    string    `json:"template"`
}

// ActionLog tracks all automated actions for rate limiting
type ActionLog struct {
	Action    string    `json:"action"` // "connection", "message", "search"
	Timestamp time.Time `json:"timestamp"`
	ProfileID string    `json:"profile_id,omitempty"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
}

// Storage handles all data persistence using JSON
type Storage struct {
	path      string
	data      *Data
	mu        sync.RWMutex
}

// Data represents the complete storage structure
type Data struct {
	Profiles   map[string]*Profile  `json:"profiles"`
	Messages   map[string]*Message  `json:"messages"`
	ActionLogs []ActionLog          `json:"action_logs"`
	LastSync   time.Time            `json:"last_sync"`
}

// New creates a new storage instance
func New(path string) (*Storage, error) {
	s := &Storage{
		path: path,
		data: &Data{
			Profiles:   make(map[string]*Profile),
			Messages:   make(map[string]*Message),
			ActionLogs: make([]ActionLog, 0),
		},
	}

	// Load existing data if available
	if err := s.load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load storage: %w", err)
		}
		// File doesn't exist, start fresh
		if err := s.save(); err != nil {
			return nil, fmt.Errorf("failed to initialize storage: %w", err)
		}
	}

	return s, nil
}

// load reads data from disk
func (s *Storage) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, s.data)
}

// save writes data to disk
func (s *Storage) save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data.LastSync = time.Now()

	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(s.path[:len(s.path)-len("/db.json")], 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	return os.WriteFile(s.path, data, 0644)
}

// SaveProfile saves or updates a profile
func (s *Storage) SaveProfile(profile *Profile) error {
	s.mu.Lock()
	s.data.Profiles[profile.ID] = profile
	s.mu.Unlock()
	return s.save()
}

// GetProfile retrieves a profile by ID
func (s *Storage) GetProfile(id string) (*Profile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	profile, exists := s.data.Profiles[id]
	if !exists {
		return nil, fmt.Errorf("profile not found: %s", id)
	}
	return profile, nil
}

// GetProfilesByState retrieves all profiles in a given state
func (s *Storage) GetProfilesByState(state ProfileState) []*Profile {
	s.mu.RLock()
	defer s.mu.RUnlock()

	profiles := make([]*Profile, 0)
	for _, profile := range s.data.Profiles {
		if profile.State == state {
			profiles = append(profiles, profile)
		}
	}
	return profiles
}

// ProfileExists checks if a profile URL has been seen before (deduplication)
func (s *Storage) ProfileExists(profileURL string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, profile := range s.data.Profiles {
		if profile.ProfileURL == profileURL {
			return true
		}
	}
	return false
}

// SaveMessage saves a message record
func (s *Storage) SaveMessage(message *Message) error {
	s.mu.Lock()
	s.data.Messages[message.ID] = message
	s.mu.Unlock()
	return s.save()
}

// GetMessagesByProfile retrieves all messages for a profile
func (s *Storage) GetMessagesByProfile(profileID string) []*Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	messages := make([]*Message, 0)
	for _, msg := range s.data.Messages {
		if msg.ProfileID == profileID {
			messages = append(messages, msg)
		}
	}
	return messages
}

// LogAction records an action for rate limiting purposes
func (s *Storage) LogAction(action, profileID string, success bool, err error) error {
	s.mu.Lock()
	
	log := ActionLog{
		Action:    action,
		Timestamp: time.Now(),
		ProfileID: profileID,
		Success:   success,
	}
	if err != nil {
		log.Error = err.Error()
	}
	
	s.data.ActionLogs = append(s.data.ActionLogs, log)
	s.mu.Unlock()
	
	return s.save()
}

// GetActionCountSince returns the count of successful actions since a given time
func (s *Storage) GetActionCountSince(action string, since time.Time) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, log := range s.data.ActionLogs {
		if log.Action == action && log.Success && log.Timestamp.After(since) {
			count++
		}
	}
	return count
}

// GetActionCountToday returns today's action count
func (s *Storage) GetActionCountToday(action string) int {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return s.GetActionCountSince(action, startOfDay)
}

// GetActionCountLastHour returns the last hour's action count
func (s *Storage) GetActionCountLastHour(action string) int {
	return s.GetActionCountSince(action, time.Now().Add(-1*time.Hour))
}

// CleanOldLogs removes action logs older than retention period (to prevent unbounded growth)
func (s *Storage) CleanOldLogs(retentionDays int) error {
	s.mu.Lock()
	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	
	filtered := make([]ActionLog, 0)
	for _, log := range s.data.ActionLogs {
		if log.Timestamp.After(cutoff) {
			filtered = append(filtered, log)
		}
	}
	s.data.ActionLogs = filtered
	s.mu.Unlock()
	
	return s.save()
}

// GetStats returns summary statistics
func (s *Storage) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := map[string]interface{}{
		"total_profiles":         len(s.data.Profiles),
		"discovered":             0,
		"requested":              0,
		"accepted":               0,
		"cooled_down":            0,
		"rejected":               0,
		"total_messages":         len(s.data.Messages),
		"connections_today":      s.GetActionCountToday("connection"),
		"messages_today":         s.GetActionCountToday("message"),
		"connections_last_hour":  s.GetActionCountLastHour("connection"),
	}

	for _, profile := range s.data.Profiles {
		switch profile.State {
		case StateDiscovered:
			stats["discovered"] = stats["discovered"].(int) + 1
		case StateRequested:
			stats["requested"] = stats["requested"].(int) + 1
		case StateAccepted:
			stats["accepted"] = stats["accepted"].(int) + 1
		case StateCooledDown:
			stats["cooled_down"] = stats["cooled_down"].(int) + 1
		case StateRejected:
			stats["rejected"] = stats["rejected"].(int) + 1
		}
	}

	return stats
}
