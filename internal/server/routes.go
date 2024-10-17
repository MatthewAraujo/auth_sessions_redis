package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", helloWeb)

	return mux
}

func helloWeb(w http.ResponseWriter, r *http.Request) {
	responseBody := map[string]string{
		"message": "healty route",
	}

	responseJSON, err := json.Marshal(responseBody)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
