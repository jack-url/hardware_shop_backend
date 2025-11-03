package routes

import (
	"hardware_shop_backend/handlers"
	"log"
	"net/http"
)

// SetupRoutes configures all API endpoints.
func SetupRoutes() {
	// --- Products ---
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			id := r.URL.Query().Get("id")
			if id != "" {
				handlers.GetProductByID(w, r)
			} else {
				handlers.GetProducts(w, r)
			}
		case http.MethodPost:
			handlers.CreateProduct(w, r)
		case http.MethodPut:
			handlers.UpdateProduct(w, r)
		case http.MethodDelete:
			handlers.DeleteProduct(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// --- Customers ---
	http.HandleFunc("/customers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			id := r.URL.Query().Get("id")
			if id != "" {
				handlers.GetCustomerByID(w, r)
			} else {
				handlers.GetCustomers(w, r)
			}
	
