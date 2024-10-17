package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MatthewAraujo/auth-limit-redis/internal/types"
	"github.com/MatthewAraujo/auth-limit-redis/internal/utils"
)

type Handler struct {
	store types.LoginStore
}

func NewHandler(store types.LoginStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var payload types.Login
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid input"))
		return
	}

	err := h.store.CreateUser(payload.Username, payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error creating user: %v", err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{"message": "User created successfully!"}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.Login
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid input"))
		return
	}

	token, err := h.store.LoginUser(&payload)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid credentials: %v", err))
		return
	}

	response := map[string]string{"token": token}
	json.NewEncoder(w).Encode(response)
}
