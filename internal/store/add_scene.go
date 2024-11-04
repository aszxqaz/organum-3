package store

import (
	"errors"
	"organum/internal/domain"
)

var ErrRoomNotLocked = errors.New("room not locked")
var ErrModelNotUploaded = errors.New("model not uploaded")

func (s *Store) AddScene(session *domain.Session, roomID string, scene *domain.Scene) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, ok := s.rooms.Get(roomID)
	if !ok {
		return ErrRoomNotFound
	}

	_, ok = s.models.Get(scene.Checksum)
	if !ok {
		return ErrModelNotUploaded
	}

	_, ok = s.roomsSessions.FirstWhere(
		func(k string, v *domain.RoomSession) bool { return v.SessionID == session.ID },
	)
	if !ok {
		return ErrRoomSessionNotFound
	}

	if room.LockOwnerID != session.ID {
		return ErrRoomNotLocked
	}

	room.Scenes = append(room.Scenes, scene)
	return nil
}
