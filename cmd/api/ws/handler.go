package ws

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"organum/internal/domain"
	"organum/internal/store"

	"github.com/olahol/melody"
)

type WsHandler struct {
	store *store.Store
}

func NewWsHandler(store *store.Store) *WsHandler {
	return &WsHandler{store: store}
}

func (h *WsHandler) HandleConnect(s *melody.Session) {
}

func (h *WsHandler) HandleDisconnect(s *melody.Session) {
	if session, ok := s.Get("session"); ok {
		if session, ok := session.(*domain.Session); ok && session != nil {
			slog.Info(fmt.Sprintf("Session id %s disconnected", session.ID))
			h.store.LeaveAllRooms(session)
			h.store.DeleteWsSession(session)
		}
	}
}

func (h *WsHandler) HandleMessage(s *melody.Session, b []byte) {
	slog.Info(fmt.Sprintf("Received ws message: %s", string(b)))

	var request WsRequest
	err := json.Unmarshal(b, &request)
	if err != nil {
		slog.Error("failed to unmarshal ws message")
		return
	}

	switch request.Method {
	case methodAuth:
		var r AuthWsRequest
		h.unmarshalOrRespondOnError(s, b, &r)
		h.handleAuthWsRequest(s, r)

	case methodTransform:
		var r TransformWsRequest
		h.unmarshalOrRespondOnError(s, b, &r)
		h.handleTransform(s, r)

	default:
		s.Close()
		return
	}
}

func (h *WsHandler) unmarshalOrRespondOnError(s *melody.Session, b []byte, v any) {
	err := json.Unmarshal(b, &v)
	if err != nil {
		s.Write([]byte(err.Error()))
		return
	}
}
