package controller

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h Handler) getAllCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errorHandler(w, r, http.StatusMethodNotAllowed, "not allowed method")
		return
	}

	categories, code, err := h.service.Category.GetAllCategorys(r.Context())
	if err != nil {
		log.Println(err)
		h.errorHandler(w, r, code, "internal server error")
		return
	}

	if err := json.NewEncoder(w).Encode(categories); err != nil {
		h.errorHandler(w, r, http.StatusInternalServerError, "internal server error")
		return
	}
}
