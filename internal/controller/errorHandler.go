package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// errorResponse is a tiny helper payload returned on API errors.
type errorResponse struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

// errorHandler logs the failure and serializes a consistent JSON payload.
func (h *Handler) errorHandler(w http.ResponseWriter, r *http.Request, status int, text string) {
	log.Printf("%s %s [%s]\t%s%s - %d - %s\n", time.Now().Format("2006/01/02 15:04:05"), r.Proto, r.Method, r.Host, r.RequestURI, status, http.StatusText(status))
	log.Println(text)
	e := errorResponse{
		Status: status,
		Msg:    text,
	}
	w.WriteHeader(e.Status)



	if err := json.NewEncoder(w).Encode(e); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
