package store

import "organum/internal/domain"

func (s *Store) GetRoomSessions(roomID string) ([]*domain.RoomSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.rooms.Get(roomID)
	if !ok {
		return nil, ErrRoomNotFound
	}

	return s.roomsSessions.ValuesWhere(
		func(k string, v *domain.RoomSession) bool {
			return v.RoomID == roomID
		},
	), nil
}
