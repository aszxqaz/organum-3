package ws

import (
	"organum/internal/domain"

	"github.com/olahol/melody"
)

func (h *WsHandler) getSessionFromMelodySession(s *melody.Session) *domain.Session {
	val, ok := s.Get("session")
	if !ok {
		return nil
	}
	session, ok := val.(*domain.Session)
	if !ok {
		return nil
	}
	return session
}
