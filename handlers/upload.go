package handlers

import (
	"io"
	"net/http"

	"github.com/ctestabu/test_task/storage"
)

// UploadAssetHandler handles file uploads.
func UploadAssetHandler(pgDB *storage.Postgres) http.HandlerFunc {
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

		// Read file content
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

		// Store asset
		if err := pgDB.StoreAsset(ctx, userID, assetName, data); err != nil {
			http.Error(w, `{"error":"server error"}`, http.StatusInternalServerError)
			return
		}

		w.Write([]byte(`{"status":"ok"}`))
	}
}
