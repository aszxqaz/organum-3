package main

import (
	"errors"
	"net/http"
	"organum/internal/domain"
	"organum/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (app *application) getRoomsHandler(w http.ResponseWriter, r *http.Request) {
	rooms := app.store.GetRooms()
	render.Status(r, http.StatusOK)
	render.JSON(w, r, rooms)
}

func (app *application) getRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")

	room, err := app.store.GetRoom(roomID)
	if err != nil {
		if errors.Is(err, store.ErrRoomNotFound) {
			app.respondNotFound(w, r, err)
			return
		}
		app.respondInternalError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, room)
}

func (app *application) postRoomHandler(w http.ResponseWriter, r *http.Request) {
	session := app.getSessionOrRespondError(w, r)
	if session == nil {
		return
	}

	room, err := app.store.CreateRoom(session)
	if err != nil {
		if errors.Is(err, store.ErrSessionAlreadyInRoom) {
			app.respondConflict(w, r, err)
			return
		}
		app.respondInternalError(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, room)
}

func (app *application) deleteRoomHandler(w http.ResponseWriter, r *http.Request) {
	session := app.getSessionOrRespondError(w, r)
	if session == nil {
		return
	}

	roomID := chi.URLParam(r, "roomID")

	if err := app.store.DeleteRoom(session, roomID); err != nil {
		if errors.Is(err, store.ErrRoomNotFound) {
			app.respondNotFound(w, r, err)
			return
		}
		if errors.Is(err, store.ErrRoomNotOwner) {
			app.respondForbidden(w, r, err)
			return
		}
		app.respondInternalError(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (app *application) joinRoomHandler(w http.ResponseWriter, r *http.Request) {
	session := app.getSessionOrRespondError(w, r)
	if session == nil {
		return
	}

	roomID := chi.URLParam(r, "roomID")
	if err := app.store.JoinRoom(session, roomID); err != nil {
		if errors.Is(err, store.ErrRoomNotFound) {
			app.respondNotFound(w, r, err)
			return
		}
		if errors.Is(err, store.ErrSessionAlreadyInRoom) {
			app.respondConflict(w, r, err)
			return
		}
		app.respondInternalError(w, r, err)
		return
	}

	room, err := app.store.GetRoom(roomID)
	if err != nil {
		app.respondInternalError(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, room)
}

func (app *application) leaveRoomHandler(w http.ResponseWriter, r *http.Request) {
	session := app.getSessionOrRespondError(w, r)
	if session == nil {
		return
	}

	roomID := chi.URLParam(r, "roomID")

	if err := app.store.LeaveRoom(session, roomID); err != nil {
		if errors.Is(err, store.ErrRoomNotFound) {
			app.respondNotFound(w, r, err)
			return
		}
		if errors.Is(err, store.ErrSessionNotInRoom) {
			app.respondNotFound(w, r, err)
			return
		}
		app.respondInternalError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
}

func (app *application) getRoomSessionsHandler(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")

	roomSessions, err := app.store.GetRoomSessions(roomID)
	if err != nil {
		if errors.Is(err, store.ErrRoomNotFound) {
			app.respondNotFound(w, r, err)
			return
		}
		app.respondInternalError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, roomSessions)
}

func (app *application) postLockHandler(w http.ResponseWriter, r *http.Request) {
	session := app.getSessionOrRespondError(w, r)
	if session == nil {
		return
	}

	roomID := chi.URLParam(r, "roomID")

	err := app.store.LockRoom(session, roomID)
	if err != nil {
		if errors.Is(err, store.ErrRoomNotFound) {
			app.respondNotFound(w, r, err)
			return
		}
		if errors.Is(err, store.ErrSessionNotInRoom) {
			app.respondForbidden(w, r, err)
			return
		}
		if errors.Is(err, store.ErrRoomAlreadyLocked) {
			app.respondConflict(w, r, err)
			return
		}
		app.respondInternalError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
}

func (app *application) deleteLockHandler(w http.ResponseWriter, r *http.Request) {
	session := app.getSessionOrRespondError(w, r)
	if session == nil {
		return
	}

	roomID := chi.URLParam(r, "roomID")

	err := app.store.UnlockRoom(session, roomID)
	if err != nil {
		if errors.Is(err, store.ErrRoomNotFound) {
			app.respondNotFound(w, r, err)
			return
		}
		if errors.Is(err, store.ErrSessionNotInRoom) {
			app.respondForbidden(w, r, err)
			return
		}
		if errors.Is(err, store.ErrRoomLockedByOther) {
			app.respondConflict(w, r, err)
			return
		}
		app.respondInternalError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
}

type PostSceneRequest struct {
	domain.Scene
}

// Bind implements render.Binder.
func (p *PostSceneRequest) Bind(r *http.Request) error {
	return nil
}

func (app *application) postSceneHandler(w http.ResponseWriter, r *http.Request) {
	session := app.getSessionOrRespondError(w, r)
	if session == nil {
		return
	}

	roomID := chi.URLParam(r, "roomID")
	var req PostSceneRequest
	render.Bind(r, &req)

	err := app.store.AddScene(session, roomID, &req.Scene)
	if err != nil {
		if errors.Is(err, store.ErrRoomNotFound) {
			app.respondNotFound(w, r, err)
			return
		}
		if errors.Is(err, store.ErrSessionNotInRoom) {
			app.respondForbidden(w, r, err)
			return
		}
		if errors.Is(err, store.ErrRoomNotLocked) {
			app.respondConflict(w, r, err)
			return
		}
		app.respondInternalError(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
}
