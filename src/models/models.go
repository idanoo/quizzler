package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Deck struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Public      bool      `json:"public"`
	CardCount   int       `json:"card_count,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Card struct {
	ID        int       `json:"id"`
	DeckID    int       `json:"deck_id"`
	Front     string    `json:"front"`
	Back      string    `json:"back"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Auth
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// Decks
type CreateDeckRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
}

type UpdateDeckRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
}

// Cards
type CreateCardRequest struct {
	Front string `json:"front"`
	Back  string `json:"back"`
}

type UpdateCardRequest struct {
	Front string `json:"front"`
	Back  string `json:"back"`
}

type ImportCardsRequest struct {
	Cards []CreateCardRequest `json:"cards"`
}

type ImportCardsResponse struct {
	Imported int `json:"imported"`
}
