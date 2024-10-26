package store

import "organum/internal/domain"

func (s *Store) GetMSessionsForRoom(roomID string) []*domain.MSession {
	s.mu.RLock()
	defer s.mu.RUnlock()

	roomSessions := s.roomsSessions.ValuesWhere(
		func(k string, v *domain.RoomSession) bool {
			return v.RoomID == roomID
		},
	)

	msessions := s.wsSessions.ValuesWhere(
		func(k string, v *domain.MSession) bool {
			for _, rs := range roomSessions {
				if rs.SessionID == v.SessionID {
					return true
				}
			}
			return false
		},
	)

	return msessions
}
