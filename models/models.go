package models

// Product in the shop.

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

// Customer who buys products.

type Customer struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

// Order  customer's purchase.
type Order struct {
	ID         int     `json:"id"`
	CustomerID int     `json:"customer_id"`
	Total      float64 `json:"total"`
	CreatedAt  string  `json:"created_at"`
}

//  and its products.
type OrderItem struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}
