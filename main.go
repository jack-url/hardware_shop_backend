package main

import (
	"hardware_shop_backend/database"
	"hardware_shop_backend/routes"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting Hardware Shop Backend Server...")

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Connect to MySQL database (reads from env variables)
	database.Connect()

	// Set up routes
	routes.SetupRoutes()

	// Define server port (default 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on http://localhost:%s", port)
	log.Println("Press CTRL+C to stop the server.")

	// Start the server
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
