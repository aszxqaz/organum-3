package ws

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/olahol/melody"
)

func (h *WsHandler) respondJSON(s *melody.Session, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("Responding: %s", string(b)))
	return s.Write(b)
}

func (h *WsHandler) broadcastRoomOthers() {

}
