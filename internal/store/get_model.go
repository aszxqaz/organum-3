package store

import (
	"errors"
	"organum/internal/domain"
)

var ErrModelNotFound = errors.New("model not found")

func (s *Store) GetModel(checksum string) (*domain.ModelJSON, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	model, ok := s.models.Get(checksum)
	if !ok {
		return nil, ErrModelNotFound
	}

	return domain.NewModelJSON(model), nil
}
