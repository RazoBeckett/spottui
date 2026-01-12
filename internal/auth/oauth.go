package auth

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

// Authenticator handles Spotify OAuth with PKCE flow
type Authenticator struct {
	auth         *spotifyauth.Authenticator
	oauth2Config *oauth2.Config
	codeVerifier string
	state        string
	tokenStore   *TokenStore
	redirectURI  string
}

// NewAuthenticator creates a new PKCE authenticator with required scopes
func NewAuthenticator(clientID, redirectURI string) (*Authenticator, error) {
	verifier, err := generateCodeVerifier()
	if err != nil {
		return nil, fmt.Errorf("failed to generate code verifier: %w", err)
	}

	scopes := []string{
		string(spotifyauth.ScopeUserReadPlaybackState),
		string(spotifyauth.ScopeUserModifyPlaybackState),
		string(spotifyauth.ScopeUserReadCurrentlyPlaying),
		string(spotifyauth.ScopeUserLibraryRead),
		string(spotifyauth.ScopeUserLibraryModify),
		string(spotifyauth.ScopePlaylistReadPrivate),
		string(spotifyauth.ScopePlaylistReadCollaborative),
		string(spotifyauth.ScopePlaylistModifyPublic),
		string(spotifyauth.ScopePlaylistModifyPrivate),
		string(spotifyauth.ScopeUserReadPrivate),
	}

	auth := spotifyauth.New(
		spotifyauth.WithClientID(clientID),
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPlaybackState,
			spotifyauth.ScopeUserModifyPlaybackState,
			spotifyauth.ScopeUserReadCurrentlyPlaying,
			spotifyauth.ScopeUserLibraryRead,
			spotifyauth.ScopeUserLibraryModify,
			spotifyauth.ScopePlaylistReadPrivate,
			spotifyauth.ScopePlaylistReadCollaborative,
			spotifyauth.ScopePlaylistModifyPublic,
			spotifyauth.ScopePlaylistModifyPrivate,
			spotifyauth.ScopeUserReadPrivate,
		),
	)

	oauth2Cfg := &oauth2.Config{
		ClientID: clientID,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},
		RedirectURL: redirectURI,
		Scopes:      scopes,
	}

	return &Authenticator{
		auth:         auth,
		oauth2Config: oauth2Cfg,
		codeVerifier: verifier,
		state:        randomString(16),
		tokenStore:   NewTokenStore(),
		redirectURI:  redirectURI,
	}, nil
}

// GetClient returns an authenticated Spotify client
// Tries cached token first, falls back to browser auth
func (a *Authenticator) GetClient(ctx context.Context) (*spotify.Client, error) {
	// Try loading existing token
	token, err := a.tokenStore.Load()
	if err == nil {
		client, err := a.clientFromToken(ctx, token)
		if err == nil {
			return client, nil
		}
		// Token invalid, continue to re-auth
	}

	// Full OAuth flow
	return a.authenticateViaBrowser(ctx)
}

func (a *Authenticator) authenticateViaBrowser(ctx context.Context) (*spotify.Client, error) {
	// Start callback server
	clientChan := make(chan *spotify.Client)
	errChan := make(chan error)
	srv := a.startCallbackServer(clientChan, errChan)
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(shutdownCtx)
	}()

	// Build auth URL with PKCE challenge
	challenge := generateCodeChallenge(a.codeVerifier)
	url := a.auth.AuthURL(a.state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", challenge),
	)

	fmt.Println("Opening browser for Spotify login...")
	fmt.Println("If browser doesn't open, visit:", url)
	openBrowser(url)

	// Wait for callback
	select {
	case client := <-clientChan:
		return client, nil
	case err := <-errChan:
		return nil, err
	case <-time.After(5 * time.Minute):
		return nil, fmt.Errorf("authentication timeout")
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (a *Authenticator) startCallbackServer(clientChan chan *spotify.Client, errChan chan error) *http.Server {
	mux := http.NewServeMux()
	srv := &http.Server{Addr: ":8080", Handler: mux}

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// Validate state
		if r.FormValue("state") != a.state {
			errChan <- fmt.Errorf("state mismatch")
			http.Error(w, "State mismatch", http.StatusBadRequest)
			return
		}

		// Exchange code for token with PKCE verifier
		token, err := a.auth.Token(r.Context(), a.state, r,
			oauth2.SetAuthURLParam("code_verifier", a.codeVerifier))
		if err != nil {
			errChan <- fmt.Errorf("token exchange failed: %w", err)
			http.Error(w, "Authentication failed", http.StatusInternalServerError)
			return
		}

		// Save token for future use
		if err := a.tokenStore.Save(token); err != nil {
			fmt.Printf("Warning: failed to save token: %v\n", err)
		}

		// Create client
		httpClient := a.auth.Client(r.Context(), token)
		client := spotify.New(httpClient)

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<!DOCTYPE html><html><body>
			<h1>Authentication successful!</h1>
			<p>You can close this window and return to the terminal.</p>
			<script>window.close();</script>
		</body></html>`)

		clientChan <- client
	})

	go srv.ListenAndServe()
	return srv
}

func (a *Authenticator) clientFromToken(ctx context.Context, token *oauth2.Token) (*spotify.Client, error) {
	if !token.Valid() && token.RefreshToken == "" {
		return nil, fmt.Errorf("token expired and no refresh token available")
	}

	tokenSource := a.oauth2Config.TokenSource(ctx, token)
	autoSaveSource := &autoSaveTokenSource{
		source:     tokenSource,
		tokenStore: a.tokenStore,
	}

	httpClient := oauth2.NewClient(ctx, autoSaveSource)
	client := spotify.New(httpClient)

	if _, err := client.CurrentUser(ctx); err != nil {
		return nil, err
	}

	return client, nil
}

type autoSaveTokenSource struct {
	source     oauth2.TokenSource
	tokenStore *TokenStore
	lastToken  *oauth2.Token
}

func (s *autoSaveTokenSource) Token() (*oauth2.Token, error) {
	token, err := s.source.Token()
	if err != nil {
		return nil, err
	}

	if s.lastToken == nil || token.AccessToken != s.lastToken.AccessToken {
		s.tokenStore.Save(token)
		s.lastToken = token
	}

	return token, nil
}

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	}

	if cmd != "" {
		exec.Command(cmd, args...).Start()
	}
}
