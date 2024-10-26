package domain

import "github.com/olahol/melody"

type MSession struct {
	SessionID string
	Melody    *melody.Session
}

func NewMSession(sessionID string, melody *melody.Session) *MSession {
	return &MSession{
		SessionID: sessionID,
		Melody:    melody,
	}
}
