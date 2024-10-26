package store

import (
	"errors"
	"io"
	"organum/internal/domain"
	"organum/internal/signer"
)

var ErrRoomSessionNotFound = errors.New("room session entry not found")
var ErrReadingFile = errors.New("error reading file")

func (s *Store) CreateModel(session *domain.Session, file io.Reader) (*domain.Model, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rs, ok := s.roomsSessions.FirstWhere(
		func(k string, v *domain.RoomSession) bool { return v.SessionID == session.ID },
	)
	if !ok {
		return nil, ErrRoomSessionNotFound
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, ErrReadingFile
	}

	model := &domain.Model{
		Checksum:   signer.Checksum(bytes),
		Bytes:      bytes,
		RoomID:     rs.RoomID,
		UploaderID: session.ID,
	}

	s.models.Set(model.Checksum, model)
	return model, nil
}
