package store

import "organum/internal/domain"

func (s *Store) GetRooms() []*domain.Room {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.rooms.Values()
}
