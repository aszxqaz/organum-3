package store

import "organum/internal/domain"

func (s *Store) JoinRoom(session *domain.Session, roomID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.roomsSessions.FirstWhere(
		func(k string, v *domain.RoomSession) bool {
			return v.SessionID == session.ID
		},
	)
	if ok {
		return ErrSessionAlreadyInRoom
	}

	_, ok = s.rooms.Get(roomID)
	if !ok {
		return ErrRoomNotFound
	}

	s.joinRoomUnsafe(session, roomID)
	return nil
}

func (s *Store) joinRoomUnsafe(session *domain.Session, roomID string) {
	roomSession := domain.NewRoomSession(roomID, session.ID)
	s.roomsSessions.Set(roomSession.ID, roomSession)
}
