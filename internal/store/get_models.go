package store

import "organum/internal/domain"

func (s *Store) GetModels() []*domain.ModelJSON {
	s.mu.RLock()
	defer s.mu.RUnlock()

	models := s.models.Values()

	modelsJson := make([]*domain.ModelJSON, len(models))
	for i, model := range models {
		modelsJson[i] = domain.NewModelJSON(model)
	}

	return modelsJson
}
