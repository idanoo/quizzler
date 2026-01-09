package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"quizzler/database"
	"quizzler/middleware"
	"quizzler/models"
)

func GetDecks(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	rows, err := database.DB.Query(`
		SELECT d.id, d.user_id, d.name, d.description, d.public, d.created_at, d.updated_at,
			   (SELECT COUNT(*) FROM cards WHERE deck_id = d.id) as card_count
		FROM decks d
		WHERE d.user_id = ?
		ORDER BY d.updated_at DESC
	`, userID)
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch decks"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	decks := []models.Deck{}
	for rows.Next() {
		var deck models.Deck
		if err := rows.Scan(&deck.ID, &deck.UserID, &deck.Name, &deck.Description, &deck.Public, &deck.CreatedAt, &deck.UpdatedAt, &deck.CardCount); err != nil {
			slog.Warn("Failed to scan deck", "error", err)
			continue
		}
		decks = append(decks, deck)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(decks)
}

func GetDeck(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	deckID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, `{"error": "Invalid deck ID"}`, http.StatusBadRequest)
		return
	}

	var deck models.Deck
	err = database.DB.QueryRow(`
		SELECT d.id, d.user_id, d.name, d.description, d.public, d.created_at, d.updated_at,
			   (SELECT COUNT(*) FROM cards WHERE deck_id = d.id) as card_count
		FROM decks d
		WHERE d.id = ? AND (d.user_id = ? OR d.public = 1)
	`, deckID, userID).Scan(&deck.ID, &deck.UserID, &deck.Name, &deck.Description, &deck.Public, &deck.CreatedAt, &deck.UpdatedAt, &deck.CardCount)
	if err != nil {
		http.Error(w, `{"error": "Deck not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deck)
}

func CreateDeck(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req models.CreateDeckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, `{"error": "Name is required"}`, http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec("INSERT INTO decks (user_id, name, description, public) VALUES (?, ?, ?, ?)", userID, req.Name, req.Description, req.Public)
	if err != nil {
		http.Error(w, `{"error": "Failed to create deck"}`, http.StatusInternalServerError)
		return
	}

	deckID, _ := result.LastInsertId()

	var deck models.Deck
	database.DB.QueryRow("SELECT id, user_id, name, description, public, created_at, updated_at FROM decks WHERE id = ?", deckID).
		Scan(&deck.ID, &deck.UserID, &deck.Name, &deck.Description, &deck.Public, &deck.CreatedAt, &deck.UpdatedAt)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(deck)
}

func UpdateDeck(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	deckID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, `{"error": "Invalid deck ID"}`, http.StatusBadRequest)
		return
	}

	var req models.UpdateDeckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, `{"error": "Name is required"}`, http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec("UPDATE decks SET name = ?, description = ?, public = ? WHERE id = ? AND user_id = ?", req.Name, req.Description, req.Public, deckID, userID)
	if err != nil {
		http.Error(w, `{"error": "Failed to update deck"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, `{"error": "Deck not found"}`, http.StatusNotFound)
		return
	}

	var deck models.Deck
	database.DB.QueryRow("SELECT id, user_id, name, description, public, created_at, updated_at FROM decks WHERE id = ?", deckID).
		Scan(&deck.ID, &deck.UserID, &deck.Name, &deck.Description, &deck.Public, &deck.CreatedAt, &deck.UpdatedAt)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deck)
}

func DeleteDeck(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	deckID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, `{"error": "Invalid deck ID"}`, http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec("DELETE FROM decks WHERE id = ? AND user_id = ?", deckID, userID)
	if err != nil {
		http.Error(w, `{"error": "Failed to delete deck"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, `{"error": "Deck not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetPublicDecks(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	rows, err := database.DB.Query(`
		SELECT d.id, d.user_id, d.name, d.description, d.public, d.created_at, d.updated_at,
			   (SELECT COUNT(*) FROM cards WHERE deck_id = d.id) as card_count
		FROM decks d
		WHERE d.public = 1 AND d.user_id != ?
		ORDER BY d.updated_at DESC
	`, userID)
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch public decks"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	decks := []models.Deck{}
	for rows.Next() {
		var deck models.Deck
		if err := rows.Scan(&deck.ID, &deck.UserID, &deck.Name, &deck.Description, &deck.Public, &deck.CreatedAt, &deck.UpdatedAt, &deck.CardCount); err != nil {
			continue
		}
		decks = append(decks, deck)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(decks)
}

// GetPublicDecksBrowse returns all public decks for unauthenticated browsing
func GetPublicDecksBrowse(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
		SELECT d.id, d.user_id, d.name, d.description, d.public, d.created_at, d.updated_at,
			   (SELECT COUNT(*) FROM cards WHERE deck_id = d.id) as card_count
		FROM decks d
		WHERE d.public = 1
		ORDER BY d.updated_at DESC
	`)
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch public decks"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	decks := []models.Deck{}
	for rows.Next() {
		var deck models.Deck
		if err := rows.Scan(&deck.ID, &deck.UserID, &deck.Name, &deck.Description, &deck.Public, &deck.CreatedAt, &deck.UpdatedAt, &deck.CardCount); err != nil {
			continue
		}
		decks = append(decks, deck)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(decks)
}

// GetPublicDeck returns a single public deck for unauthenticated access
func GetPublicDeck(w http.ResponseWriter, r *http.Request) {
	deckID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, `{"error": "Invalid deck ID"}`, http.StatusBadRequest)
		return
	}

	var deck models.Deck
	err = database.DB.QueryRow(`
		SELECT d.id, d.user_id, d.name, d.description, d.public, d.created_at, d.updated_at,
			   (SELECT COUNT(*) FROM cards WHERE deck_id = d.id) as card_count
		FROM decks d
		WHERE d.id = ? AND d.public = 1
	`, deckID).Scan(&deck.ID, &deck.UserID, &deck.Name, &deck.Description, &deck.Public, &deck.CreatedAt, &deck.UpdatedAt, &deck.CardCount)
	if err != nil {
		http.Error(w, `{"error": "Deck not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deck)
}

// GetPublicDeckCards returns cards for a public deck without authentication
func GetPublicDeckCards(w http.ResponseWriter, r *http.Request) {
	deckID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, `{"error": "Invalid deck ID"}`, http.StatusBadRequest)
		return
	}

	// Verify deck is public
	var isPublic bool
	err = database.DB.QueryRow("SELECT public FROM decks WHERE id = ?", deckID).Scan(&isPublic)
	if err != nil || !isPublic {
		http.Error(w, `{"error": "Deck not found"}`, http.StatusNotFound)
		return
	}

	rows, err := database.DB.Query("SELECT id, deck_id, front, back, created_at, updated_at FROM cards WHERE deck_id = ? ORDER BY created_at", deckID)
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch cards"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	cards := []models.Card{}
	for rows.Next() {
		var card models.Card
		if err := rows.Scan(&card.ID, &card.DeckID, &card.Front, &card.Back, &card.CreatedAt, &card.UpdatedAt); err != nil {
			continue
		}
		cards = append(cards, card)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}
