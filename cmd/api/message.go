package main

import "organum/internal/domain"

type WsBroadcast struct {
	Method    string `json:"method"`
	SessionID string `json:"sessionId"`
	RoomID    string `json:"roomId"`
}

type WsBroadcastLock struct {
	WsBroadcast
	LockOwnerId string `json:"lockOwnerId"`
}

func NewWsBroadcastLock(session *domain.Session, roomID string, lockOwnerID string) *WsBroadcastLock {
	return &WsBroadcastLock{
		WsBroadcast: WsBroadcast{
			Method:    "lock",
			SessionID: session.ID,
			RoomID:    roomID,
		},
		LockOwnerId: lockOwnerID,
	}
}

type WsBroadcastScene struct {
	WsBroadcast
	Scene *domain.Scene `json:"scene"`
}

func NewWsBroadcastScene(session *domain.Session, roomID string, scene *domain.Scene) *WsBroadcastScene {
	return &WsBroadcastScene{
		WsBroadcast: WsBroadcast{
			Method:    "scene",
			SessionID: session.ID,
			RoomID:    roomID,
		},
		Scene: scene,
	}
}

type WsBroadcastToggleJoined struct {
	WsBroadcast
	Checksum string `json:"checksum"`
	IsJoined bool   `json:"isJoined"`
}

func NewWsBroadcastToggleJoinedScene(session *domain.Session, roomID string, checksum string, isJoined bool) *WsBroadcastToggleJoined {
	return &WsBroadcastToggleJoined{
		Checksum: checksum,
		IsJoined: isJoined,
		WsBroadcast: WsBroadcast{
			Method:    "toggleJoined",
			SessionID: session.ID,
			RoomID:    roomID,
		},
	}
}
