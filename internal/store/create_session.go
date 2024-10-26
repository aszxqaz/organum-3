package store

import (
	"organum/internal/domain"
	"organum/internal/signer"
)

func (s *Store) CreateSession(secret string) (*domain.Session, *domain.Token) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session := domain.NewSession()
	token := domain.NewToken(session.ID, signer.Sign(session.ID, secret))

	s.sessions.Set(session.ID, session)
	s.tokens.Set(token.ID, token)

	return session, token
}
