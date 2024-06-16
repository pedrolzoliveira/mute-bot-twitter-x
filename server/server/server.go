package server

import (
	"encoding/json"
	"mutebotx/database"
	"net/http"

	"github.com/rs/cors"
)

type getMutedBotsResponse []string

func getMutedBots(w http.ResponseWriter, _ *http.Request) {
	db := database.CreateDatabase()
	var bots []database.MutedBots
	var response getMutedBotsResponse = make([]string, 0)

	err := db.Select(&bots, "SELECT * FROM muted_bots")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

	if bots == nil {
		json.NewEncoder(w).Encode(response)
		return
	}

	for _, bot := range bots {
		response = append(response, bot.UserHandle)
	}

	json.NewEncoder(w).Encode(response)
}

type MuteBotRequest struct {
	UserHandle string `json:"user_handle"`
}

func muteBot(w http.ResponseWriter, r *http.Request) {
	db := database.CreateDatabase()

	var req MuteBotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.UserHandle == "" {
		http.Error(w, "user_handle is required", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT OR IGNORE INTO muted_bots (user_handle) VALUES (?)", req.UserHandle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}

func CreateServer() http.Handler {
	server := http.NewServeMux()

	server.HandleFunc("GET /muted-accounts", getMutedBots)
	server.HandleFunc("POST /mute", muteBot)

	handler := cors.AllowAll().Handler(server)

	return wrap(handler)
}
