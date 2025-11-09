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

// ORDER HANDLERS 


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

	if order.CustomerID <= 0 {
		errorResponse(w, http.StatusBadRequest, "invalid customer ID")
		return
	}

	//FIXED: Pass actual fields, not struct

	res, err := database.DB.Exec(
		"INSERT INTO orders (customer_id, total) VALUES (?, ?)",
		order.CustomerID, order.Total,
	)
	if err != nil {
		log.Printf("CreateOrder: %v", err)
		errorResponse(w, http.StatusInternalServerError, "failed to create order")
		return
	}

	id, _ := res.LastInsertId()
	order.ID = int(id)
	jsonResponse(w, http.StatusCreated, order)
}

// GET /orders
func GetOrders(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
SELECT o.id, o.customer_id, c.name, o.total
FROM orders o
JOIN customers c ON o.customer_id = c.id
`)
	if err != nil {
		log.Printf("GetOrders: %v", err)
		errorResponse(w, http.StatusInternalServerError, "database query failed")
		return
	}
	defer rows.Close()

	var orders []map[string]interface{}
	for rows.Next() {
		var id, customerID int
		var name string
		var total float64

		if err := rows.Scan(&id, &customerID, &name, &total); err != nil {
			errorResponse(w, http.StatusInternalServerError, "scan error")
			return
		}

		order := map[string]interface{}{
			"id":          id,
			"customer_id": customerID,
			"customer":    name,
			"total":       total,
		}
		orders = append(orders, order)
	}

	jsonResponse(w, http.StatusOK, orders)
}

// GET /orders
func GetOrderByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		errorResponse(w, http.StatusBadRequest, "invalid order ID")
		return
	}

	var order models.Order
	var customerName string

	err = database.DB.QueryRow(`
SELECT o.id, o.customer_id, c.name, o.total
FROM orders o
JOIN customers c ON o.customer_id = c.id
WHERE o.id = ?`, id).Scan(&order.ID, &order.CustomerID, &customerName, &order.Total)

	if err != nil {
		if err == sql.ErrNoRows {
			errorResponse(w, http.StatusNotFound, "order not found")
		} else {
			errorResponse(w, http.StatusInternalServerError, "database error")
		}
		return
	}

	resp := map[string]interface{}{
		"id":          order.ID,
		"customer_id": order.CustomerID,
		"customer":    customerName,
		"total":       order.Total,
	}

	jsonResponse(w, http.StatusOK, resp)
}

// PUT /orders?id=#
func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		errorResponse(w, http.StatusBadRequest, "invalid order ID")
		return
	}

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid JSON data")
		return
	}

	_, err = database.DB.Exec(
		"UPDATE orders SET customer_id=?, total=? WHERE id=?",
		order.CustomerID, order.Total, id,
	)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "failed to update order")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "order updated"})
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

	res, err := database.DB.Exec("DELETE FROM orders WHERE id=?", id)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "failed to delete order")
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		errorResponse(w, http.StatusNotFound, "order not found")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "order deleted"})
}

