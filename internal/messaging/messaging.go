package messaging

import (
	"fmt"
	"strings"
	"time"

	"subspace/internal/browser"
	"subspace/internal/config"
	"subspace/internal/logger"
	"subspace/internal/stealth"
	"subspace/internal/storage"
)

// Messenger handles message sending operations
type Messenger struct {
	browser   browser.Controller
	stealth   *stealth.Stealth
	storage   *storage.Storage
	limits    config.LimitsConfig
	templates map[string]string
	log       *logger.ContextLogger
}

// New creates a new messenger with default templates
func New(b browser.Controller, s *stealth.Stealth, storage *storage.Storage, limits config.LimitsConfig) *Messenger {
	m := &Messenger{
		browser:   b,
		stealth:   s,
		storage:   storage,
		limits:    limits,
		templates: make(map[string]string),
		log:       logger.NewContext("messaging"),
	}

	// Load default templates
	m.loadDefaultTemplates()

	return m
}

// loadDefaultTemplates sets up default message templates
func (m *Messenger) loadDefaultTemplates() {
	m.templates["follow_up"] = `Hi {{.Name}},

Thanks for connecting! I noticed your background in {{.Title}} at {{.Company}}.

I'm always interested in connecting with professionals in the field. Would love to stay in touch!

Best regards`

	m.templates["introduction"] = `Hi {{.Name}},

I came across your profile and was impressed by your experience in {{.Title}}.

I'm working on some interesting projects and thought we might have synergies to explore.

Looking forward to connecting!`

	m.templates["follow_up_short"] = `Hi {{.Name}}, thanks for connecting! Looking forward to staying in touch.`

	m.log.Info("Loaded message templates", "count", len(m.templates))
}

// SendMessage sends a message to a connected profile
func (m *Messenger) SendMessage(profile *storage.Profile, templateName string) error {
	m.log.Info("Sending message", "profile", profile.Name, "template", templateName)
	start := time.Now()

	// Check message limits
	messagesToday := m.storage.GetActionCountToday("message")
	if messagesToday >= m.limits.MessagesPerDay {
		err := fmt.Errorf("daily message limit reached: %d/%d", messagesToday, m.limits.MessagesPerDay)
		m.log.Warn("Cannot send message", "error", err)
		return err
	}

	// Check if profile has accepted connection
	if profile.State != storage.StateAccepted && profile.State != storage.StateCooledDown {
		return fmt.Errorf("cannot message profile in state: %s", profile.State)
	}

	// Check if we've already messaged this profile
	existingMessages := m.storage.GetMessagesByProfile(profile.ID)
	if len(existingMessages) > 0 {
		m.log.Info("Profile already messaged", "count", len(existingMessages))
		// Could add logic to allow follow-up messages after certain time
	}

	// Generate personalized message
	content, err := m.renderTemplate(templateName, profile)
	if err != nil {
		logger.Timing("messaging", "send_message", start, err)
		return fmt.Errorf("failed to render template: %w", err)
	}

	m.log.Debug("Rendered message", "length", len(content))

	// Navigate to messaging with profile
	if err := m.navigateToConversation(profile); err != nil {
		logger.Timing("messaging", "send_message", start, err)
		return fmt.Errorf("failed to navigate: %w", err)
	}

	// Type and send message
	if err := m.typeAndSend(content); err != nil {
		logger.Timing("messaging", "send_message", start, err)
		return fmt.Errorf("failed to send message: %w", err)
	}

	// Save message record
	message := &storage.Message{
		ID:        fmt.Sprintf("msg-%d", time.Now().UnixNano()),
		ProfileID: profile.ID,
		Content:   content,
		SentAt:    time.Now(),
		Template:  templateName,
	}

	if err := m.storage.SaveMessage(message); err != nil {
		m.log.Error("Failed to save message record", "error", err)
		// Don't fail the operation, message was sent
	}

	// Log action for rate limiting
	m.storage.LogAction("message", profile.ID, true, nil)

	logger.Timing("messaging", "send_message", start, nil)
	m.log.Info("Message sent successfully", "profile", profile.Name)

	return nil
}

// renderTemplate fills in template variables with profile data
func (m *Messenger) renderTemplate(templateName string, profile *storage.Profile) (string, error) {
	template, exists := m.templates[templateName]
	if !exists {
		return "", fmt.Errorf("template not found: %s", templateName)
	}

	content := template
	content = strings.ReplaceAll(content, "{{.Name}}", profile.Name)
	content = strings.ReplaceAll(content, "{{.Title}}", profile.Title)
	content = strings.ReplaceAll(content, "{{.Company}}", profile.Company)

	return content, nil
}

