package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ctestabu/test_task/storage"
)

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func AuthHandler(pgDB *storage.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "invalid request"}`, http.StatusBadRequest)
			return
		}

		userID, err := pgDB.ValidateUser(r.Context(), req.Login, req.Password)
		if err != nil {
			http.Error(w, `{"error": "invalid login/password"}`, http.StatusUnauthorized)
			return
		}

		if err := pgDB.DeleteUserSession(r.Context(), userID); err != nil {
			http.Error(w, `{"error": "DeleteUserSession server error"}`, http.StatusInternalServerError)
			return
		}

		sessionID, err := pgDB.CreateSession(r.Context(), userID, r.RemoteAddr)
		if err != nil {
			http.Error(w, `{"error": "CreateSession server error"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"token": sessionID})
	}
}
