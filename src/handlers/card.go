package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"quizzler/database"
	"quizzler/middleware"
	"quizzler/models"
)

func GetCards(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	deckID, err := strconv.Atoi(r.PathValue("deckId"))
	if err != nil {
		http.Error(w, `{"error": "Invalid deck ID"}`, http.StatusBadRequest)
		return
	}

	// Verify deck belongs to user or is public
	var deckUserID int
	var isPublic bool
	err = database.DB.QueryRow("SELECT user_id, public FROM decks WHERE id = ?", deckID).Scan(&deckUserID, &isPublic)
	if err != nil || (deckUserID != userID && !isPublic) {
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

func GetCard(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	cardID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, `{"error": "Invalid card ID"}`, http.StatusBadRequest)
		return
	}

	var card models.Card
	var deckUserID int
	err = database.DB.QueryRow(`
		SELECT c.id, c.deck_id, c.front, c.back, c.created_at, c.updated_at, d.user_id
		FROM cards c
		JOIN decks d ON c.deck_id = d.id
		WHERE c.id = ?
	`, cardID).Scan(&card.ID, &card.DeckID, &card.Front, &card.Back, &card.CreatedAt, &card.UpdatedAt, &deckUserID)
	if err != nil || deckUserID != userID {
		http.Error(w, `{"error": "Card not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(card)
}

func CreateCard(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	deckID, err := strconv.Atoi(r.PathValue("deckId"))
	if err != nil {
		http.Error(w, `{"error": "Invalid deck ID"}`, http.StatusBadRequest)
		return
	}

	// Verify deck belongs to user
	var deckUserID int
	err = database.DB.QueryRow("SELECT user_id FROM decks WHERE id = ?", deckID).Scan(&deckUserID)
	if err != nil || deckUserID != userID {
		http.Error(w, `{"error": "Deck not found"}`, http.StatusNotFound)
		return
	}

	var req models.CreateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Front == "" || req.Back == "" {
		http.Error(w, `{"error": "Front and back are required"}`, http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec("INSERT INTO cards (deck_id, front, back) VALUES (?, ?, ?)", deckID, req.Front, req.Back)
	if err != nil {
		http.Error(w, `{"error": "Failed to create card"}`, http.StatusInternalServerError)
		return
	}

	cardID, _ := result.LastInsertId()

	var card models.Card
	database.DB.QueryRow("SELECT id, deck_id, front, back, created_at, updated_at FROM cards WHERE id = ?", cardID).
		Scan(&card.ID, &card.DeckID, &card.Front, &card.Back, &card.CreatedAt, &card.UpdatedAt)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(card)
}

func UpdateCard(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	cardID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, `{"error": "Invalid card ID"}`, http.StatusBadRequest)
		return
	}

	// Verify card belongs to user's deck
	var deckUserID int
	err = database.DB.QueryRow(`
		SELECT d.user_id FROM cards c
		JOIN decks d ON c.deck_id = d.id
		WHERE c.id = ?
	`, cardID).Scan(&deckUserID)
	if err != nil || deckUserID != userID {
		http.Error(w, `{"error": "Card not found"}`, http.StatusNotFound)
		return
	}

	var req models.UpdateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Front == "" || req.Back == "" {
		http.Error(w, `{"error": "Front and back are required"}`, http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("UPDATE cards SET front = ?, back = ? WHERE id = ?", req.Front, req.Back, cardID)
	if err != nil {
		http.Error(w, `{"error": "Failed to update card"}`, http.StatusInternalServerError)
		return
	}

	var card models.Card
	database.DB.QueryRow("SELECT id, deck_id, front, back, created_at, updated_at FROM cards WHERE id = ?", cardID).
		Scan(&card.ID, &card.DeckID, &card.Front, &card.Back, &card.CreatedAt, &card.UpdatedAt)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(card)
}

func DeleteCard(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	cardID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, `{"error": "Invalid card ID"}`, http.StatusBadRequest)
		return
	}

	// Verify card belongs to user's deck
	var deckUserID int
	err = database.DB.QueryRow(`
		SELECT d.user_id FROM cards c
		JOIN decks d ON c.deck_id = d.id
		WHERE c.id = ?
	`, cardID).Scan(&deckUserID)
	if err != nil || deckUserID != userID {
		http.Error(w, `{"error": "Card not found"}`, http.StatusNotFound)
		return
	}

	_, err = database.DB.Exec("DELETE FROM cards WHERE id = ?", cardID)
	if err != nil {
		http.Error(w, `{"error": "Failed to delete card"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func ImportCards(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	deckID, err := strconv.Atoi(r.PathValue("deckId"))
	if err != nil {
		http.Error(w, `{"error": "Invalid deck ID"}`, http.StatusBadRequest)
		return
	}

	// Verify deck belongs to user
	var deckUserID int
	err = database.DB.QueryRow("SELECT user_id FROM decks WHERE id = ?", deckID).Scan(&deckUserID)
	if err != nil || deckUserID != userID {
		http.Error(w, `{"error": "Deck not found"}`, http.StatusNotFound)
		return
	}

	var req models.ImportCardsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if len(req.Cards) == 0 {
		http.Error(w, `{"error": "No cards to import"}`, http.StatusBadRequest)
		return
	}

	imported := 0
	for _, card := range req.Cards {
		if card.Front == "" || card.Back == "" {
			continue
		}
		_, err := database.DB.Exec("INSERT INTO cards (deck_id, front, back) VALUES (?, ?, ?)", deckID, card.Front, card.Back)
		if err == nil {
			imported++
		}
	}

	json.NewEncoder(w).Encode(models.ImportCardsResponse{Imported: imported})
}
