package store

import (
	"errors"
	"organum/internal/domain"
)

var ErrRoomAlreadyLocked = errors.New("room already locked")
var ErrRoomLockedByOther = errors.New("room locked by other session")

func (s *Store) LockRoom(session *domain.Session, roomID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, ok := s.rooms.Get(roomID)
	if !ok {
		return ErrRoomNotFound
	}

	_, ok = s.roomsSessions.FirstWhere(
		func(k string, v *domain.RoomSession) bool {
			return v.SessionID == session.ID && v.RoomID == roomID
		},
	)
	if !ok {
		return ErrSessionNotInRoom
	}

	if room.LockOwnerID != "" {
		return ErrRoomAlreadyLocked
	}

	room.LockOwnerID = session.ID
	s.rooms.Set(room.ID, room)

	return nil
}

func (s *Store) UnlockRoom(session *domain.Session, roomID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, ok := s.rooms.Get(roomID)
	if !ok {
		return ErrRoomNotFound
	}

	_, ok = s.roomsSessions.FirstWhere(
		func(k string, v *domain.RoomSession) bool {
			return v.SessionID == session.ID && v.RoomID == roomID
		},
	)
	if !ok {
		return ErrSessionNotInRoom
	}

	if room.LockOwnerID != session.ID {
		return ErrRoomLockedByOther
	}

	room.LockOwnerID = ""
	s.rooms.Set(room.ID, room)

	return nil
}
