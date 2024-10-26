package main

import (
	"net/http"
	"organum/cmd/api/ws"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/models", app.getModelsHandler)
	r.Post("/models", app.requireAuth(app.postModelHandler))
	r.Get("/models/{checksum}", app.getModelHandler)
	r.Get("/models/{checksum}/bytes", app.getModelBytesHandler)

	r.Get("/sessions", app.getSessionsHandler)
	r.Post("/sessions", app.postSessionHandler)

	r.Get("/rooms", app.getRoomsHandler)
	r.Post("/rooms", app.requireAuth(app.postRoomHandler))
	r.Get("/rooms/{roomID}", app.getRoomHandler)
	r.Delete("/rooms/{roomID}", app.requireAuth(app.deleteRoomHandler))

	r.Post("/rooms/{roomID}/lock", app.requireAuth(app.postLockHandler))
	r.Delete("/rooms/{roomID}/lock", app.requireAuth(app.deleteLockHandler))

	r.Get("/rooms/{roomID}/sessions", app.getRoomSessionsHandler)
	r.Post("/rooms/{roomID}/sessions", app.requireAuth(app.joinRoomHandler))
	r.Delete("/rooms/{roomID}/sessions", app.requireAuth(app.leaveRoomHandler))

	r.Post("/rooms/{roomID}/scenes", app.requireAuth(app.postSceneHandler))

	wsHandler := ws.NewWsHandler(app.store)

	app.melody.HandleConnect(wsHandler.HandleConnect)
	app.melody.HandleDisconnect(wsHandler.HandleDisconnect)
	app.melody.HandleMessage(wsHandler.HandleMessage)

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		app.melody.HandleRequest(w, r)
	})

	return app.authenticate(r)
}
