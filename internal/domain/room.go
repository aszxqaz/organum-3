package domain

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Users    []string `json:"users"`
	Scenes      []*Scene  `json:"scenes"`
	OwnerID     string    `json:"ownerId"`
	LockOwnerID string    `json:"lockOwnerId"`
	CreatedAt   time.Time `json:"createdAt"`
}

// func (r *Room) Join(userId string) *errors.Error {
// 	if slices.Contains(r.Users, userId) {
// 		return errors.NewBadRequestError("user already joined room")
// 	}
// 	r.Users = append(r.Users, userId)
// 	return nil
// }

// func (r *Room) UnJoin(userId string) *errors.Error {
// 	index := slices.Index(r.Users, userId)
// 	if index == -1 {
// 		return errors.NewBadRequestError("user not found in room")
// 	}

// 	r.Users = append(r.Users[:index], r.Users[index+1:]...)
// 	return nil
// }

func (r *Room) AddScene(s *Scene) {
	r.Scenes = append(r.Scenes, s)
}

func NewRoom(sessionID string) *Room {
	return &Room{
		ID:          uuid.NewString(),
		Name:        "",
		Scenes:      []*Scene{},
		OwnerID:     sessionID,
		LockOwnerID: sessionID,
		CreatedAt:   time.Now(),
	}
}
