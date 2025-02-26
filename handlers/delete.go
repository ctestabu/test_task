package handlers

import (
	"net/http"

	"github.com/ctestabu/test_task/storage"
)

func DeleteAssetHandler(pgDB *storage.Postgres) http.HandlerFunc {
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

		assetName := getAssetName(r)
		if assetName == "" {
			http.Error(w, `{"error":"invalid asset name"}`, http.StatusBadRequest)
			return
		}

		err = pgDB.DeleteAsset(r.Context(), userID, assetName)
		if err != nil {
			if err.Error() == "asset not found" {
				http.Error(w, `{"error": "asset not found"}`, http.StatusNotFound)
				return
			}
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"deleted"}`))
	}
}
