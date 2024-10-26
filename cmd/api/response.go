package main

import "organum/internal/domain"

type CreateSessionResponse struct {
	SessionID string `json:"sessionId"`
	Token     string `json:"token"`
}

type JoinRoomResponse struct {
	Room *domain.Room `json:"room"`
}

type CreateModelResponse struct {
	Checksum string `json:"checksum"`
	RoomID   string `json:"roomId"`
}
