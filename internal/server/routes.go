package server

import (
	"encoding/json"
	"net/http"

	"github.com/MatthewAraujo/auth-limit-redis/cmd/service/limit"
	"github.com/MatthewAraujo/auth-limit-redis/cmd/service/user"
	"github.com/MatthewAraujo/auth-limit-redis/internal/database"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	db := database.New()
	userStore := user.NewRedisUserStore(db.GetRedisClient())
	userHandler := user.NewHandler(userStore)

	limitStore := limit.NewTokenStore(db.GetRedisClient())
	limitHandler := limit.NewHandler(limitStore)

	mux.HandleFunc("/", helloWeb)
	mux.HandleFunc("/login", userHandler.HandleLogin)
	mux.HandleFunc("/create", userHandler.HandleCreateUser)
	mux.HandleFunc("/limit", limitHandler.HandleLimit)

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
