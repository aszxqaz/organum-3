package store

import (
	"errors"
	"organum/internal/domain"
	"slices"
)

var (
	ErrSceneNotFound  = errors.New("scene not found")
	ErrObjectNotFound = errors.New("object not found")
)

func (s *Store) UpdateTransform(session *domain.Session, transform *domain.Transform, roomID string, checksum string, name string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

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

	index := slices.IndexFunc(room.Scenes, func(s *domain.Scene) bool { return s.Checksum == checksum })
	if index == -1 {
		return ErrSceneNotFound
	}
	scene := room.Scenes[index]

	if name != "" {
		index := slices.IndexFunc(scene.Objects, func(s *domain.Object) bool { return s.Name == name })
		if index == -1 {
			return ErrObjectNotFound
		}
		scene.Objects[index].Transform = *transform
	} else {
		scene.Transform = *transform
	}

	return nil
}
