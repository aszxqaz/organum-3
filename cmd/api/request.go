package main

type SetNameRequest struct {
	Name string `json:"name"`
}

type MeRequest struct {
}

type JoinRoomRequest struct {
	RoomID string `json:"roomId"`
}

type CreateRoomRequest struct {
	Name string `json:"name"`
}
