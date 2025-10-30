package routes

import (
	"hardware_shop_backend/handlers"
	"log"
	"net/http"
)

// SetupRoutes registers all routes for the API calling 

func SetupRoutes() {
	log.Println("Setting up routes...")

	// PRODUCT ROUTES 

	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("📦 %s /products called", r.Method)
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
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	// CUSTOMER ROUTES 

	http.HandleFunc("/customers", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("🧍 %s /customers called", r.Method)
		switch r.Method {
		case http.MethodGet:
			id := r.URL.Query().Get("id")
			if id != "" {
				handlers.GetCustomerByID(w, r)
			} else {
				handlers.GetCustomers(w, r)
			}
		case http.MethodPost:
			handlers.CreateCustomer(w, r)
		case http.MethodPut:
			handlers.UpdateCustomer(w, r)
		case http.MethodDelete:
			handlers.DeleteCustomer(w, r)
		default:
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	// ORDER ROUTES
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		log.Printf(" %s /orders called", r.Method)
		switch r.Method {
		case http.MethodGet:
			id := r.URL.Query().Get("id")
			if id != "" {
				handlers.GetOrderDetails(w, r)
			} else {
				handlers.GetOrders(w, r)
			}
		case http.MethodPost:
			handlers.CreateOrder(w, r)
		default:
			handlers.ErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	log.Println(" Routes setup complete.")
}
