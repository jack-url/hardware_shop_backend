package handlers

import (
	"database/sql"
	"encoding/json"
	"hardware_shop_backend/database"
	"hardware_shop_backend/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Utility Response Helpers 

func jsonResponse(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func errorResponse(w http.ResponseWriter, status int, msg string) {
	jsonResponse(w, status, map[string]string{"error": msg})
}

// CUSTOMER HANDLERS 

// POST /customers
func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var c models.Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid JSON data")
		return
	}

	c.Name = strings.TrimSpace(c.Name)
	if c.Name == "" || c.Email == "" {
		errorResponse(w, http.StatusBadRequest, "name and email required")
		return
	}

	res, err := database.DB.Exec(
		"INSERT INTO customers (name, email, phone, address) VALUES (?, ?, ?, ?)",
		c.Name, c.Email, c.Phone, c.Address,
	)
	if err != nil {
		log.Printf("CreateCustomer: %v", err)
		errorResponse(w, http.StatusInternalServerError, "failed to create customer")
		return
	}

	id, _ := res.LastInsertId()
	c.ID = int(id)
	log.Printf("Customer created: ID=%d, Name=%s", c.ID, c.Name)
	jsonResponse(w, http.StatusCreated, c)
}

// GET /customers

func GetCustomers(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, email, phone, address FROM customers")
	if err != nil {
		log.Printf("GetCustomers: %v", err)
		errorResponse(w, http.StatusInternalServerError, "database query error")
		return
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var c models.Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Address); err != nil {
			log.Printf("GetCustomers scan: %v", err)
			errorResponse(w, http.StatusInternalServerError, "scan error")
			return
		}
		customers = append(customers, c)
	}

	jsonResponse(w, http.StatusOK, customers)
}

// GET /customers?id=#

func GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		errorResponse(w, http.StatusBadRequest, "invalid customer ID")
		return
	}

	var c models.Customer
	err = database.DB.QueryRow("SELECT id, name, email, phone, address FROM customers WHERE id = ?", id).
		Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			errorResponse(w, http.StatusNotFound, "customer not found")
		} else {
			log.Printf("GetCustomerByID: %v", err)
			errorResponse(w, http.StatusInternalServerError, "database error")
		}
		return
	}

	jsonResponse(w, http.StatusOK, c)
}

// PUT /customers?id=#

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		errorResponse(w, http.StatusBadRequest, "invalid customer ID")
		return
	}

	var c models.Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid JSON data")
		return
	}

	res, err := database.DB.Exec(
		"UPDATE customers SET name=?, email=?, phone=?, address=? WHERE id=?",
		c.Name, c.Email, c.Phone, c.Address, id,
	)
	if err != nil {
		log.Printf("UpdateCustomer: %v", err)
		errorResponse(w, http.StatusInternalServerError, "failed to update customer")
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		errorResponse(w, http.StatusNotFound, "customer not found")
		return
	}

	log.Printf("Customer updated: ID=%d", id)
	jsonResponse(w, http.StatusOK, map[string]string{"message": "customer updated"})
}

// DELETE /customers?id=#

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		errorResponse(w, http.StatusBadRequest, "invalid customer ID")
		return
	}

	res, err := database.DB.Exec("DELETE FROM customers WHERE id = ?", id)
	if err != nil {
		log.Printf("DeleteCustomer: %v", err)
		errorResponse(w, http.StatusInternalServerError, "failed to delete customer")
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		errorResponse(w, http.StatusNotFound, "customer not found")
		return
	}

	log.Printf("Customer deleted: ID=%d", id)
	jsonResponse(w, http.StatusOK, map[string]string{"message": "customer deleted"})
}
