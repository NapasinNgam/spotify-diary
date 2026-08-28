package model

import "time"

type User struct {
	ID              int       `json:"id"`
	SpotifyID       string    `json:"spotify_id"`
	DisplayName     string    `json:"display_name"`
	Email           string    `json:"email"`
	AvatarURL       string    `json:"avatar_url"`
	AccessToken     string    `json:"-"` // never expose
	RefreshToken    string    `json:"-"` // never expose
	TokenExpiresAt  time.Time `json:"-"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
