package controller

import (
	"encoding/json"
	"net/http"

	"rlf/internal/entity"
)

func (h *Handler) createComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errorHandler(w, r, http.StatusMethodNotAllowed, "not allowed method")
		return
	}
	userID := r.Context().Value(h.service.IDKey).(int)
	if userID < 0 {
		h.errorHandler(w, r, http.StatusUnauthorized, "invalid id")
		return
	}
	input := entity.Comment{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.errorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}
	input.UserID = uint(userID)
	if status, err := h.service.Comment.CreateComment(r.Context(), input); err != nil {
		h.errorHandler(w, r, status, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) voteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		h.errorHandler(w, r, http.StatusMethodNotAllowed, "not allowed method")
		return
	}
	userID := r.Context().Value(h.service.IDKey).(int)
	if userID < 0 {
		h.errorHandler(w, r, http.StatusUnauthorized, "invalid id")
		return
	}
	var input entity.CommentVote
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.errorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}
	input.UserID = uint(userID)
	if status, err := h.service.Comment.UpsertCommentVote(r.Context(), input); err != nil {
		h.errorHandler(w, r, status, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
