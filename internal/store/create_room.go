package store

import (
	"errors"
	"organum/internal/domain"
)

var (
	ErrSessionAlreadyInRoom = errors.New("session already in room")
)

func (s *Store) CreateRoom(session *domain.Session) (*domain.Room, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.roomsSessions.FirstWhere(func(k string, v *domain.RoomSession) bool { return v.SessionID == session.ID })
	if ok {
		return nil, ErrSessionAlreadyInRoom
	}

	return s.createRoomUnsafe(session), nil
}

func (s *Store) createRoomUnsafe(session *domain.Session) *domain.Room {
	room := domain.NewRoom(session.ID)
	roomSession := domain.NewRoomSession(room.ID, session.ID)
	s.rooms.Set(room.ID, room)
	s.roomsSessions.Set(roomSession.ID, roomSession)

	return room
}
