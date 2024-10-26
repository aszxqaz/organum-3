package main

import (
	"net/http"

	"github.com/go-chi/render"
)

func (app *application) postSessionHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := app.contextGetSession(r)
	if ok {
		render.Status(r, http.StatusConflict)
		render.PlainText(w, r, "Request session already exists")
		return
	}

	session, token := app.store.CreateSession(app.config.secret)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, CreateSessionResponse{
		SessionID: session.ID,
		Token:     token.Token,
	})
}

func (app *application) getSessionsHandler(w http.ResponseWriter, r *http.Request) {
	sessions := app.store.GetSessions()
	render.Status(r, http.StatusOK)
	render.JSON(w, r, sessions)
}
