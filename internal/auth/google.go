package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleUser holds the minimal profile data we need from Google.
type GoogleUser struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type GoogleAuth struct {
	cfg *config.Config
}

func NewGoogle(cfg *config.Config) *GoogleAuth {
	return &GoogleAuth{cfg: cfg}
}

func (g *GoogleAuth) oauthConfig(redirectURI string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     g.cfg.GoogleClientID,
		ClientSecret: g.cfg.GoogleClientSecret,
		RedirectURL:  redirectURI,
		Scopes: []string{
			"openid",
			"email",
			"profile",
		},
		Endpoint: google.Endpoint,
	}
}

func redirectURI(cfg *config.Config, desktop bool) string {
	base := fmt.Sprintf("http://127.0.0.1:%d/auth/callback", cfg.ServerPort)
	if !desktop {
		base = "https://httptoolkitpro.vercel.app/auth/callback"
	}
	return base
}

func (g *GoogleAuth) AuthURL(desktop bool, state string) (string, error) {
	if g.cfg.GoogleClientID == "" || g.cfg.GoogleClientSecret == "" {
		return "", fmt.Errorf("Google OAuth not configured")
	}
	cfg := g.oauthConfig(redirectURI(g.cfg, desktop))
	return cfg.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce), nil
}

func (g *GoogleAuth) Exchange(ctx context.Context, code string, desktop bool) (*GoogleUser, string, string, error) {
	if g.cfg.GoogleClientID == "" || g.cfg.GoogleClientSecret == "" {
		return nil, "", "", fmt.Errorf("Google OAuth not configured")
	}

	cfg := g.oauthConfig(redirectURI(g.cfg, desktop))
	token, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, "", "", fmt.Errorf("oauth exchange: %w", err)
	}

	idToken, _ := token.Extra("id_token").(string)
	refreshToken := token.RefreshToken
	if idToken == "" {
		// Fall back to the access token for userinfo; id_token may not always
		// be present depending on the OAuth client configuration.
		idToken = token.AccessToken
	}

	user, err := userinfo(ctx, token.AccessToken)
	if err != nil {
		return nil, "", "", err
	}
	return user, idToken, refreshToken, nil
}

// userinfo calls Google's userinfo endpoint with the provided access token.
func userinfo(ctx context.Context, accessToken string) (*GoogleUser, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("userinfo request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("userinfo %d: %s", resp.StatusCode, string(body))
	}

	var user GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("userinfo decode: %w", err)
	}
	return &user, nil
}

// VerifyIDToken checks a Google ID token with the tokeninfo endpoint and
// validates that it was issued for this app.
func (g *GoogleAuth) VerifyIDToken(ctx context.Context, idToken string) (*GoogleUser, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	u := "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + url.QueryEscape(idToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tokeninfo request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("tokeninfo %d: %s", resp.StatusCode, string(body))
	}

	var payload struct {
		Email string `json:"email"`
		Aud   string `json:"aud"`
		Name  string `json:"name"`
		Pic   string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("tokeninfo decode: %w", err)
	}

	// The token must have been issued for our OAuth client.
	if !strings.Contains(g.cfg.GoogleClientID, payload.Aud) && payload.Aud != g.cfg.GoogleClientID {
		return nil, fmt.Errorf("token aud mismatch")
	}

	return &GoogleUser{Email: payload.Email, Name: payload.Name, Picture: payload.Pic}, nil
}
