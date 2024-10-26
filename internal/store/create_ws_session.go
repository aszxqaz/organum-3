package store

import (
	"organum/internal/domain"

	"github.com/olahol/melody"
)

func (s *Store) CreateWsSession(session *domain.Session, ms *melody.Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.wsSessions.Set(session.ID, domain.NewMSession(session.ID, ms))
}

func (s *Store) DeleteWsSession(session *domain.Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.wsSessions.Delete(session.ID)
}
