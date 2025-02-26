package handlers

import (
	"errors"
	"net/http"

	"github.com/ctestabu/test_task/storage"
)

// DownloadAssetHandler handles file retrieval.
func DownloadAssetHandler(pgDB *storage.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Extract token
		token := getToken(r)
		if token == "" {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		// Validate session
		userID, err := pgDB.ValidateSession(ctx, token)
		if err != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		// Extract asset name
		assetName := getAssetName(r)
		if assetName == "" {
			http.Error(w, `{"error":"invalid asset name"}`, http.StatusBadRequest)
			return
		}

		// Retrieve asset
		data, err := pgDB.GetAsset(ctx, userID, assetName)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			} else {
				http.Error(w, `{"error":"server error"}`, http.StatusInternalServerError)
			}
			return
		}

		// Send file
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(data)
	}
}
