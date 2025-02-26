package handlers

import (
	"net/http"
	"strings"
)

func getToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}

// getAssetNam extracts the asset name from the request URL.
func getAssetName(r *http.Request) string {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) > 2 {
		return parts[len(parts)-1]
	}
	return ""
}
