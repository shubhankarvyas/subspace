package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"subspace/internal/auth"
	"subspace/internal/browser"
	"subspace/internal/config"
	"subspace/internal/connect"
	"subspace/internal/logger"
	"subspace/internal/messaging"
	"subspace/internal/search"
	"subspace/internal/stealth"
	"subspace/internal/storage"
)

/*
SUBSPACE AUTOMATION POC - MAIN APPLICATION

⚠️  EDUCATIONAL PROOF-OF-CONCEPT ONLY ⚠️

This application demonstrates advanced browser automation and anti-detection
engineering principles. It is NOT intended for production use or to violate
any platform's terms of service.

PURPOSE:
- Showcase engineering capabilities in automation
- Demonstrate stealth and human behavior simulation
- Illustrate clean software architecture patterns
- Provide educational value for automation engineers

LIMITATIONS:
- No real selectors or working automation
- Mock data and simulated flows
- Educational commentary throughout
- Deliberately incomplete for ethical reasons

See README.md for complete documentation and ethical guidelines.
*/

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	demoMode := flag.Bool("demo", false, "Run in demo mode (shows stealth techniques)")
	statsOnly := flag.Bool("stats", false, "Show statistics and exit")
	flag.Parse()

	// Banner
	printBanner()

	// 1. Load Configuration
	fmt.Println("📋 Loading configuration...")
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Printf("❌ Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. Initialize Logger
	logger.Init(cfg.App.LogLevel)
	logger.Info("Starting Subspace Automation PoC",
		"version", "1.0.0",
		"mode", getMode(*demoMode, *statsOnly))

	// 3. Initialize Storage
	logger.Info("Initializing storage", "path", cfg.App.DataDir)
	db, err := storage.New(cfg.App.DataDir + "/db.json")
	if err != nil {
		logger.Error("Failed to initialize storage", "error", err)
		os.Exit(1)
	}

	// Show stats if requested
	if *statsOnly {
		showStats(db)
		return
	}

	// 4. Initialize Browser
	logger.Info("Initializing browser", "headless", cfg.App.Headless)
	b, err := browser.New(cfg.App)
	if err != nil {
		logger.Error("Failed to initialize browser", "error", err)
		os.Exit(1)
	}
	defer func() {
		logger.Info("Shutting down browser")
		if err := b.Close(); err != nil {
			logger.Error("Error closing browser", "error", err)
		}
	}()

	// 5. Initialize Stealth Engine
	logger.Info("Initializing stealth engine")
	s := stealth.New(cfg.Stealth, b.Page)
	logger.Info(s.Summary())

	// Apply fingerprint masking
	if err := s.MaskFingerprint(); err != nil {
		logger.Warn("Failed to apply fingerprint masking", "error", err)
	}

	// 6. Initialize Modules
	logger.Info("Initializing automation modules")
	authenticator := auth.New(b, s, db)
	searcher := search.New(b, s, db)
	connector := connect.New(b, s, db, cfg.Limits)
	messenger := messaging.New(b, s, db, cfg.Limits)

	// 7. Run Demo or Automation Flow
	if *demoMode {
		runDemo(s, b)
	} else {
		runAutomation(cfg, s, authenticator, searcher, connector, messenger)
	}

	logger.Info("Application shutdown complete")
}

