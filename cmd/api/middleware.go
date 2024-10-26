package main

import (
	"errors"
	"net/http"
	"organum/internal/store"
	"strings"
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		token, err := extractToken(r)
		if err != nil {
			if errors.Is(err, ErrInvalidToken) {
				app.respondUnauthorized(w, r, err)
				return
			}
			if errors.Is(err, ErrMissingToken) {
				next.ServeHTTP(w, r)
				return
			}
		}

		session, err := app.store.GetSessionByToken(token)
		if err != nil {
			if errors.Is(err, store.ErrTokenNotFound) {
				app.respondUnauthorized(w, r, err)
				return
			}
			app.respondInternalError(w, r, err)
			return
		}

		r = app.contextSetSession(r, session)
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := app.contextGetSession(r)
		if !ok {
			app.respondUnauthorized(w, r, errors.New("must be authenticated to access"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

var ErrMissingToken = errors.New("missing token")
var ErrInvalidToken = errors.New("invalid token")

func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrMissingToken
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", ErrInvalidToken
	}
	return parts[1], nil
}
