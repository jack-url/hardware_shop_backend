# Hardware Shop Backend API

A simple and efficient **Go (Golang)** backend system for managing a hardware shop.  
It provides RESTful API endpoints to manage **Customers**, **Products**, and **Orders** with full CRUD operations.

##  Project Overview

This backend project enables a hardware store to:
- Register and manage customers  
- Add, update, and delete products  
- Record and view customer orders  
- Handle database persistence using **MySQL**  
- Communicate via structured **JSON APIs**

Built using:
- **Go (Golang)**
- **MySQL**
- **net/http**
- **encoding/json**

---

## Base URL
When running locally:
http://localhost:8080

##  API Endpoints

Below are all the available endpoints.  
You can test them using **Postman**
---

### Customer Endpoints

| Method | Endpoint | Description |
|--------|-----------|-------------|
| POST | `/customers` | Create a new customer |
| GET | `/customers` | Get all customers |
| GET | `/customers?id={id}` | Get a customer by ID |
| PUT | `/customers?id={id}` | Update a customer |
| DELETE | `/customers?id={id}` | Delete a customer |

**Examples:**
POST http://localhost:8080/customers
GET http://localhost:8080/customers
GET http://localhost:8080/customers?id=1
PUT http://localhost:8080/customers?id=1
DELETE http://localhost:8080/customers?id=1
---

### Product Endpoints

| Method | Endpoint | Description |
|--------|-----------|-------------|
| POST | `/products` | Create a new product |
| GET | `/products` | Get all products |
| GET | `/products?id={id}` | Get a product by ID |
| PUT | `/products?id={id}` | Update a product |
| DELETE | `/products?id={id}` | Delete a product |

**Examples:**
POST http://localhost:8080/products
GET http://localhost:8080/products
GET http://localhost:8080/products?id=1
PUT http://localhost:8080/products?id=1
DELETE http://localhost:8080/products?id=1
---

### Order Endpoints

| Method | Endpoint | Description |
|--------|-----------|-------------|
| POST | `/orders` | Create a new order |
| GET | `/orders` | Get all orders |
| GET | `/orders?id={id}` | Get an order by ID |
| PUT | `/orders?id={id}` | Update an order |
| DELETE | `/orders?id={id}` | Delete an order |

**Examples:**
POST http://localhost:8080/orders
GET http://localhost:8080/orders
GET http://localhost:8080/orders?id=1
PUT http://localhost:8080/orders?id=1
DELETE http://localhost:8080/orders?id=1

##  Example JSON Bodies

### ➕ Create Customer

{
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "0712345678"

}
➕ Create Product
json
{
  "name": "cement",
  "price": 500,
  "stock": 10
}
➕ Create Order
json
  {
  "customer_id": 1,
  "total": 1200.50,
  "created_at": "2025-11-04T10:15:00Z"
}
 How to Test
 Using Postman
Open Postman

Enter the base URL http://localhost:8080

Choose the correct HTTP method (GET, POST, PUT, DELETE)

Enter the endpoint (e.g. /customers)

If required, go to the Body tab → select raw → choose JSON and paste example data

Click Send
