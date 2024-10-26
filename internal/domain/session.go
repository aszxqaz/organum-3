package domain

import "github.com/google/uuid"

type Session struct {
	ID string `json:"id"`
}

func NewSession() *Session {
	return &Session{
		ID: uuid.NewString(),
	}
}

type Token struct {
	ID        string `json:"id"`
	SessionID string `json:"sessionId"`
	Token     string `json:"token"`
}

func NewToken(sessionID string, token string) *Token {
	return &Token{
		ID:        uuid.NewString(),
		SessionID: sessionID,
		Token:     token,
	}
}
