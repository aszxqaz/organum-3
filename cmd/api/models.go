package main

import (
	"errors"
	"net/http"
	"organum/internal/domain"
	"organum/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (app *application) postModelHandler(w http.ResponseWriter, r *http.Request) {
	session := app.getSessionOrRespondError(w, r)
	if session == nil {
		return
	}

	model, err := app.store.CreateModel(session, r.Body)
	if err != nil {
		if errors.Is(err, store.ErrRoomSessionNotFound) {
			app.respondBadRequest(w, r, err)
			return
		}
		if errors.Is(err, store.ErrReadingFile) {
			app.respondInternalError(w, r, err)
			return
		}
		app.respondInternalError(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, domain.NewModelJSON(model))
}

func (app *application) getModelsHandler(w http.ResponseWriter, r *http.Request) {
	models := app.store.GetModels()
	render.Status(r, http.StatusOK)
	render.JSON(w, r, models)
}

func (app *application) getModelBytesHandler(w http.ResponseWriter, r *http.Request) {
	session := app.getSessionOrRespondError(w, r)
	if session == nil {
		return
	}

	checksum := chi.URLParam(r, "checksum")
	bytes, err := app.store.GetModelBytes(session, checksum)
	if err != nil {
		if errors.Is(err, store.ErrModelNotFound) {
			app.respondNotFound(w, r, err)
			return
		}
		if errors.Is(err, store.ErrModelInOtherRoom) {
			app.respondForbidden(w, r, err)
			return
		}
		app.respondInternalError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func (app *application) getModelHandler(w http.ResponseWriter, r *http.Request) {
	checksum := chi.URLParam(r, "checksum")
	model, err := app.store.GetModel(checksum)
	if err != nil {
		if errors.Is(err, store.ErrModelNotFound) {
			app.respondNotFound(w, r, err)
			return
		}
		app.respondInternalError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, model)
}
