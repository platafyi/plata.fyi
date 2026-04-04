package handlers

import (
	"net/http"

	"github.com/platafyi/plata.fyi/internal/database"
)

type AuthHandler struct {
	store database.Store
}

func NewAuthHandler(store database.Store) *AuthHandler {
	return &AuthHandler{store: store}
}

func (h *AuthHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		jsonError(w, "Метод не е дозволен", http.StatusMethodNotAllowed)
		return
	}

	token := bearerToken(r)
	if token == "" {
		jsonError(w, "Потребна е автентикација", http.StatusUnauthorized)
		return
	}

	if err := h.store.DeleteToken(r.Context(), token); err != nil {
		jsonError(w, "Грешка при одјавување", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
