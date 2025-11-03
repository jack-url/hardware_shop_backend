package main

import (
	"hardware_shop_backend/database"
	"hardware_shop_backend/routes"
	"log"
	"net/http"
	"os"
)

func main() {
	//  Step 1: Start server setup log

	log.Println(" Starting Hardware Shop Backend Server...")

	//  Step 2: Connect to MySQL database
	database.Connect()

	//  Step 3: Set up routes (Products, Customers, Orders)
	routes.SetupRoutes()

	//  Step 4: Define server port (default 8080)

	port := ":8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = ":" + fromEnv
	}

	//  Step 5: Start the server
	
	log.Printf(" Server running on http://localhost%s", port)
	log.Println(" Press CTRL+C to stop the server.")

	// Listen and serve
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
