package handlers

import (
	"database/sql"
	"encoding/json"
	"hardware_shop_backend/database"
	"hardware_shop_backend/models"
	"log"
	"net/http"
	"strconv"
)

// ----------------------- Utility Response Helpers -----------------------

func jsonResponse(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func errorResponse(w http.ResponseWriter, status int, msg string) {
	jsonResponse(w, status, map[string]string{"error": msg})
}

// ----------------------- ORDER HANDLERS -----------------------

// POST /orders
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid JSON data")
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "could not start transaction")
		return
	}

	// Insert the order
	res, err := tx.Exec("INSERT INTO orders (customer_id, total, created_at) VALUES (?, ?, NOW())",
		order.CustomerID, order.Total)
	if err != nil {
		tx.Rollback()
		log.Printf("CreateOrder insert order: %v", err)
		errorResponse(w, http.StatusInternalServerError, "failed to create order")
		return
	}

	orderID, _ := res.LastInsertId()

	// Insert each order item
	for _, item := range order.Items {
		_, err := tx.Exec("INSERT INTO order_items (order_id, product_id, quantity, unit_price) VALUES (?, ?, ?, ?)",
			orderID, item.ProductID, item.Quantity, item.UnitPrice)
		if err != nil {
			tx.Rollback()
			log.Printf("CreateOrder insert item: %v", err)
			errorResponse(w, http.StatusInternalServerError, "failed to insert order item")
			return
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("CreateOrder commit: %v", err)
		errorResponse(w, http.StatusInternalServerError, "failed to commit order")
		return
	}

	order.ID = int(orderID)
	log.Printf("Order created successfully: ID=%d", order.ID)
	jsonResponse(w, http.StatusCreated, order)
}

// GET /orders
func GetOrders(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
		SELECT o.id, o.customer_id, c.name, o.total, o.created_at
		FROM orders o
		JOIN customers c ON o.customer_id = c.id
	`)
	if err != nil {
		log.Printf("GetOrders: %v", err)
		errorResponse(w, http.StatusInternalServerError, "database query error")
		return
	}
	defer rows.Close()

	var orders []map[string]interface{}
	for rows.Next() {
		var (
			id, customerID int
			customerName   string
			total          float64
			createdAt      string
		)
		if err := rows.Scan(&id, &customerID, &customerName, &total, &createdAt); err != nil {
			errorResponse(w, http.StatusInternalServerError, "failed to scan row")
			return
		}

		orders = append(orders, map[string]interface{}{
			"id":            id,
			"customer_id":   customerID,
			"customer_name": customerName,
			"total":         total,
			"created_at":    createdAt,
		})
	}

	jsonResponse(w, http.StatusOK, orders)
}

// GET /orders?id=#
func GetOrderByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		errorResponse(w, http.StatusBadRequest, "invalid order ID")
		return
	}

	var order models.Order
	err = database.DB.QueryRow("SELECT id, customer_id, total, created_at FROM orders WHERE id = ?", id).
		Scan(&order.ID, &order.CustomerID, &order.Total, &order.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			errorResponse(w, http.StatusNotFound, "order not found")
		} else {
			log.Printf("GetOrderByID: %v", err)
			errorResponse(w, http.StatusInternalServerError, "database error")
		}
		return
	}

	// Get order items
	rows, err := database.DB.Query("SELECT id, order_id, product_id, quantity, unit_price FROM order_items WHERE order_id = ?", id)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "failed to fetch order items")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.UnitPrice); err == nil {
			order.Items = append(order.Items, item)
		}
	}

	jsonResponse(w, http.StatusOK, order)
}

// DELETE /orders?id=#
func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		errorResponse(w, http.StatusBadRequest, "invalid order ID")
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "could not start transaction")
		return
	}

	// Delete items first (foreign key constraint)
	_, err = tx.Exec("DELETE FROM order_items WHERE order_id = ?", id)
	if err != nil {
		tx.Rollback()
		errorResponse(w, http.StatusInternalServerError, "failed to delete order items")
		return
	}

	// Then delete the order
	res, err := tx.Exec("DELETE FROM orders WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		errorResponse(w, http.StatusInternalServerError, "failed to delete order")
		return
	}

	tx.Commit()
	rows, _ := res.RowsAffected()
	if rows == 0 {
		errorResponse(w, http.StatusNotFound, "order not found")
		return
	}

	log.Printf("Order deleted: ID=%d", id)
	jsonResponse(w, http.StatusOK, map[string]string{"message": "order deleted"})
}
