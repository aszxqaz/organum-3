package ws

import (
	"encoding/json"
	"organum/internal/domain"

	"github.com/olahol/melody"
)

const methodTransform = "transform"

type TransformWsRequest struct {
	WsRequest
	Transform domain.Transform `json:"transform"`
	RoomID    string           `json:"roomId"`
	Checksum  string           `json:"checksum"`
	Name      string           `json:"name"`
}

type TransformWsBroadcast struct {
	TransformWsRequest
	SessionId string `json:"sessionId"`
}

type TransformWsResponseError struct {
	WsResponse
}

func NewTransformWsResponseError(err error) AuthWsResponse {
	return AuthWsResponse{
		WsResponse: WsResponse{
			Method:  methodTransform,
			Success: false,
			Error:   err.Error(),
		},
	}
}

func (h *WsHandler) handleTransform(s *melody.Session, r TransformWsRequest) {
	session := h.getSessionFromMelodySession(s)
	if session == nil {
		s.Write([]byte("Not authenticated"))
		return
	}

	err := h.store.UpdateTransform(session, &r.Transform, r.RoomID, r.Checksum, r.Name)
	if err != nil {
		h.respondJSON(s, NewTransformWsResponseError(err))
		return
	}

	brdMsg, _ := json.Marshal(TransformWsBroadcast{
		TransformWsRequest: r,
		SessionId:          session.ID,
	})

	msessions := h.store.GetMSessionsForRoom(r.RoomID)
	for _, msession := range msessions {
		if msession.SessionID != session.ID {
			msession.Melody.Write(brdMsg)
		}
	}
}
