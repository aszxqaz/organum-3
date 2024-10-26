package store

import (
	"errors"
	"organum/internal/domain"
)

var ErrTokenNotFound = errors.New("token not found")
var ErrSessionNotFound = errors.New("session not found")

func (s *Store) GetSessionByToken(plain string) (*domain.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	token, ok := s.tokens.FirstWhere(func(k string, v *domain.Token) bool { return v.Token == plain })
	if !ok {
		return nil, ErrTokenNotFound
	}

	session, ok := s.sessions.Get(token.SessionID)
	if !ok {
		return nil, ErrSessionNotFound
	}

	return session, nil
}
