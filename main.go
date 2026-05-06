package main

import (
	"context"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/razobeckett/spottui/internal/auth"
	"github.com/razobeckett/spottui/internal/config"
	"github.com/razobeckett/spottui/internal/tui"
)

func main() {
	// Load configuration from config.json
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Check if client ID is set (from config or environment)
	if cfg.ClientID == "" {
		fmt.Println("Error: Spotify client ID not configured")
		fmt.Println()
		fmt.Println("To use SpotTUI, you need to:")
		fmt.Println("1. Create an app at https://developer.spotify.com/dashboard")
		fmt.Println("2. Set the redirect URI to: http://localhost:8080/callback")
		fmt.Println("3. Add your client ID to ~/.config/spottui/config.json:")
		fmt.Println(`   {"client_id": "your_client_id_here"}`)
		fmt.Println()
		fmt.Println("Or set the SPOTIFY_CLIENT environment variable:")
		fmt.Println("   SPOTIFY_CLIENT=your_client_id ./spottui")
		os.Exit(1)
	}

	fmt.Println("🎵 SpotTUI - Spotify Terminal Client")
	fmt.Println()

	// Authenticate
	authenticator, err := auth.NewAuthenticator(cfg.ClientID, cfg.CallbackURL)
	if err != nil {
		fmt.Printf("Failed to create authenticator: %v\n", err)
		os.Exit(1)
	}

	client, err := authenticator.GetClient(context.Background())
	if err != nil {
		fmt.Printf("Authentication failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Authenticated successfully!")
	fmt.Println("Starting TUI...")

	// Create and run TUI
	model := tui.NewModel(client, &cfg)
	p := tea.NewProgram(model)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running app: %v\n", err)
		os.Exit(1)
	}
}
