package limit

import (
	"fmt"
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

	tokenParts := strings.Split(token, " ")
	if len(tokenParts) < 2 {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid token format"))
		return
	}
	token = tokenParts[1]

	expired, err := h.store.TokenIsExpired(token)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("internal server error"))
		return
	}

	if expired {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("token is expired"))
		return
	}

	err = h.store.IncrementTokenCount(token)
	if err != nil {
		if err == ErrLimitRequest {
			utils.WriteError(w, http.StatusTooManyRequests, ErrLimitRequest)
			return
		}

		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to increment token count"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Token received successfully"))
}
