package ws

import (
	"github.com/olahol/melody"
)

const methodAuth = "auth"

type AuthWsRequest struct {
	WsRequest
	Token string `json:"token"`
}

type AuthWsResponse struct {
	WsResponse
}

func NewAuthWsResponseSuccess() AuthWsResponse {
	return AuthWsResponse{
		WsResponse: WsResponse{
			Method:  methodAuth,
			Success: true,
		},
	}
}

func NewAuthWsResponseError(err error) AuthWsResponse {
	return AuthWsResponse{
		WsResponse: WsResponse{
			Method:  methodAuth,
			Success: false,
			Error:   err.Error(),
		},
	}
}

func (h *WsHandler) handleAuthWsRequest(s *melody.Session, r AuthWsRequest) {
	session, err := h.store.GetSessionByToken(r.Token)
	if err != nil {
		h.respondJSON(s, NewAuthWsResponseError(err))
		return
	}

	h.store.CreateWsSession(session, s)
	s.Set("session", session)

	h.respondJSON(s, NewAuthWsResponseSuccess())
}
