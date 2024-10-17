package limit

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/MatthewAraujo/auth-limit-redis/internal/types"
	"github.com/MatthewAraujo/auth-limit-redis/internal/utils"
)

type Handler struct {
	store types.LimitStore
}

func NewHandler(store types.LimitStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) HandleLimit(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	if token == "" {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("authorization token is missing"))
		return
	}

	token = strings.Split(token, " ")[1]
	log.Printf("token: %s", token)

	err := h.store.IncrementTokenCount(token)
	if err != nil {

		if err == ErrLimitRequest {
			utils.WriteError(w, http.StatusTooManyRequests, ErrLimitRequest)
			return
		}

		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("authorization token is missing"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Token received successfully"))
}
