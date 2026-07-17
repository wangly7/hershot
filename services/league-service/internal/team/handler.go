package team

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) ListTeams(w http.ResponseWriter, r *http.Request) {
	teams, err := h.service.ListTeams(r.Context())
	if err != nil {
		log.Printf("list teams failed: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(teams); err != nil {
		log.Printf("encode teams response failed: %v", err)
	}
}

func (h *Handler) GetTeam(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	foundTeam, err := h.service.GetTeam(r.Context(), id)

	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, foundTeam)
}

func (h *Handler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var request CreateTeamRequest

	if err := decodeJSON(r, &request); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	createdTeam, err := h.service.CreateTeam(r.Context(), request)
	if err != nil {
		h.handleServiceError(w, err)
	}

	writeJSON(
		w,
		http.StatusOK,
		createdTeam,
	)
}

func (h *Handler) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var request UpdateTeamRequest

	if err := decodeJSON(r, &request); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	updatedTeam, err := h.service.UpdateTeam(r.Context(), id, request)
	if err != nil {
		h.handleServiceError(w, err)
	}

	writeJSON(
		w,
		http.StatusOK,
		updatedTeam,
	)
}

func (h *Handler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.DeleteTeam(r.Context(), id); err != nil {
		h.handleServiceError(w, err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTeamNotFound):
		writeError(
			w,
			http.StatusNotFound,
			ErrTeamNotFound.Error(),
		)

	case errors.Is(err, ErrInvalidID),
		errors.Is(err, ErrInvalidTeamName),
		errors.Is(err, ErrInvalidTeamCity),
		errors.Is(err, ErrInvalidTeamAbbreviation):
		writeError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
	default:
		log.Printf("team request failed: %v", err)
		writeError(
			w,
			http.StatusInternalServerError,
			"internal service error",
		)
	}
}

func decodeJSON(r *http.Request, destination any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		return err
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}

	return nil
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{
		Error: message,
	})
}
