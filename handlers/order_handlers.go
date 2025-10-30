package handlers

import (
	"encoding/json"
	"hardware_shop_backend/database"
	"hardware_shop_backend/models"
	"log"
	"net/http"
	"strconv"
)

// CreateOrder handles POST /orders
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	// Insert into orders table
	res, err := database.DB.Exec(
		"INSERT INTO orders (customer_id, total, created_at) VALUES (?, ?, NOW())",
		order.CustomerID, order.Total,
	)
	if err != nil {
		log.Printf("CreateOrder error: %v", err)
		errorResponse(w, http.StatusInternalServerError, "failed to create order")
		return
	}

	id, _ := res.LastInsertId()
	order.ID = int(id)
	jsonResponse(w, http.StatusCreated, order)
}

// GetOrders handles GET /orders
func GetOrders(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, customer_id, total, created_at FROM orders")
	if err != nil {
		log.Printf("GetOrders error: %v", err)
		errorResponse(w, http.StatusInternalServerError, "database error")
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.CustomerID, &o.Total, &o.CreatedAt); err != nil {
			log.Printf("GetOrders scan error: %v", err)
			errorResponse(w, http.StatusInternalServerError, "scan error")
			return
		}
		orders = append(orders, o)
	}

	jsonResponse(w, http.StatusOK, orders)
}

// GetOrderByID handles GET /orders?id=#
func GetOrderByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		errorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	var o models.Order
	err = database.DB.QueryRow(
		"SELECT id, customer_id, total, created_at FROM orders WHERE id = ?", id,
	).Scan(&o.ID, &o.CustomerID, &o.Total, &o.CreatedAt)
	if err != nil {
		log.Printf("GetOrderByID error: %v", err)
		errorResponse(w, http.StatusNotFound, "order not found")
		return
	}

	jsonResponse(w, http.StatusOK, o)
}
