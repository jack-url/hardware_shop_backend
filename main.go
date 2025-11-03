package main

import (
	"hardware_shop_backend/database"
	"hardware_shop_backend/routes"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Starting Hardware Shop Backend Server...")

	// Connect to MySQL database
	database.Connect()

	// Set up routes (Products, Customers, Orders)
	routes.SetupRoutes()

	// Define server port (default 8080)
	port := ":8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = ":" + fromEnv
	}

	log.Printf("Server running on http://localhost%s", port)
	log.Println("Press CTRL+C to stop the server.")

	// Start the server
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
