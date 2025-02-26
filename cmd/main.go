package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ctestabu/test_task/handlers"
	"github.com/ctestabu/test_task/storage"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Read database URL
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	// Initialize DB connection
	db, err := storage.NewPG(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Test DB connection
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Middleware for authentication
	// authMiddleware := storage.AuthMiddleware(db)

	// // Register routes
	http.Handle("/api/auth", handlers.AuthHandler(db))
	// http.Handle("/api/protected-endpoint", authMiddleware(http.HandlerFunc(protectedHandler)))

	http.Handle("/api/upload-asset/", handlers.UploadAssetHandler(db))
	http.Handle("/api/asset/", handlers.DownloadAssetHandler(db))
	http.Handle("/api/delete-asset/", handlers.DeleteAssetHandler(db))
	http.Handle("/api/list-assets", handlers.ListAssetsHandler(db))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
