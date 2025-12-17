package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the complete application configuration
type Config struct {
	App     AppConfig     `yaml:"app"`
	Stealth StealthConfig `yaml:"stealth"`
	Limits  LimitsConfig  `yaml:"limits"`
	Auth    AuthConfig    `yaml:"auth"`
	Search  SearchConfig  `yaml:"search"`
}

// AppConfig contains general application settings
type AppConfig struct {
	DataDir   string `yaml:"data_dir"`
	LogLevel  string `yaml:"log_level"`
	Headless  bool   `yaml:"headless"`
	UserAgent string `yaml:"user_agent"`
}

// StealthConfig contains anti-detection configuration
// Each technique can be fine-tuned independently
type StealthConfig struct {
	// Mouse Movement Configuration
	MouseSpeed         float64 `yaml:"mouse_speed"`          // Pixels per second (200-400 is human-like)
	MouseWanderEnabled bool    `yaml:"mouse_wander_enabled"` // Random hover movements
	MouseWanderChance  float64 `yaml:"mouse_wander_chance"`  // 0.0-1.0 probability

	// Typing Configuration
	TypingSpeedMin int     `yaml:"typing_speed_min"` // Milliseconds per keystroke
	TypingSpeedMax int     `yaml:"typing_speed_max"`
	TypoChance     float64 `yaml:"typo_chance"`      // 0.0-1.0 probability of making a typo
	TypoCorrection bool    `yaml:"typo_correction"`  // Auto-correct typos with backspace

	// Timing & Jitter
	ActionDelayMin int `yaml:"action_delay_min"` // Milliseconds between actions
	ActionDelayMax int `yaml:"action_delay_max"`
	ThinkTimeMin   int `yaml:"think_time_min"`   // Longer pauses simulating "thinking"
	ThinkTimeMax   int `yaml:"think_time_max"`

	// Scrolling Behavior
	ScrollEnabled      bool    `yaml:"scroll_enabled"`
	ScrollChance       float64 `yaml:"scroll_chance"`        // Chance to scroll before action
	ScrollDistance     int     `yaml:"scroll_distance"`      // Pixels per scroll
	ScrollAcceleration float64 `yaml:"scroll_acceleration"`  // Simulate acceleration/deceleration

	// Business Hours & Scheduling
	BusinessHoursEnabled bool   `yaml:"business_hours_enabled"`
	BusinessHoursStart   string `yaml:"business_hours_start"` // HH:MM format
	BusinessHoursEnd     string `yaml:"business_hours_end"`
	BreakTimeEnabled     bool   `yaml:"break_time_enabled"`
	BreakTimeStart       string `yaml:"break_time_start"`
	BreakTimeEnd         string `yaml:"break_time_end"`

	// Fingerprint Masking
	MaskWebDriver    bool `yaml:"mask_webdriver"`     // Hide webdriver flag
	MaskChrome       bool `yaml:"mask_chrome"`        // Hide automation indicators
	RandomViewport   bool `yaml:"random_viewport"`    // Randomize browser window size
	ViewportWidthMin int  `yaml:"viewport_width_min"`
	ViewportWidthMax int  `yaml:"viewport_width_max"`
	ViewportHeightMin int  `yaml:"viewport_height_min"`
	ViewportHeightMax int  `yaml:"viewport_height_max"`
}

// LimitsConfig enforces rate limiting and safety boundaries
type LimitsConfig struct {
	ConnectionsPerDay  int `yaml:"connections_per_day"`
	ConnectionsPerHour int `yaml:"connections_per_hour"`
	MessagesPerDay     int `yaml:"messages_per_day"`
	SearchesPerDay     int `yaml:"searches_per_day"`
	CooldownMinutes    int `yaml:"cooldown_minutes"` // After daily limit reached
}

// AuthConfig contains authentication-related settings
type AuthConfig struct {
	SessionCookiePath string `yaml:"session_cookie_path"`
	ReuseSession      bool   `yaml:"reuse_session"`
	CheckpointRetries int    `yaml:"checkpoint_retries"`
}

// SearchConfig contains search behavior settings
type SearchConfig struct {
	ResultsPerPage      int      `yaml:"results_per_page"`
	MaxPages            int      `yaml:"max_pages"`
	DeduplicationWindow int      `yaml:"deduplication_window"` // Days to remember seen profiles
	DefaultKeywords     []string `yaml:"default_keywords"`
}

// Load reads and parses the configuration file
func Load(path string) (*Config, error) {
	// Set defaults
	cfg := &Config{
		App: AppConfig{
			DataDir:   "./data",
			LogLevel:  "info",
			Headless:  false,
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		},
		Stealth: StealthConfig{
			MouseSpeed:            300.0,
			MouseWanderEnabled:    true,
			MouseWanderChance:     0.15,
			TypingSpeedMin:        80,
			TypingSpeedMax:        200,
			TypoChance:            0.03,
			TypoCorrection:        true,
			ActionDelayMin:        500,
			ActionDelayMax:        2000,
			ThinkTimeMin:          2000,
			ThinkTimeMax:          5000,
			ScrollEnabled:         true,
			ScrollChance:          0.3,
			ScrollDistance:        300,
			ScrollAcceleration:    0.8,
			BusinessHoursEnabled:  true,
			BusinessHoursStart:    "09:00",
			BusinessHoursEnd:      "17:00",
			BreakTimeEnabled:      true,
			BreakTimeStart:        "12:00",
			BreakTimeEnd:          "13:00",
			MaskWebDriver:         true,
			MaskChrome:            true,
			RandomViewport:        true,
			ViewportWidthMin:      1200,
			ViewportWidthMax:      1920,
			ViewportHeightMin:     800,
			ViewportHeightMax:     1080,
		},
		Limits: LimitsConfig{
			ConnectionsPerDay:  50,
			ConnectionsPerHour: 10,
			MessagesPerDay:     30,
			SearchesPerDay:     20,
			CooldownMinutes:    60,
		},
		Auth: AuthConfig{
			SessionCookiePath: "./data/session.json",
			ReuseSession:      true,
			CheckpointRetries: 3,
		},
		Search: SearchConfig{
			ResultsPerPage:      25,
			MaxPages:            10,
			DeduplicationWindow: 30,
			DefaultKeywords:     []string{"software engineer", "golang developer"},
		},
	}

	// Override with file if exists
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Validate checks configuration values for correctness
func (c *Config) Validate() error {
	// Validate log level
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[c.App.LogLevel] {
		return fmt.Errorf("invalid log level: %s (must be debug, info, warn, or error)", c.App.LogLevel)
	}

	// Validate business hours format
	if c.Stealth.BusinessHoursEnabled {
		if _, err := time.Parse("15:04", c.Stealth.BusinessHoursStart); err != nil {
			return fmt.Errorf("invalid business_hours_start format: %s (use HH:MM)", c.Stealth.BusinessHoursStart)
		}
		if _, err := time.Parse("15:04", c.Stealth.BusinessHoursEnd); err != nil {
			return fmt.Errorf("invalid business_hours_end format: %s (use HH:MM)", c.Stealth.BusinessHoursEnd)
		}
	}

	// Validate limits
	if c.Limits.ConnectionsPerDay <= 0 || c.Limits.ConnectionsPerDay > 100 {
		return fmt.Errorf("connections_per_day must be between 1 and 100")
	}
	if c.Limits.ConnectionsPerHour > c.Limits.ConnectionsPerDay {
		return fmt.Errorf("connections_per_hour cannot exceed connections_per_day")
	}

	return nil
}

// GetEnv reads environment variables with fallback
func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
