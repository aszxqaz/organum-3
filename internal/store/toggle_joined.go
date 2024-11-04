package store

import (
	"errors"
	"organum/internal/domain"
	"slices"
)

var ErrSceneAlreadyJoined = errors.New("scene already joined")
var ErrSceneAlreadyUnjoined = errors.New("scene already unjoined")
var ErrSceneNotMultiobject = errors.New("scene is not multiobject")

func (s *Store) ToggleJoined(session *domain.Session, roomID string, checksum string, isJoined bool) error {
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

	if room.LockOwnerID == "" {
		return ErrRoomNotLocked
	}

	if room.LockOwnerID != session.ID {
		return ErrRoomLockedByOther
	}

	index := slices.IndexFunc(room.Scenes, func(s *domain.Scene) bool { return s.Checksum == checksum })
	if index == -1 {
		return ErrSceneNotFound
	}
	scene := room.Scenes[index]

	if len(scene.Objects) < 2 {
		return ErrSceneNotMultiobject
	}

	if scene.IsJoined && isJoined {
		return ErrSceneAlreadyJoined
	}

	if !scene.IsJoined && !isJoined {
		return ErrSceneAlreadyUnjoined
	}

	scene.IsJoined = isJoined

	return nil
}
