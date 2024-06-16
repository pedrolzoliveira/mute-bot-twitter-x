package server

import (
	"encoding/json"
	"mutebotx/database"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func clearMutedBotsTable() {
	db := database.CreateDatabase()
	db.MustExec("DELETE FROM muted_bots")
}

func TestCreateServer(t *testing.T) {
	server := CreateServer()
	assert.NotNil(t, server)
}

func TestGetMutedBotsWithNoMutedBots(t *testing.T) {
	clearMutedBotsTable()

	req := httptest.NewRequest("GET", "/muted-bots", nil)
	w := httptest.NewRecorder()
	getMutedBots(w, req)

	var response []string
	json.NewDecoder(w.Body).Decode(&response)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, []string{}, response)
}

func TestGetMutedBotsWithMutedBots(t *testing.T) {
	clearMutedBotsTable()

	db := database.CreateDatabase()
	db.MustExec("INSERT INTO muted_bots (user_handle) VALUES ('test')")

	req := httptest.NewRequest("GET", "/muted-bots", nil)
	w := httptest.NewRecorder()
	getMutedBots(w, req)

	var response []string
	json.NewDecoder(w.Body).Decode(&response)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, []string{"test"}, response)
}

func TestMuteBot(t *testing.T) {
	clearMutedBotsTable()

	bodyReader := strings.NewReader(`{"user_handle": "any_user_handle"}`)

	req := httptest.NewRequest("POST", "/mute", bodyReader)
	w := httptest.NewRecorder()

	muteBot(w, req)

	var mutedBots []database.MutedBots
	database.CreateDatabase().Select(&mutedBots, "SELECT * FROM muted_bots")

	assert.Equal(t, 201, w.Code)
	assert.Equal(
		t,
		[]database.MutedBots{{UserHandle: "any_user_handle"}},
		mutedBots,
	)
}

func TestMuteBotWithExistingMutedBot(t *testing.T) {
	clearMutedBotsTable()

	db := database.CreateDatabase()
	db.MustExec("INSERT INTO muted_bots (user_handle) VALUES ('any_user_handle')")

	bodyReader := strings.NewReader(`{"user_handle": "any_user_handle"}`)

	req := httptest.NewRequest("POST", "/mute", bodyReader)
	w := httptest.NewRecorder()

	muteBot(w, req)

	var mutedBots []database.MutedBots
	database.CreateDatabase().Select(&mutedBots, "SELECT * FROM muted_bots")

	assert.Equal(t, 201, w.Code)
	assert.Equal(
		t,
		[]database.MutedBots{{UserHandle: "any_user_handle"}},
		mutedBots,
	)

}

func TestMuteBotWithInvalidJSON(t *testing.T) {
	clearMutedBotsTable()

	bodyReader := strings.NewReader(`invalid: "json"}`)

	req := httptest.NewRequest("POST", "/mute", bodyReader)
	w := httptest.NewRecorder()

	muteBot(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestMuteBotWithNoUserHandle(t *testing.T) {
	clearMutedBotsTable()

	bodyReader := strings.NewReader(`{}`)

	req := httptest.NewRequest("POST", "/mute", bodyReader)
	w := httptest.NewRecorder()

	muteBot(w, req)

	assert.Equal(t, 400, w.Code)
}
