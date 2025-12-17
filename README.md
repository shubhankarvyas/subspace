# Subspace - Advanced Browser Automation PoC

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-Educational-yellow.svg)](LICENSE)
[![Status](https://img.shields.io/badge/Status-PoC-orange.svg)]()

## ⚠️ CRITICAL DISCLAIMER

**THIS IS AN EDUCATIONAL PROOF-OF-CONCEPT ONLY**

This project demonstrates advanced browser automation and anti-detection engineering principles for **EDUCATIONAL and TECHNICAL EVALUATION purposes ONLY**.

### What This Is

- ✅ Educational demonstration of automation engineering
- ✅ Showcase of stealth and behavior simulation techniques
- ✅ Example of clean, modular Go architecture
- ✅ Reference implementation for automation engineers

### What This Is NOT

- ❌ Production-ready automation software
- ❌ Tool for violating platform Terms of Service
- ❌ Working LinkedIn (or any platform) automation
- ❌ Complete or optimized for real-world usage

### Ethical Boundaries

- **NO real selectors** are provided for any platform
- **NO production-grade detection bypasses** are implemented
- **NO encouragement** to violate any platform's ToS
- **CLEAR limitations** documented throughout

This software is deliberately incomplete to prevent misuse. The goal is to demonstrate engineering capabilities, not to enable ToS violations.

---

## 📋 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Stealth Techniques](#-stealth-techniques)
- [Installation](#-installation)
- [Configuration](#-configuration)
- [Usage](#-usage)
- [Demo Mode](#-demo-mode)
- [Technical Details](#-technical-details)
- [Limitations](#-limitations)
- [Contributing](#-contributing)
- [License](#-license)

---

## ✨ Features

### Core Automation

- **Authentication**: Session management, cookie persistence, checkpoint detection
- **Search & Discovery**: Mock search logic, pagination, deduplication
- **Connection Management**: State machine with rate limiting
- **Messaging**: Template-driven messages with personalization
- **Storage**: JSON-based persistence, resumable workflows

### Stealth & Anti-Detection (8+ Techniques)

1. **Bézier Curve Mouse Movement** - Natural mouse paths
2. **Randomized Timing Patterns** - Variable delays and jitter
3. **Browser Fingerprint Masking** - Hide automation indicators
4. **Random Scrolling with Acceleration** - Physics-based scrolling
5. **Typing Simulation with Typos** - Human-like typing errors
6. **Mouse Hover Wandering** - Simulated reading behavior
7. **Activity Scheduling** - Business hours enforcement
8. **Rate Limiting & Cooldown** - Prevent suspicious patterns

### Engineering Excellence

- Clean, modular architecture
- Interface-based design (no direct Rod calls in business logic)
- Structured JSON logging
- Comprehensive error handling
- Resume-safe state management
- Extensive configuration options

---

## 🏗️ Architecture

```
subspace/
├── cmd/
│   └── app/
│       └── main.go              # Application entry point
│
├── internal/
│   ├── auth/                    # Authentication & session management
│   │   └── auth.go              # Login, cookie persistence, checkpoint detection
│   │
│   ├── browser/                 # Browser abstraction layer
│   │   └── browser.go           # Rod wrapper with clean interface
│   │
│   ├── config/                  # Configuration management
│   │   └── config.go            # YAML loading, validation, defaults
│   │
│   ├── connect/                 # Connection request management
│   │   └── connect.go           # State machine, rate limiting
│   │
│   ├── logger/                  # Structured logging
│   │   └── logger.go            # JSON logging, context support
│   │
│   ├── messaging/               # Message sending
│   │   └── messaging.go         # Templates, personalization, history
│   │
│   ├── search/                  # Search & discovery
│   │   └── search.go            # Mock search, pagination, parsing
│   │
│   ├── stealth/                 # Anti-detection engine
│   │   └── stealth.go           # 8+ stealth techniques
│   │
│   └── storage/                 # Data persistence
│       └── storage.go           # JSON storage, state tracking
│
├── config.yaml                  # Main configuration file
├── .env.example                 # Environment variables template
├── go.mod                       # Go module definition
└── README.md                    # This file
```

### Design Principles

1. **Separation of Concerns**: Each package has a single, clear responsibility
2. **Interface Abstraction**: Browser package prevents direct Rod usage
3. **Stealth-First Design**: Human behavior simulation is a first-class system
4. **Configuration-Driven**: All behavior tunable via config
5. **Resume-Safe**: State persists, automation can recover from interruptions

---

## 🎭 Stealth Techniques

### 1. Bézier Curve Mouse Movement

**Why**: Instant mouse teleportation is a strong automation signal.

**How**: Uses cubic Bézier curves to generate smooth, natural mouse paths between points.

**Implementation**:
```go
// Calculate point on curve using control points
x, y := cubicBezier(start, cp1, cp2, end, t)
```

**Configuration**:
```yaml
stealth:
  mouse_speed: 300.0  # Pixels per second
```

**Tradeoff**: Slower than instant movement, but far more human-like.

---

### 2. Randomized Timing Patterns

**Why**: Fixed delays are easily detectable. Humans have variable reaction times.

**How**: Add jitter to all delays, simulate "thinking" with longer pauses.

**Implementation**:
```go
delay := randomInt(config.ActionDelayMin, config.ActionDelayMax)
time.Sleep(time.Duration(delay) * time.Millisecond)
```

**Configuration**:
```yaml
stealth:
  action_delay_min: 500    # Quick actions
  action_delay_max: 2000
  think_time_min: 2000     # Longer pauses
  think_time_max: 5000
```

**Tradeoff**: Makes automation slower but more realistic.

---

### 3. Browser Fingerprint Masking

**Why**: Automation tools leave fingerprints (`navigator.webdriver`, etc.).

**How**: Inject JavaScript to hide/modify detection indicators.

**Implementation**:
```javascript
Object.defineProperty(navigator, 'webdriver', {
    get: () => undefined
});
```

**Configuration**:
```yaml
stealth:
  mask_webdriver: true
  mask_chrome: true
  random_viewport: true
```

**Tradeoff**: Basic defense; advanced detection may still identify automation.

---

### 4. Random Scrolling with Acceleration

**Why**: Humans scroll to view content; bots often don't.

**How**: Random scroll movements with acceleration/deceleration physics.

**Implementation**:
```go
// Ease-in-out acceleration curve
acceleration := easeInOutCubic(progress)
```

**Configuration**:
```yaml
stealth:
  scroll_enabled: true
  scroll_chance: 0.3
  scroll_distance: 300
  scroll_acceleration: 0.8
```

**Tradeoff**: Adds time, increases realism significantly.

---

### 5. Typing Simulation with Typos

**Why**: Instant text appearance is unnatural; perfect typing is rare.

**How**: Character-by-character typing with random delays and occasional typos.

**Implementation**:
```go
if rand.Float64() < config.TypoChance {
    typeWrongChar()
    time.Sleep(randomDelay())
    pressBackspace()
}
```

**Configuration**:
```yaml
stealth:
  typing_speed_min: 80
  typing_speed_max: 200
  typo_chance: 0.03
  typo_correction: true
```

**Tradeoff**: Much slower than instant input, but highly realistic.

---

### 6. Mouse Hover Wandering

**Why**: Humans move mouse while reading/thinking; bots keep it still.

**How**: Occasional small, random mouse movements.

**Implementation**:
```go
// Small random offsets simulating reading
offsetX := randomFloat(-30, 30)
offsetY := randomFloat(-30, 30)
```

**Configuration**:
```yaml
stealth:
  mouse_wander_enabled: true
  mouse_wander_chance: 0.15
```

**Tradeoff**: Minimal performance impact, significant realism boost.

---

### 7. Activity Scheduling (Business Hours & Breaks)

**Why**: Activity at 3 AM is suspicious; humans work during business hours.

**How**: Enforce time-of-day restrictions and break periods.

**Implementation**:
```go
inBusinessHours := isTimeInRange(now, start, end)
inBreakTime := isTimeInRange(now, breakStart, breakEnd)
```

**Configuration**:
```yaml
stealth:
  business_hours_enabled: true
  business_hours_start: "09:00"
  business_hours_end: "17:00"
  break_time_start: "12:00"
  break_time_end: "13:00"
```

**Tradeoff**: Limits automation windows, dramatically reduces suspicion.

---

### 8. Rate Limiting & Cooldown

**Why**: Rapid-fire actions are robotic; humans have natural pacing.

**How**: Enforce minimum time between actions, with cooldown periods.

**Implementation**:
```go
if elapsed < requiredCooldown {
    time.Sleep(requiredCooldown - elapsed)
}
```

**Configuration**:
```yaml
limits:
  connections_per_day: 50
  connections_per_hour: 10
  cooldown_minutes: 60
```

**Tradeoff**: Slower throughput, essential for avoiding detection.

---

## 📦 Installation

### Prerequisites

- **Go 1.21+** (required)
- **Google Chrome** or **Chromium** (for Rod)
- **Git** (for cloning)

### Steps

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd subspace
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Copy environment template**:
   ```bash
   cp .env.example .env
   ```

4. **Build the application**:
   ```bash
   go build -o subspace cmd/app/main.go
   ```

---

## ⚙️ Configuration

### Configuration Files

1. **config.yaml** - Main configuration
   - Stealth technique parameters
   - Rate limits and safety boundaries
   - Browser settings
   - Search and messaging options

2. **.env** - Environment variables (sensitive data)
   - Credentials (mock for PoC)
   - Session paths
   - Log levels

### Key Configuration Options

#### Stealth Tuning

```yaml
stealth:
  # Make automation faster (less stealthy)
  mouse_speed: 500.0
  action_delay_min: 200
  action_delay_max: 800
  
  # Make automation slower (more stealthy)
  mouse_speed: 250.0
  action_delay_min: 1000
  action_delay_max: 3000
```

#### Rate Limits

```yaml
limits:
  connections_per_day: 50     # Conservative: 20-30, Aggressive: 80-100
  connections_per_hour: 10    # Adjust based on account age
  messages_per_day: 30
  searches_per_day: 20
```

#### Business Hours

```yaml
stealth:
  business_hours_enabled: true
  business_hours_start: "09:00"
  business_hours_end: "17:00"
  
  # Disable for 24/7 operation (less stealthy)
  business_hours_enabled: false
```

---

## 🚀 Usage

### Normal Mode (Mock Automation)

Run the full automation workflow with mock data:

```bash
./subspace
```

or

```bash
go run cmd/app/main.go
```

This will:
1. Load configuration
2. Initialize browser and stealth engine
3. Mock login (no real credentials needed)
4. Simulate search
5. Simulate connection requests
6. Simulate messaging
7. Display statistics

### Demo Mode (Stealth Showcase)

See all stealth techniques in action:

```bash
./subspace -demo
```

This demonstrates:
- Bézier mouse movement
- Human-like typing with typos
- Random scrolling
- Mouse wandering
- Timing variations
- Business hours checking
- Rate limiting

### Statistics Mode

View current automation statistics:

```bash
./subspace -stats
```

Shows:
- Profile counts by state
- Activity counters
- Rate limit status

### Custom Configuration

Use a different config file:

```bash
./subspace -config=custom-config.yaml
```

---

## 🎬 Demo Mode

Demo mode showcases each stealth technique individually with visual feedback:

```bash
go run cmd/app/main.go -demo
```

### What You'll See

1. **Mouse Movement**: Watch the cursor move along curved paths
2. **Typing Simulation**: See character-by-character typing with pauses
3. **Scrolling**: Observe smooth, accelerated scrolling
4. **Wandering**: Notice subtle mouse movements
5. **Timing**: See varied delays between actions
6. **Business Hours**: Check current time restrictions
7. **Fingerprinting**: Confirmation of masking techniques
8. **Cooldowns**: Experience rate limiting enforcement

### Recording a Demo Video

For showcasing the PoC:

1. Run in non-headless mode: Set `headless: false` in config.yaml
2. Start screen recording
3. Run demo mode: `./subspace -demo`
4. Show browser window and terminal output side-by-side
5. Highlight:
   - Smooth mouse movements
   - Typing with variable speed
   - Random scrolling
   - Log output showing technique execution

---

## 🔧 Technical Details

### State Machine (Connection Management)

```
discovered → requested → accepted → cooled_down
         ↓
      rejected
```

- **discovered**: Profile found in search
- **requested**: Connection request sent
- **accepted**: Request accepted by target
- **cooled_down**: Follow-up complete, in cooldown
- **rejected**: Request declined or withdrawn

### Storage Format

Data is persisted in JSON format:

```json
{
  "profiles": {
    "profile-123": {
      "id": "profile-123",
      "name": "John Doe",
      "state": "requested",
      "discovered_at": "2024-01-15T10:00:00Z",
      "requested_at": "2024-01-15T10:30:00Z"
    }
  },
  "messages": {},
  "action_logs": []
}
```

### Logging Format

Structured JSON logs for easy parsing:

```json
{
  "timestamp": "2024-01-15T10:00:00Z",
  "level": "INFO",
  "message": "Moving mouse with Bézier curve",
  "fields": {
    "module": "stealth",
    "to_x": 800,
    "to_y": 600
  }
}
```

### Performance Characteristics

- **Mouse Movement**: 10-100 steps, ~2-5 seconds per movement
- **Typing**: 80-200ms per character, ~10-20s per sentence
- **Scrolling**: 10 steps with easing, ~200ms total
- **Page Navigation**: 2-5 seconds with delays

---

## ⚠️ Limitations

### By Design (Ethical)

1. **No Real Selectors**: All DOM selectors are mocked or commented
2. **No Platform-Specific Code**: No optimizations for any real platform
3. **Incomplete Flows**: Critical steps are simulated, not implemented
4. **Mock Data**: Profiles, results, and responses are generated

### Technical

1. **Single Account**: No multi-account support
2. **No CAPTCHA Solving**: Deliberately excluded
3. **No Proxy Rotation**: Not implemented
4. **Limited Error Recovery**: Basic retry logic only
5. **No Advanced Targeting**: Simple keyword search only

### Operational

1. **Not Production-Ready**: Missing error handling for edge cases
2. **No Monitoring**: No metrics, alerts, or health checks
3. **No Scalability**: Single-threaded, local execution only
4. **No Updates**: Static implementation, no self-updating

---

## 🔬 Educational Value

### For Automation Engineers

- Learn modular architecture for automation projects
- Understand anti-detection technique implementation
- See practical examples of behavior simulation
- Study state management and persistence patterns

### For Security Researchers

- Understand what automation looks like from the inside
- Learn detection surface areas and signals
- Study common anti-detection approaches
- Evaluate effectiveness of various techniques

### For Software Architects

- See clean separation of concerns in practice
- Learn interface-based design patterns
- Study configuration-driven architecture
- Understand resume-safe workflow design

---

## 🛠️ Extending the PoC

### Adding New Stealth Techniques

1. Add configuration to `internal/config/config.go`:
   ```go
   NewTechniqueEnabled bool `yaml:"new_technique_enabled"`
   ```

2. Implement in `internal/stealth/stealth.go`:
   ```go
   func (s *Stealth) NewTechnique() error {
       // Implementation
   }
   ```

3. Document in README with Why/How/Tradeoff

### Adding New Modules

1. Create package in `internal/`
2. Define interface for browser abstraction
3. Integrate with stealth engine
4. Add configuration support
5. Update main workflow

---

## 📚 Further Reading

### Browser Automation

- [Rod Documentation](https://go-rod.github.io/)
- [Chrome DevTools Protocol](https://chromedevtools.github.io/devtools-protocol/)
- [Puppeteer Stealth](https://github.com/berstend/puppeteer-extra/tree/master/packages/puppeteer-extra-plugin-stealth)

### Anti-Detection

- [FingerprintJS](https://github.com/fingerprintjs/fingerprintjs)
- [Bot Detection Techniques](https://www.imperva.com/learn/application-security/bot-detection/)
- [Browser Fingerprinting](https://browserleaks.com/)

### Behavior Simulation

- [Bézier Curves](https://en.wikipedia.org/wiki/B%C3%A9zier_curve)
- [Human-Computer Interaction](https://www.interaction-design.org/literature/topics/human-computer-interaction)
- [Keystroke Dynamics](https://en.wikipedia.org/wiki/Keystroke_dynamics)

---

## 🤝 Contributing

This is an educational project. Contributions should focus on:

- ✅ Improving documentation
- ✅ Adding new stealth technique examples
- ✅ Enhancing code architecture
- ✅ Fixing bugs in demo mode

**NOT accepted**:
- ❌ Real selectors for any platform
- ❌ Production-grade detection bypasses
- ❌ Removal of ethical limitations
- ❌ Encouragement of ToS violations

---

## 📄 License

This project is provided for **EDUCATIONAL PURPOSES ONLY**.

By using this software, you agree:

1. To use it solely for learning and technical evaluation
2. NOT to use it to violate any platform's Terms of Service
3. NOT to use it for commercial purposes without explicit permission
4. To assume full responsibility for any consequences of use

The authors and contributors accept NO liability for misuse.

---

## 🙏 Acknowledgments

- **Rod Library**: For excellent Go browser automation
- **Go Community**: For robust tooling and libraries
- **Automation Engineers**: For sharing knowledge and techniques
- **Security Researchers**: For advancing detection understanding

---

## 📧 Contact

For questions about this educational project:

- **Purpose**: Educational demonstration only
- **Support**: None provided (PoC/reference implementation)
- **Issues**: GitHub Issues for documentation/architecture questions only

---

## ⚖️ Final Reminder

**This is an EDUCATIONAL tool demonstrating engineering principles.**

**Do NOT use for:**
- Production automation
- Terms of Service violations
- Unauthorized access or actions
- Any activity without explicit platform permission

**Use responsibly and ethically.**

The goal is to learn about automation engineering and anti-detection techniques, not to misuse them.

**When in doubt: DON'T AUTOMATE IT.**

---

*Built with Go 🐹 for educational purposes 📚*

*Last Updated: December 2024*
