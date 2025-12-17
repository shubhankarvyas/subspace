package stealth

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	
	"subspace/internal/config"
	"subspace/internal/logger"
)

type Stealth struct {
	config config.StealthConfig
	page   *rod.Page
	log    *logger.ContextLogger
	rng    *rand.Rand
}

// New creates a new stealth engine
func New(cfg config.StealthConfig, page *rod.Page) *Stealth {
	return &Stealth{
		config: cfg,
		page:   page,
		log:    logger.NewContext("stealth"),
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type Point struct {
	X, Y float64
}

// MoveMouse moves the mouse from current position to target using Bézier curves
func (s *Stealth) MoveMouse(toX, toY float64) error {
	s.log.Debug("Moving mouse with Bézier curve", "to_x", toX, "to_y", toY)
	start := time.Now()

	// Get current mouse position (mock for PoC)
	// EDUCATIONAL NOTE: In production, track actual cursor position
	fromX, fromY := s.getCurrentMousePosition()

	// Generate control points for Bézier curve
	cp1, cp2 := s.generateBezierControlPoints(fromX, fromY, toX, toY)

	// Calculate movement steps
	steps := s.calculateSteps(fromX, fromY, toX, toY)
	
	// Move along the curve
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		
		// Calculate point on cubic Bézier curve
		x, y := s.cubicBezier(
			Point{fromX, fromY},
			cp1,
			cp2,
			Point{toX, toY},
			t,
		)

		// EDUCATIONAL NOTE: In production, use:
		// s.page.Mouse.Move(x, y, steps)
		_ = x // Used in production
		_ = y
		
		// Add slight delay between movements
		delay := time.Duration(1000/s.config.MouseSpeed) * time.Millisecond
		time.Sleep(delay)
	}

	logger.Timing("stealth", "move_mouse", start, nil)
	return nil
}

// generateBezierControlPoints creates random control points for natural curves
func (s *Stealth) generateBezierControlPoints(x1, y1, x2, y2 float64) (Point, Point) {
	// Add randomness to control points for variation
	distX := x2 - x1
	distY := y2 - y1
	
	// Control points offset by 30-60% of distance
	cp1 := Point{
		X: x1 + distX*0.3 + s.randomFloat(-50, 50),
		Y: y1 + distY*0.3 + s.randomFloat(-50, 50),
	}
	
	cp2 := Point{
		X: x1 + distX*0.7 + s.randomFloat(-50, 50),
		Y: y1 + distY*0.7 + s.randomFloat(-50, 50),
	}
	
	return cp1, cp2
}

// cubicBezier calculates a point on a cubic Bézier curve
func (s *Stealth) cubicBezier(p0, p1, p2, p3 Point, t float64) (float64, float64) {
	// B(t) = (1-t)³P₀ + 3(1-t)²tP₁ + 3(1-t)t²P₂ + t³P₃
	u := 1 - t
	tt := t * t
	uu := u * u
	uuu := uu * u
	ttt := tt * t

	x := uuu*p0.X + 3*uu*t*p1.X + 3*u*tt*p2.X + ttt*p3.X
	y := uuu*p0.Y + 3*uu*t*p1.Y + 3*u*tt*p2.Y + ttt*p3.Y

	return x, y
}

// calculateSteps determines how many steps needed for smooth movement
func (s *Stealth) calculateSteps(x1, y1, x2, y2 float64) int {
	distance := math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
	// More steps for longer distances
	steps := int(distance / 5)
	if steps < 10 {
		steps = 10
	}
	if steps > 100 {
		steps = 100
	}
	return steps
}

// getCurrentMousePosition returns mock current position
func (s *Stealth) getCurrentMousePosition() (float64, float64) {
	// In production, track actual position
	return 100, 100
}

func (s *Stealth) RandomDelay() {
	delay := s.randomInt(s.config.ActionDelayMin, s.config.ActionDelayMax)
	s.log.Debug("Random delay", "ms", delay)
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

// ThinkingPause simulates a human "thinking" or reading
func (s *Stealth) ThinkingPause() {
	delay := s.randomInt(s.config.ThinkTimeMin, s.config.ThinkTimeMax)
	s.log.Debug("Thinking pause", "ms", delay)
	time.Sleep(time.Duration(delay) * time.Millisecond)
}


func (s *Stealth) MaskFingerprint() error {
	s.log.Info("Applying fingerprint masking")
	
	//  NOTE: The go-rod/stealth package already handles much of this
	// Additional custom masking would be done via JavaScript injection:
	
	if s.config.MaskWebDriver {
		script := `
			// Hide navigator.webdriver
			Object.defineProperty(navigator, 'webdriver', {
				get: () => undefined
			});
		`
		_ = script // In production: s.page.Eval(script)
		s.log.Debug("WebDriver flag masked")
	}

	if s.config.RandomViewport {
		width := s.randomInt(s.config.ViewportWidthMin, s.config.ViewportWidthMax)
		height := s.randomInt(s.config.ViewportHeightMin, s.config.ViewportHeightMax)
		
		//  NOTE: In production:
		// s.page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
		//     Width: width, Height: height,
		// })
		
		s.log.Debug("Viewport randomized", "width", width, "height", height)
	}

	return nil
}


func (s *Stealth) RandomScroll() error {
	if !s.config.ScrollEnabled {
		return nil
	}

	if s.rng.Float64() > s.config.ScrollChance {
		return nil // Don't scroll this time
	}

	s.log.Debug("Performing random scroll")
	
	// Random scroll distance (can be negative for scroll up)
	distance := s.randomInt(-s.config.ScrollDistance, s.config.ScrollDistance*2)
	
	// Simulate scroll with acceleration
	steps := 10
	for i := 0; i < steps; i++ {
		// Ease-in-out acceleration curve
		progress := float64(i) / float64(steps)
		acceleration := s.easeInOutCubic(progress)
		
		stepDistance := float64(distance) * acceleration / float64(steps)
		
		// NOTE: In production:
		// s.page.Mouse.Scroll(0, stepDistance, steps)
		_ = stepDistance // Used in production
		
		time.Sleep(20 * time.Millisecond)
	}

	return nil
}

// easeInOutCubic provides smooth acceleration curve
func (s *Stealth) easeInOutCubic(t float64) float64 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	return 1 - math.Pow(-2*t+2, 3)/2
}

// WHY: Instant text appearance is unnatural; perfect typing is rare.
// HOW: Character-by-character typing with random delays and occasional typos.
// TRADEOFF: Much slower than instant input, but highly realistic.

// TypeHumanLike types text character by character with human-like behavior
func (s *Stealth) TypeHumanLike(selector, text string) error {
	s.log.Debug("Typing with human simulation", "length", len(text))
	start := time.Now()

	for i, char := range text {
		// Check if we should make a typo
		if s.config.TypoChance > 0 && s.rng.Float64() < s.config.TypoChance {
			s.makeTypo(selector)
		}

		// Type the character
		// EDUCATIONAL NOTE: In production:
		// element.Input(string(char))
		
		// Variable delay between keystrokes
		delay := s.randomInt(s.config.TypingSpeedMin, s.config.TypingSpeedMax)
		
		// Longer pause at word boundaries (spaces, commas)
		if char == ' ' || char == ',' || char == '.' {
			delay += s.randomInt(50, 200)
		}
		
		time.Sleep(time.Duration(delay) * time.Millisecond)

		s.log.Debug("Typed character", "index", i, "char", string(char))
	}

	logger.Timing("stealth", "type_human", start, nil)
	return nil
}

// makeTypo simulates a typing error and correction
func (s *Stealth) makeTypo(selector string) {
	if !s.config.TypoCorrection {
		return
	}

	s.log.Debug("Simulating typo")
	
	// Type wrong character
	wrongChar := string(rune(s.randomInt(97, 122))) // Random lowercase letter
	// In production: element.Input(wrongChar)
	_ = wrongChar // Used in production
	
	time.Sleep(time.Duration(s.randomInt(100, 300)) * time.Millisecond)
	
	// "Notice" the error and backspace
	// In production: element.Input("\b")
	
	time.Sleep(time.Duration(s.randomInt(50, 150)) * time.Millisecond)
}

func (s *Stealth) WanderMouse() error {
	if !s.config.MouseWanderEnabled {
		return nil
	}

	if s.rng.Float64() > s.config.MouseWanderChance {
		return nil
	}

	s.log.Debug("Mouse wandering")
	
	// Small random movements (simulate reading or hovering)
	for i := 0; i < s.randomInt(2, 5); i++ {
		offsetX := s.randomFloat(-30, 30)
		offsetY := s.randomFloat(-30, 30)
		
		// Get current position and move slightly
		currentX, currentY := s.getCurrentMousePosition()
		s.MoveMouse(currentX+offsetX, currentY+offsetY)
		
		time.Sleep(time.Duration(s.randomInt(200, 800)) * time.Millisecond)
	}

	return nil
}

func (s *Stealth) CheckBusinessHours() bool {
	if !s.config.BusinessHoursEnabled {
		return true // Always allowed if not enabled
	}

	now := time.Now()
	currentTime := now.Format("15:04")

	// Check if in business hours
	inBusinessHours := s.isTimeInRange(currentTime, s.config.BusinessHoursStart, s.config.BusinessHoursEnd)
	
	// Check if in break time
	inBreakTime := false
	if s.config.BreakTimeEnabled {
		inBreakTime = s.isTimeInRange(currentTime, s.config.BreakTimeStart, s.config.BreakTimeEnd)
	}

	allowed := inBusinessHours && !inBreakTime
	
	if !allowed {
		s.log.Warn("Outside allowed activity hours", 
			"current_time", currentTime,
			"in_business_hours", inBusinessHours,
			"in_break_time", inBreakTime)
	}

	return allowed
}

// isTimeInRange checks if time is between start and end
func (s *Stealth) isTimeInRange(current, start, end string) bool {
	return current >= start && current <= end
}

// WaitForBusinessHours blocks until business hours resume
func (s *Stealth) WaitForBusinessHours() {
	for !s.CheckBusinessHours() {
		s.log.Info("Waiting for business hours to resume...")
		time.Sleep(15 * time.Minute) // Check every 15 minutes
	}
}

var lastActionTime time.Time

// EnforceCooldown ensures minimum time between actions
func (s *Stealth) EnforceCooldown(actionType string, minDelaySeconds int) {
	if lastActionTime.IsZero() {
		lastActionTime = time.Now()
		return
	}

	elapsed := time.Since(lastActionTime)
	required := time.Duration(minDelaySeconds) * time.Second

	if elapsed < required {
		remaining := required - elapsed
		s.log.Info("Enforcing cooldown", 
			"action", actionType,
			"wait_seconds", remaining.Seconds())
		time.Sleep(remaining)
	}

	lastActionTime = time.Now()
}
func (s *Stealth) randomInt(min, max int) int {
	if min >= max {
		return min
	}
	return min + s.rng.Intn(max-min+1)
}

// randomFloat returns a random float64 between min and max
func (s *Stealth) randomFloat(min, max float64) float64 {
	return min + s.rng.Float64()*(max-min)
}

// ShouldProceed checks if action should proceed based on random chance
func (s *Stealth) ShouldProceed(probability float64) bool {
	return s.rng.Float64() < probability
}

// Summary logs a summary of active stealth techniques
func (s *Stealth) Summary() string {
	active := []string{}
	
	if s.config.MouseWanderEnabled {
		active = append(active, "Mouse Wandering")
	}
	if s.config.ScrollEnabled {
		active = append(active, "Random Scrolling")
	}
	if s.config.TypoCorrection {
		active = append(active, "Typo Simulation")
	}
	if s.config.BusinessHoursEnabled {
		active = append(active, "Business Hours")
	}
	if s.config.MaskWebDriver {
		active = append(active, "Fingerprint Masking")
	}
	
	return fmt.Sprintf("Active stealth techniques: %v", active)
}