// navigateToConversation opens the messaging conversation with a profile
func (m *Messenger) navigateToConversation(profile *storage.Profile) error {
	m.log.Debug("Navigating to conversation", "profile", profile.Name)

	// EDUCATIONAL NOTE: In production:
	// 1. Navigate to LinkedIn messaging
	// 2. Search for the profile by name
	// 3. Open the conversation
	//
	// Alternative: Construct direct message URL if profile ID is known
	// messageURL := fmt.Sprintf("https://www.linkedin.com/messaging/thread/xxx/")

	// Mock navigation
	m.stealth.RandomDelay()
	m.stealth.WaitForPageLoad()

	return nil
}

// typeAndSend types the message and sends it
func (m *Messenger) typeAndSend(content string) error {
	m.log.Debug("Typing and sending message")

	// Step 1: Focus on message box
	m.stealth.MoveMouse(500, 600) // Mock coordinates
	m.stealth.RandomDelay()
	// In production: m.browser.Click(".msg-form__contenteditable")

	// Step 2: Type message with human-like behavior
	m.stealth.ThinkingPause() // Pause before typing (composing message)
	m.stealth.TypeHumanLike("mock-message-input", content)

	// Step 3: Pause before sending (reviewing message)
	m.stealth.ThinkingPause()

	// Step 4: Move to send button
	m.stealth.MoveMouse(700, 700)
	m.stealth.RandomDelay()

	// Step 5: Click send
	// In production: m.browser.Click(".msg-form__send-button")
	m.log.Debug("Message sent")

	return nil
}

// SendBulkMessages sends messages to multiple profiles
func (m *Messenger) SendBulkMessages(profiles []*storage.Profile, templateName string) error {
	m.log.Info("Starting bulk messaging", "count", len(profiles), "template", templateName)
	
	sent := 0
	failed := 0

	for i, profile := range profiles {
		m.log.Info("Processing profile", "index", i+1, "total", len(profiles))

		// Check if we've hit daily limit
		messagesToday := m.storage.GetActionCountToday("message")
		if messagesToday >= m.limits.MessagesPerDay {
			m.log.Warn("Daily limit reached, stopping bulk send",
				"sent", sent,
				"remaining", len(profiles)-i)
			break
		}

		// Send message
		if err := m.SendMessage(profile, templateName); err != nil {
			m.log.Error("Failed to send message", "profile", profile.Name, "error", err)
			failed++
			continue
		}

		sent++

		// Enforce cooldown between messages
		m.stealth.EnforceCooldown("message", 60) // 60 seconds minimum between messages
	}

	m.log.Info("Bulk messaging complete",
		"sent", sent,
		"failed", failed,
		"total", len(profiles))

	return nil
}

// ProcessAcceptedConnections sends follow-up messages to newly accepted connections
func (m *Messenger) ProcessAcceptedConnections() error {
	m.log.Info("Processing accepted connections for messaging")

	// Get accepted connections that haven't been messaged yet
	accepted := m.storage.GetProfilesByState(storage.StateAccepted)
	
	unmessaged := make([]*storage.Profile, 0)
	for _, profile := range accepted {
		messages := m.storage.GetMessagesByProfile(profile.ID)
		if len(messages) == 0 {
			unmessaged = append(unmessaged, profile)
		}
	}

	m.log.Info("Found unmessaged connections", "count", len(unmessaged))

	if len(unmessaged) == 0 {
		return nil
	}

	// Send follow-up messages
	return m.SendBulkMessages(unmessaged, "follow_up")
}

// AddTemplate adds a custom message template
func (m *Messenger) AddTemplate(name, content string) {
	m.templates[name] = content
	m.log.Info("Added template", "name", name)
}

// GetTemplate retrieves a template by name
func (m *Messenger) GetTemplate(name string) (string, error) {
	template, exists := m.templates[name]
	if !exists {
		return "", fmt.Errorf("template not found: %s", name)
	}
	return template, nil
}

// ListTemplates returns all available template names
func (m *Messenger) ListTemplates() []string {
	names := make([]string, 0, len(m.templates))
	for name := range m.templates {
		names = append(names, name)
	}
	return names
}

// GetMessageHistory returns all messages for a profile
func (m *Messenger) GetMessageHistory(profileID string) []*storage.Message {
	return m.storage.GetMessagesByProfile(profileID)
}

// CanSendMore checks if more messages can be sent today
func (m *Messenger) CanSendMore() bool {
	today := m.storage.GetActionCountToday("message")
	return today < m.limits.MessagesPerDay
}

// GetStats returns messaging statistics
func (m *Messenger) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"messages_today":   m.storage.GetActionCountToday("message"),
		"limit_daily":      m.limits.MessagesPerDay,
		"can_send_more":    m.CanSendMore(),
		"templates_loaded": len(m.templates),
	}
}
