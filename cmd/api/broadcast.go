package main

import (
	"encoding/json"
	"log/slog"
	"organum/internal/domain"
)

func (app *application) broadcast(session *domain.Session, roomID string, v any) {
	msessions := app.store.GetMSessionsForRoom(roomID)
	if len(msessions) > 0 {
		msg, _ := json.Marshal(v)
		for _, msession := range msessions {
			if msession.SessionID != session.ID {
				slog.Info("Broadcasting...", "roomID", roomID, "sessionID", msession.SessionID, "message", string(msg))
				msession.Melody.Write(msg)
			}
		}
	}
}
