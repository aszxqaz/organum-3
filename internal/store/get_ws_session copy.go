package store

import (
	"errors"
	"organum/internal/domain"
)

var ErrWsSessionNotFound = errors.New("ws session not found")

func (s *Store) GetWsSession(session *domain.Session) (*domain.MSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ms, ok := s.wsSessions.Get(session.ID)
	if !ok {
		return nil, ErrWsSessionNotFound
	}

	return ms, nil
}
