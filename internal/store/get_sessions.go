package store

import "organum/internal/domain"

func (s *Store) GetSessions() []*domain.Session {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessions.Values()
}
