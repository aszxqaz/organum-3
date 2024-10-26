package store

import (
	"errors"
	"organum/internal/domain"
)

var (
	ErrRoomNotFound = errors.New("room not found")
	ErrRoomNotOwner = errors.New("session is not room owner")
)

func (s *Store) DeleteRoom(session *domain.Session, roomID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, ok := s.rooms.Get(roomID)
	if !ok {
		return ErrRoomNotFound
	}

	if room.OwnerID != session.ID {
		return ErrRoomNotOwner
	}

	s.deleteRoomUnsafe(roomID)
	return nil
}

func (s *Store) deleteRoomUnsafe(roomID string) {
	s.roomsSessions.DeleteWhere(func(k string, v *domain.RoomSession) bool { return v.RoomID == roomID })
	s.models.DeleteWhere(func(k string, v *domain.Model) bool { return v.RoomID == roomID })
	s.rooms.Delete(roomID)
}
