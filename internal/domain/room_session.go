package domain

import "github.com/google/uuid"

type RoomSession struct {
	ID        string `json:"id"`
	RoomID    string `json:"roomId"`
	SessionID string `json:"sessionId"`
}

func NewRoomSession(roomID string, sessionID string) *RoomSession {
	return &RoomSession{
		ID:        uuid.NewString(),
		RoomID:    roomID,
		SessionID: sessionID,
	}
}
