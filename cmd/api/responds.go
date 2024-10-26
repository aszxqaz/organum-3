package main

import (
	"net/http"

	"github.com/go-chi/render"
)

func (app *application) respondInternalError(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusInternalServerError)
	render.PlainText(w, r, err.Error())
}

func (app *application) respondConflict(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusConflict)
	render.PlainText(w, r, err.Error())
}

func (app *application) respondNotFound(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusNotFound)
	render.PlainText(w, r, err.Error())
}

func (app *application) respondForbidden(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusForbidden)
	render.PlainText(w, r, err.Error())
}

func (app *application) respondBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusBadRequest)
	render.PlainText(w, r, err.Error())
}

func (app *application) respondUnauthorized(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusUnauthorized)
	render.PlainText(w, r, err.Error())
}
