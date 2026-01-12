package main

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"

	"github.com/razobeckett/spottui/internal/auth"
	"github.com/razobeckett/spottui/internal/tui"
)

func main() {
	// Load .env file if present
	godotenv.Load()

	// Get client ID from environment
	clientID := os.Getenv("SPOTIFY_CLIENT")
	if clientID == "" {
		fmt.Println("Error: SPOTIFY_CLIENT environment variable not set")
		fmt.Println()
		fmt.Println("To use SpotTUI, you need to:")
		fmt.Println("1. Create an app at https://developer.spotify.com/dashboard")
		fmt.Println("2. Set the redirect URI to: http://localhost:8080/callback")
		fmt.Println("3. Create a .env file with:")
		fmt.Println("   SPOTIFY_CLIENT=your_client_id")
		fmt.Println("   SPOTIFY_CALLBACK=http://localhost:8080/callback")
		os.Exit(1)
	}

	// Get callback URL from environment
	callbackURL := os.Getenv("SPOTIFY_CALLBACK")
	if callbackURL == "" {
		callbackURL = "http://localhost:8080/callback"
	}

	// Authenticate
	authenticator, err := auth.NewAuthenticator(clientID, callbackURL)
	if err != nil {
		fmt.Printf("Failed to create authenticator: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🎵 SpotTUI - Spotify Terminal Client")
	fmt.Println()

	client, err := authenticator.GetClient(context.Background())
	if err != nil {
		fmt.Printf("Authentication failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Authenticated successfully!")
	fmt.Println("Starting TUI...")

	// Create and run TUI
	model := tui.NewModel(client)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running app: %v\n", err)
		os.Exit(1)
	}
}
