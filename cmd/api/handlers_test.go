package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"organum/internal/domain"
	"organum/internal/jsonlog"
	"organum/internal/store"
	"slices"
	"testing"

	"github.com/olahol/melody"
	"github.com/test-go/testify/assert"
)

func TestSessionHandler(t *testing.T) {
	h := getHandler()

	rsp := createSession(h, nil)
	assert.Equal(t, rsp.Code, http.StatusCreated)
	sessions := getSessions(h)
	assert.Len(t, sessions, 1)

	session := mustCreateSession(h, nil)
	assert.NotEmpty(t, session.SessionID)
	assert.NotEmpty(t, session.Token)
	sessions = getSessions(h)
	assert.Len(t, sessions, 2)

	rsp = createSession(h, session)
	assert.Equal(t, http.StatusConflict, rsp.Code)
}

func TestCreateRoom(t *testing.T) {
	h := getHandler()

	rsp := createRoom(h, nil)
	assert.Equal(t, http.StatusUnauthorized, rsp.Code)

	session := mustCreateSession(h, nil)
	rsp = createRoom(h, session)
	assert.Equal(t, http.StatusCreated, rsp.Code)
	rooms := getRooms(h)
	assert.Len(t, rooms, 1)

	rsp = createRoom(h, session)
	assert.Equal(t, http.StatusConflict, rsp.Code)
	rooms = getRooms(h)
	assert.Len(t, rooms, 1)

	session = mustCreateSession(h, nil)
	rsp = createRoom(h, session)
	assert.Equal(t, http.StatusCreated, rsp.Code)
	rooms = getRooms(h)
	assert.Len(t, rooms, 2)

	session = mustCreateSession(h, nil)
	room := mustCreateRoom(h, session)
	assert.Equal(t, session.SessionID, room.OwnerID)
	assert.Equal(t, session.SessionID, room.LockOwnerID)
}

func TestJoinRoom(t *testing.T) {
	h := getHandler()

	csr := mustCreateSession(h, nil)
	room := mustCreateRoom(h, csr)
	rsp := joinRoom(h, room.ID, csr)
	assert.Equal(t, http.StatusConflict, rsp.Code)
	joins := getSessionsForRoom(h, room.ID)
	assert.Len(t, joins, 1)
	assert.Equal(t, csr.SessionID, joins[0].SessionID)

	csr = mustCreateSession(h, nil)
	rsp = joinRoom(h, room.ID, csr)
	assert.Equal(t, http.StatusCreated, rsp.Code)
	joins = getSessionsForRoom(h, room.ID)
	assert.Len(t, joins, 2)
	assert.True(t, slices.ContainsFunc(joins, func(rs *domain.RoomSession) bool { return rs.SessionID == csr.SessionID }))

	rsp = joinRoom(h, room.ID, csr)
	assert.Equal(t, http.StatusConflict, rsp.Code)
	joins = getSessionsForRoom(h, room.ID)
	assert.Len(t, joins, 2)
}

func TestLeaveRoom(t *testing.T) {
	h := getHandler()

	csr := mustCreateSession(h, nil)
	room := mustCreateRoom(h, csr)
	rsp := leaveRoom(h, room.ID, csr)
	assert.Equal(t, http.StatusOK, rsp.Code)
	joins := getSessionsForRoom(h, room.ID)
	assert.Len(t, joins, 0)

	rsp = leaveRoom(h, "WRONG_ID", csr)
	assert.Equal(t, http.StatusNotFound, rsp.Code)

	room = mustCreateRoom(h, csr)
	csr = mustCreateSession(h, nil)
	mustJoinRoom(h, room.ID, csr)
	joins = getSessionsForRoom(h, room.ID)
	assert.Len(t, joins, 2)
	rsp = leaveRoom(h, room.ID, csr)
	assert.Equal(t, http.StatusOK, rsp.Code)
	joins = getSessionsForRoom(h, room.ID)
	assert.Len(t, joins, 1)

	room = mustCreateRoom(h, csr)
	mustLeaveRoom(h, room.ID, csr)
	room = getRoom(h, room.ID)
	assert.Nil(t, room)
}

func TestModels(t *testing.T) {
	h := getHandler()

	csr := mustCreateSession(h, nil)
	rsp := createModel(h, []byte{1, 2, 3}, csr)
	assert.Equal(t, http.StatusBadRequest, rsp.Code)

	room := mustCreateRoom(h, csr)
	rsp = createModel(h, []byte{1, 2, 3}, csr)
	assert.Equal(t, http.StatusCreated, rsp.Code)
	models := getModels(h)
	assert.Len(t, models, 1)
	assert.Equal(t, csr.SessionID, models[0].UploaderID)
	assert.Equal(t, room.ID, models[0].RoomID)
	model := getModel(h, models[0].Checksum)
	assert.Equal(t, csr.SessionID, model.UploaderID)
	assert.Equal(t, room.ID, model.RoomID)

	model = mustCreateModel(h, []byte{1, 2, 3}, csr)
	assert.Equal(t, csr.SessionID, model.UploaderID)
	assert.Equal(t, room.ID, model.RoomID)
	model = getModel(h, model.Checksum)
	assert.Equal(t, csr.SessionID, model.UploaderID)
	assert.Equal(t, room.ID, model.RoomID)
	models = getModels(h)
	assert.Len(t, models, 1)

	model = mustCreateModel(h, []byte{5, 6, 7}, csr)
	assert.Equal(t, csr.SessionID, model.UploaderID)
	assert.Equal(t, room.ID, model.RoomID)
	model = getModel(h, model.Checksum)
	assert.Equal(t, csr.SessionID, model.UploaderID)
	assert.Equal(t, room.ID, model.RoomID)
	models = getModels(h)
	assert.Len(t, models, 2)

	deleteRoom(h, room.ID, csr)
	model = getModel(h, model.Checksum)
	assert.Nil(t, model)
	models = getModels(h)
	assert.Len(t, models, 0)

	room = mustCreateRoom(h, csr)
	model = mustCreateModel(h, []byte{5, 6, 7}, csr)
	leaveRoom(h, room.ID, csr)
	model = getModel(h, model.Checksum)
	assert.Nil(t, model)
	models = getModels(h)
	assert.Len(t, models, 0)
}

