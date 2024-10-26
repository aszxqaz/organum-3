package store

import (
	"errors"
	"organum/internal/domain"
)

var (
	ErrSessionNotInRoom = errors.New("session not in room")
)

func (s *Store) LeaveRoom(session *domain.Session, roomID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.roomsSessions.FirstWhere(
		func(k string, v *domain.RoomSession) bool {
			return v.SessionID == session.ID && v.RoomID == roomID
		},
	)
	if !ok {
		return ErrSessionNotInRoom
	}

	_, ok = s.rooms.Get(roomID)
	if !ok {
		return ErrRoomNotFound
	}

	s.leaveRoomUnsafe(session, roomID)
	return nil
}

func (s *Store) leaveRoomUnsafe(session *domain.Session, roomID string) {
	s.roomsSessions.DeleteWhere(
		func(k string, v *domain.RoomSession) bool {
			return v.SessionID == session.ID && v.RoomID == roomID
		},
	)

	if rs := s.roomsSessions.ValuesWhere(
		func(k string, v *domain.RoomSession) bool {
			return v.RoomID == roomID
		},
	); len(rs) == 0 {
		s.models.DeleteWhere(func(k string, v *domain.Model) bool { return v.RoomID == roomID })
		s.rooms.Delete(roomID)
	}
}
