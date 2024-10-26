package main

import (
	"context"
	"net/http"
	"organum/internal/domain"

	"github.com/go-chi/render"
)

type ctxKey string

const sessionCtxKey ctxKey = "session"

func (app *application) contextSetSession(r *http.Request, session *domain.Session) *http.Request {
	ctx := context.WithValue(r.Context(), sessionCtxKey, session)
	return r.WithContext(ctx)
}

func (app *application) contextGetSession(r *http.Request) (*domain.Session, bool) {
	session, ok := r.Context().Value(sessionCtxKey).(*domain.Session)
	if !ok {
		return nil, false
	}
	return session, true
}

func (app *application) getSessionOrRespondError(w http.ResponseWriter, r *http.Request) *domain.Session {
	session, ok := app.contextGetSession(r)
	if !ok || session == nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, "request session not found")
		return nil
	}
	return session
}