func getSessions(h http.Handler) []*domain.Session {
	req, _ := http.NewRequest("GET", "/sessions", nil)
	rsp := execReq(req, h)
	var sessions []*domain.Session
	json.NewDecoder(rsp.Body).Decode(&sessions)
	return sessions
}

func getRoom(h http.Handler, roomID string) *domain.Room {
	req, _ := http.NewRequest("GET", "/rooms/"+roomID, nil)
	rsp := execReq(req, h)
	var room *domain.Room
	json.NewDecoder(rsp.Body).Decode(&room)
	return room
}

func getRooms(h http.Handler) []*domain.Room {
	req, _ := http.NewRequest("GET", "/rooms", nil)
	rsp := execReq(req, h)
	var rooms []*domain.Room
	json.NewDecoder(rsp.Body).Decode(&rooms)
	return rooms
}

func getSessionsForRoom(h http.Handler, roomID string) []*domain.RoomSession {
	req, _ := http.NewRequest("GET", "/rooms/"+roomID+"/sessions", nil)
	rsp := execReq(req, h)
	var roomSessions []*domain.RoomSession
	json.NewDecoder(rsp.Body).Decode(&roomSessions)
	return roomSessions
}

func getModel(h http.Handler, checksum string) *domain.ModelJSON {
	req, _ := http.NewRequest("GET", "/models/"+checksum, nil)
	rsp := execReq(req, h)
	var model *domain.ModelJSON
	json.NewDecoder(rsp.Body).Decode(&model)
	return model
}

func getModels(h http.Handler) []*domain.ModelJSON {
	req, _ := http.NewRequest("GET", "/models", nil)
	rsp := execReq(req, h)
	var models []*domain.ModelJSON
	json.NewDecoder(rsp.Body).Decode(&models)
	return models
}

func leaveRoom(h http.Handler, roomID string, csr *CreateSessionResponse) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("DELETE", "/rooms/"+roomID+"/sessions", nil)
	setAuthHeader(req, csr)
	return execReq(req, h)
}

func mustLeaveRoom(h http.Handler, roomID string, csr *CreateSessionResponse) {
	leaveRoom(h, roomID, csr)
}

func joinRoom(h http.Handler, roomID string, csr *CreateSessionResponse) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", "/rooms/"+roomID+"/sessions", nil)
	setAuthHeader(req, csr)
	return execReq(req, h)
}

func mustJoinRoom(h http.Handler, roomID string, csr *CreateSessionResponse) *domain.Room {
	rsp := joinRoom(h, roomID, csr)
	var room domain.Room
	json.NewDecoder(rsp.Body).Decode(&room)
	return &room
}

func createModel(h http.Handler, file []byte, csr *CreateSessionResponse) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", "/models", bytes.NewBuffer(file))
	setAuthHeader(req, csr)
	return execReq(req, h)
}

func mustCreateModel(h http.Handler, file []byte, csr *CreateSessionResponse) *domain.ModelJSON {
	rsp := createModel(h, file, csr)
	var model *domain.ModelJSON
	json.NewDecoder(rsp.Body).Decode(&model)
	return model
}

func deleteRoom(h http.Handler, roomID string, csr *CreateSessionResponse) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("DELETE", "/rooms/"+roomID, nil)
	setAuthHeader(req, csr)
	return execReq(req, h)
}

func createRoom(h http.Handler, csr *CreateSessionResponse) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", "/rooms", nil)
	setAuthHeader(req, csr)
	return execReq(req, h)
}

func mustCreateRoom(h http.Handler, csr *CreateSessionResponse) *domain.Room {
	rsp := createRoom(h, csr)
	var room *domain.Room
	json.NewDecoder(rsp.Body).Decode(&room)
	return room
}

func createSession(h http.Handler, csr *CreateSessionResponse) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", "/sessions", nil)
	setAuthHeader(req, csr)
	return execReq(req, h)
}

func mustCreateSession(h http.Handler, csr *CreateSessionResponse) *CreateSessionResponse {
	rsp := createSession(h, csr)
	var session CreateSessionResponse
	json.NewDecoder(rsp.Body).Decode(&session)
	return &session
}

func setAuthHeader(r *http.Request, csr *CreateSessionResponse) {
	if csr != nil {
		r.Header.Set("Authorization", "Bearer "+csr.Token)
	}
}

func execReq(req *http.Request, h http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}

func getHandler() http.Handler {
	logger := jsonlog.New(io.Discard, jsonlog.LevelInfo)
	app := &application{
		config: &config{
			secret: "secret",
		},
		logger: logger,
		store:  store.NewStore(),
		melody: melody.New(),
	}
	return app.routes()
}
