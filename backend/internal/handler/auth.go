package handler

import (
	"context"
	"fmt"

	"github.com/NapasinNgam/spotify-diary/internal/config"
	"github.com/NapasinNgam/spotify-diary/internal/middleware"
	"github.com/NapasinNgam/spotify-diary/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	cfg       *config.Config
	auth      *spotifyauth.Authenticator
	userRepo  *repository.UserRepository
}

func NewAuthHandler(cfg *config.Config, userRepo *repository.UserRepository) *AuthHandler {
	auth := spotifyauth.New(
		spotifyauth.WithClientID(cfg.SpotifyClientID),
		spotifyauth.WithClientSecret(cfg.SpotifyClientSecret),
		spotifyauth.WithRedirectURL(cfg.SpotifyRedirectURI),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadRecentlyPlayed,
			spotifyauth.ScopeUserReadEmail,
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopeUserTopRead,
			spotifyauth.ScopePlaylistReadPrivate,
			spotifyauth.ScopePlaylistReadCollaborative,
		),
	)

	return &AuthHandler{
		cfg:      cfg,
		auth:     auth,
		userRepo: userRepo,
	}
}

// Login redirects user to Spotify authorization page
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	state := "spotify-diary-state" // TODO: use random state + session store
	url := h.auth.AuthURL(state)
	return c.JSON(fiber.Map{
		"url": url,
	})
}

// Callback handles Spotify OAuth callback
func (h *AuthHandler) Callback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing authorization code",
		})
	}

	// Exchange code for token
	token, err := h.auth.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to exchange token: %v", err),
		})
	}

	// Create Spotify client and get user profile
	httpClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	client := spotify.New(httpClient)

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to get user profile: %v", err),
		})
	}

	// Get avatar URL
	avatarURL := ""
	if len(user.Images) > 0 {
		avatarURL = user.Images[0].URL
	}

	// Upsert user in database
	dbUser, err := h.userRepo.UpsertUser(context.Background(), repository.UpsertUserParams{
		SpotifyID:      string(user.ID),
		DisplayName:    user.DisplayName,
		Email:          user.Email,
		AvatarURL:      avatarURL,
		AccessToken:    token.AccessToken,
		RefreshToken:   token.RefreshToken,
		TokenExpiresAt: token.Expiry,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to save user: %v", err),
		})
	}

	// Generate JWT for frontend
	jwtToken, err := middleware.GenerateJWT(
		h.cfg.JWTSecret,
		dbUser.ID,
		dbUser.SpotifyID,
		dbUser.Email,
		h.cfg.JWTExpiryHours,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Redirect to frontend with token
	redirectURL := fmt.Sprintf("%s/auth/callback?token=%s", h.cfg.FrontendURL, jwtToken)
	return c.Redirect(redirectURL, fiber.StatusTemporaryRedirect)
}

// Me returns the current authenticated user
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	user, err := h.userRepo.GetByID(context.Background(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"user": user,
	})
}
