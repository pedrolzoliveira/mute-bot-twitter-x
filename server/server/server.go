package server

import (
	"encoding/json"
	"mutebotx/database"
	"net/http"
)

type getMutedBotsResponse []string

func getMutedBots(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

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
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
	w.Header().Set("Content-Type", "application/json")

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

func CreateServer() *http.ServeMux {
	server := http.NewServeMux()
	server.HandleFunc("GET /muted-accounts", getMutedBots)
	server.HandleFunc("/mute", muteBot)
	return server
}
