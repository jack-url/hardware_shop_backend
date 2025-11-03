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

// POST /products
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
	if p.Name == "" {
		errorResponse(w, http.StatusBadRequest, "name is required")
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

	// Log the created product
	log.Printf("Product created: ID=%d, Name=%s, Price=%.2f, Stock=%d", p.ID, p.Name, p.Price, p.Stock)

	jsonResponse(w, http.StatusCreated, p)
}

	

	



// GET /products
func GetProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, price, stock FROM products")
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "database query failed")
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

// GET /products?id=#
func GetProductByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		errorResponse(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	var p models.Product
	err = database.DB.QueryRow("SELECT id, name, price, stock FROM products WHERE id=?", id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if err != nil {
		if err == sql.ErrNoRows {
			errorResponse(w, http.StatusNotFound, "product not found")
		} else {
			errorResponse(w, http.StatusInternalServerError, "query error")
		}
		return
	}

	jsonResponse(w, http.StatusOK, p)
}

// PUT /products?id=#
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		errorResponse(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	_, err = database.DB.Exec(
		"UPDATE products SET name=?, price=?, stock=? WHERE id=?",
		p.Name, p.Price, p.Stock, id,
	)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "update failed")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "product updated"})
}

// DELETE /products?id=#
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		errorResponse(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	res, err := database.DB.Exec("DELETE FROM products WHERE id=?", id)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "delete failed")
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		errorResponse(w, http.StatusNotFound, "product not found")
		return
	}

	// Log deleted product
	log.Printf("Product deleted: ID=%d", id)

	jsonResponse(w, http.StatusOK, map[string]string{"message": "product deleted"})
}
