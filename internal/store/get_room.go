package store

import "organum/internal/domain"

func (s *Store) GetRoom(roomID string) (*domain.Room, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	room, ok := s.rooms.Get(roomID)
	if !ok {
		return nil, ErrRoomNotFound
	}
	return room, nil
}
