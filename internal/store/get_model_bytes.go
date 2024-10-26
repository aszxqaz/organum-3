package store

import (
	"errors"
	"organum/internal/domain"
)

var (
	ErrModelInOtherRoom = errors.New("model not in session room")
)

func (s *Store) GetModelBytes(session *domain.Session, checksum string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	model, ok := s.models.Get(checksum)
	if !ok {
		return nil, ErrModelNotFound
	}

	_, ok = s.roomsSessions.FirstWhere(func(k string, v *domain.RoomSession) bool {
		return v.SessionID == session.ID && v.RoomID == model.RoomID
	})
	if !ok {
		return nil, ErrModelInOtherRoom
	}

	return model.Bytes, nil
}
