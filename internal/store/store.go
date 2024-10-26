package store

import (
	"organum/internal/domain"
	"sync"

	"github.com/aszxqaz/automap"
)

type Store struct {
	mu            sync.RWMutex
	sessions      automap.Map[string, *domain.Session]
	rooms         automap.Map[string, *domain.Room]
	models        automap.Map[string, *domain.Model]
	tokens        automap.Map[string, *domain.Token]
	roomsSessions automap.Map[string, *domain.RoomSession]
	wsSessions    automap.Map[string, *domain.MSession]
}

func NewStore() *Store {
	return &Store{}
}