// runAutomation executes the main automation workflow
func runAutomation(
	cfg *config.Config,
	s *stealth.Stealth,
	authenticator *auth.Authenticator,
	searcher *search.Searcher,
	connector *connect.Connector,
	messenger *messaging.Messenger,
) {
	logger.Info("Starting automation workflow")

	// Check Business Hours
	if !s.CheckBusinessHours() {
		logger.Warn("Outside business hours")
		fmt.Println("\n⏰ Current time is outside configured business hours")
		fmt.Println("   Configure business_hours in config.yaml to adjust")
		return
	}

	// Step 1: Authentication
	fmt.Println("\n🔐 Step 1: Authentication")
	logger.Info("Attempting login")
	
	if err := authenticator.Login(); err != nil {
		logger.Error("Login failed", "error", err)
		fmt.Printf("❌ Login failed: %v\n", err)
		fmt.Println("   NOTE: This is expected in PoC mode (no real credentials)")
		// Continue anyway for demo purposes
	} else {
		fmt.Println("✅ Login successful (session restored or mock login)")
	}

	// Small delay between major steps
	s.ThinkingPause()

	// Step 2: Search
	fmt.Println("\n🔍 Step 2: Search & Discovery")
	logger.Info("Running search")
	
	keywords := "Software Engineer"
	if err := searcher.RunSearch(keywords, 2); err != nil {
		logger.Error("Search failed", "error", err)
		fmt.Printf("❌ Search failed: %v\n", err)
	} else {
		fmt.Println("✅ Search completed - profiles discovered")
	}

	s.ThinkingPause()

	// Step 3: Connections
	fmt.Println("\n🤝 Step 3: Connection Requests")
	logger.Info("Processing connections")
	
	if connector.CanSendMore() {
		if err := connector.ProcessDailyConnections(); err != nil {
			logger.Error("Connection processing failed", "error", err)
			fmt.Printf("❌ Connection processing failed: %v\n", err)
		} else {
			fmt.Println("✅ Connection requests processed")
		}
	} else {
		fmt.Println("⚠️  Daily connection limit reached")
	}

	s.ThinkingPause()

	// Step 4: Check for accepted connections
	fmt.Println("\n✉️  Step 4: Check Accepted Connections")
	logger.Info("Checking for acceptances")
	
	if err := connector.CheckAcceptedConnections(); err != nil {
		logger.Error("Acceptance check failed", "error", err)
	} else {
		accepted := connector.GetAcceptedConnections()
		fmt.Printf("✅ Found %d accepted connections\n", len(accepted))
	}

	s.ThinkingPause()

	// Step 5: Messaging
	fmt.Println("\n💬 Step 5: Follow-up Messaging")
	logger.Info("Processing messages")
	
	if messenger.CanSendMore() {
		if err := messenger.ProcessAcceptedConnections(); err != nil {
			logger.Error("Messaging failed", "error", err)
			fmt.Printf("❌ Messaging failed: %v\n", err)
		} else {
			fmt.Println("✅ Follow-up messages sent")
		}
	} else {
		fmt.Println("⚠️  Daily message limit reached")
	}

	// Final Summary
	fmt.Println("\n📊 Workflow Summary")
	connStats := connector.GetStats()
	msgStats := messenger.GetStats()
	
	fmt.Printf("   Connections today: %v/%v\n", 
		connStats["connections_today"], 
		connStats["limit_daily"])
	fmt.Printf("   Messages today: %v/%v\n", 
		msgStats["messages_today"], 
		msgStats["limit_daily"])
	fmt.Printf("   Pending requests: %v\n", 
		connStats["pending_requests"])
	fmt.Printf("   Accepted connections: %v\n", 
		connStats["accepted_connections"])

	logger.Info("Automation cycle complete")

	// Keep browser open briefly in non-headless mode
	if !cfg.App.Headless {
		fmt.Println("\n⏳ Keeping browser open for 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}

// runDemo showcases stealth techniques
func runDemo(s *stealth.Stealth, b *browser.Browser) {
	logger.Info("Running demonstration mode")
	fmt.Println("\n🎭 STEALTH TECHNIQUES DEMONSTRATION\n")

	// Demo 1: Mouse Movement
	fmt.Println("1️⃣  Bézier Curve Mouse Movement")
	fmt.Println("   Moving mouse from (100,100) to (800,600)...")
	s.MoveMouse(800, 600)
	fmt.Println("   ✓ Smooth, curved path demonstrated\n")
	time.Sleep(1 * time.Second)

	// Demo 2: Typing with Typos
	fmt.Println("2️⃣  Human-like Typing Simulation")
	fmt.Println("   Typing: 'Hello, this is a test message'")
	s.TypeHumanLike("demo", "Hello, this is a test message")
	fmt.Println("   ✓ Variable speed + occasional typos demonstrated\n")
	time.Sleep(1 * time.Second)

	// Demo 3: Random Scrolling
	fmt.Println("3️⃣  Natural Scrolling Behavior")
	fmt.Println("   Performing random scroll...")
	s.RandomScroll()
	fmt.Println("   ✓ Accelerated scroll with physics demonstrated\n")
	time.Sleep(1 * time.Second)

	// Demo 4: Mouse Wandering
	fmt.Println("4️⃣  Mouse Hover Wandering")
	fmt.Println("   Simulating reading behavior...")
	s.WanderMouse()
	fmt.Println("   ✓ Random micro-movements demonstrated\n")
	time.Sleep(1 * time.Second)

	// Demo 5: Timing Patterns
	fmt.Println("5️⃣  Randomized Timing")
	fmt.Println("   Action delay...")
	start := time.Now()
	s.RandomDelay()
	fmt.Printf("   ✓ Delayed %dms (randomized)\n\n", time.Since(start).Milliseconds())
	
	fmt.Println("   Thinking pause...")
	start = time.Now()
	s.ThinkingPause()
	fmt.Printf("   ✓ Paused %dms (simulating thought)\n\n", time.Since(start).Milliseconds())

	// Demo 6: Business Hours
	fmt.Println("6️⃣  Business Hours Enforcement")
	if s.CheckBusinessHours() {
		fmt.Println("   ✓ Currently within business hours\n")
	} else {
		fmt.Println("   ⚠️  Currently outside business hours\n")
	}

	// Demo 7: Fingerprint Masking
	fmt.Println("7️⃣  Browser Fingerprint Masking")
	fmt.Println("   Applied WebDriver flag masking")
	fmt.Println("   Applied viewport randomization")
	fmt.Println("   ✓ Fingerprint techniques active\n")

	// Demo 8: Rate Limiting
	fmt.Println("8️⃣  Rate Limiting & Cooldown")
	fmt.Println("   Enforcing 5-second cooldown...")
	start = time.Now()
	s.EnforceCooldown("demo", 5)
	fmt.Printf("   ✓ Cooldown enforced (%dms)\n\n", time.Since(start).Milliseconds())

	fmt.Println("✅ Demo complete! All 8+ stealth techniques showcased.")
	fmt.Println("\nℹ️  Check logs for detailed timing and execution data")
	
	// Keep browser open
	fmt.Println("\n⏳ Keeping browser open for 10 seconds...")
	time.Sleep(10 * time.Second)
}

// showStats displays current statistics
func showStats(db *storage.Storage) {
	fmt.Println("\n📊 AUTOMATION STATISTICS\n")
	
	stats := db.GetStats()
	
	fmt.Println("Profile States:")
	fmt.Printf("  Discovered:  %v\n", stats["discovered"])
	fmt.Printf("  Requested:   %v\n", stats["requested"])
	fmt.Printf("  Accepted:    %v\n", stats["accepted"])
	fmt.Printf("  Cooled Down: %v\n", stats["cooled_down"])
	fmt.Printf("  Rejected:    %v\n", stats["rejected"])
	fmt.Printf("  TOTAL:       %v\n\n", stats["total_profiles"])
	
	fmt.Println("Activity Today:")
	fmt.Printf("  Connections: %v\n", stats["connections_today"])
	fmt.Printf("  Messages:    %v\n", stats["messages_today"])
	fmt.Printf("  Total Msgs:  %v\n\n", stats["total_messages"])
	
	fmt.Println("Recent Activity:")
	fmt.Printf("  Connections (last hour): %v\n", stats["connections_last_hour"])
}

// printBanner displays the application banner
func printBanner() {
	banner := `
╔═══════════════════════════════════════════════════════════════╗
║                                                               ║
║   ███████╗██╗   ██╗██████╗ ███████╗██████╗  █████╗  ██████╗ ║
║   ██╔════╝██║   ██║██╔══██╗██╔════╝██╔══██╗██╔══██╗██╔════╝ ║
║   ███████╗██║   ██║██████╔╝███████╗██████╔╝███████║██║      ║
║   ╚════██║██║   ██║██╔══██╗╚════██║██╔═══╝ ██╔══██║██║      ║
║   ███████║╚██████╔╝██████╔╝███████║██║     ██║  ██║╚██████╗ ║
║   ╚══════╝ ╚═════╝ ╚═════╝ ╚══════╝╚═╝     ╚═╝  ╚═╝ ╚═════╝ ║
║                                                               ║
║              Browser Automation PoC v1.0.0                   ║
║                  EDUCATIONAL USE ONLY                         ║
║                                                               ║
╚═══════════════════════════════════════════════════════════════╝

⚠️  IMPORTANT DISCLAIMER ⚠️

This is an EDUCATIONAL proof-of-concept demonstrating:
  • Advanced browser automation techniques
  • Anti-detection and stealth engineering
  • Human behavior simulation systems
  • Clean software architecture patterns

This software is NOT intended for:
  ✗ Production use or deployment
  ✗ Violating any platform's terms of service
  ✗ Automated account actions on real services
  ✗ Bypassing security measures

The purpose is to showcase ENGINEERING CAPABILITIES for
educational and technical evaluation purposes only.

See README.md for complete documentation and guidelines.
`
	fmt.Println(banner)
}

// getMode returns a description of the current running mode
func getMode(demo, stats bool) string {
	if demo {
		return "demonstration"
	}
	if stats {
		return "statistics"
	}
	return "automation"
}
