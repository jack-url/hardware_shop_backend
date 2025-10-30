package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"hardware_shop_backend/database"
	"hardware_shop_backend/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)


// POST /products → create a new product
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	p.Name = strings.TrimSpace(p.Name)
	if p.Name == "" || p.Price <= 0 {
		errorResponse(w, http.StatusBadRequest, "name and positive price required")
		return
	}

	res, err := database.DB.Exec(
		"INSERT INTO products (name, price, stock) VALUES (?, ?, ?)",
		p.Name, p.Price, p.Stock,
	)
	if err != nil {
		log.Printf("CreateProduct: %v", err)
		errorResponse(w, http.StatusInternalServerError, "failed to create product")
		return
	}

	id, _ := res.LastInsertId()
	p.ID = int(id)
	jsonResponse(w, http.StatusCreated, p)
}

// GET /products → list all products

func GetProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, price, stock FROM products")
	if err != nil {
		log.Printf("GetProducts: %v", err)
		errorResponse(w, http.StatusInternalServerError, "database error")
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			errorResponse(w, http.StatusInternalServerError, "scan error")
			return
		}
		products = append(products, p)
	}

	jsonResponse(w, http.StatusOK, products)
}

// GET /products id

func GetProductByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		errorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	var p models.Product
	err = database.DB.QueryRow(
		"SELECT id, name, price, stock FROM products WHERE id = ?", id,
	).Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponse(w, http.StatusNotFound, "product not found")
			return
		}
		errorResponse(w, http.StatusInternalServerError, "query error")
		return
	}

	jsonResponse(w, http.StatusOK, p)
}

// PUT /products?id=#
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		errorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	p.Name = strings.TrimSpace(p.Name)
	if p.Name == "" || p.Price <= 0 {
		errorResponse(w, http.StatusBadRequest, "name and positive price required")
		return
	}

	res, err := database.DB.Exec(
		"UPDATE products SET name=?, price=?, stock=? WHERE id=?",
		p.Name, p.Price, p.Stock, id,
	)
	if err != nil {
		log.Printf("UpdateProduct: %v", err)
		errorResponse(w, http.StatusInternalServerError, "update failed")
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		errorResponse(w, http.StatusNotFound, "product not found")
		return
	}

	p.ID = id
	jsonResponse(w, http.StatusOK, p)
}

// DELETE /products?id=#
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		errorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	res, err := database.DB.Exec("DELETE FROM products WHERE id=?", id)
	if err != nil {
		log.Printf("DeleteProduct: %v", err)
		errorResponse(w, http.StatusInternalServerError, "delete failed")
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		errorResponse(w, http.StatusNotFound, "product not found")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"deleted": "ok"})
}
