package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ctestabu/test_task/storage"
)

func ListAssetsHandler(pgDB *storage.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token := getToken(r)
		if token == "" {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		userID, err := pgDB.ValidateSession(ctx, token)
		if err != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		assets, err := pgDB.ListAssets(ctx, userID)
		if err != nil {
			http.Error(w, `{"error":"server error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(assets); err != nil {
			http.Error(w, `{"error":"server error"}`, http.StatusInternalServerError)
		}
	}
}
