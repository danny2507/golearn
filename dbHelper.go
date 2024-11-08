// dbhelper.go
package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

// InitDB initializes the database connection and sets up tables if missing
func InitDB() *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open a DB connection: %v", err)
	}

	// Verify connection
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	fmt.Println("Database connection successful")

	// Check if the database has the necessary tables
	if isDatabaseEmpty(db) {
		fmt.Println("Database is empty, initializing tables...")
		err := executeSQLFile(db, "db.sql")
		if err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}
		fmt.Println("Database tables created successfully.")
	}

	return db
}

// isDatabaseEmpty checks if the required tables exist in the database
func isDatabaseEmpty(db *sql.DB) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users')").Scan(&exists)
	if err != nil {
		log.Fatalf("Error checking for table existence: %v", err)
	}
	return !exists
}

// executeSQLFile reads an SQL file and executes its commands
func executeSQLFile(db *sql.DB, filepath string) error {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", filepath, err)
	}

	_, err = db.Exec(string(file))
	if err != nil {
		return fmt.Errorf("failed to execute SQL commands: %v", err)
	}

	return nil
}

// CloseDB closes the database connection
func CloseDB() {
	if db != nil {
		db.Close()
	}
}

// Product represents a product in the database

// AddProduct inserts a new product into the database
func AddProduct(name, description string, price float64, userID int) (int, error) {
	var productID int
	err := db.QueryRow(`
		INSERT INTO products (name, description, price, user_id)
		VALUES ($1, $2, $3, $4) RETURNING id`, name, description, price, userID).Scan(&productID)
	if err != nil {
		return 0, fmt.Errorf("could not insert product: %v", err)
	}
	return productID, nil
}

// GetProduct retrieves a single product by ID
func GetProduct(id int) (*Product, error) {
	product := &Product{}
	err := db.QueryRow(`
		SELECT id, name, description, price, user_id FROM products WHERE id = $1`, id).
		Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Product not found
		}
		return nil, fmt.Errorf("could not get product: %v", err)
	}
	return product, nil
}

// ListProducts retrieves all products from the database
func ListProducts() ([]Product, error) {
	rows, err := db.Query(`SELECT id, name, description, price, user_id FROM products`)
	if err != nil {
		return nil, fmt.Errorf("could not get products: %v", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.UserID); err != nil {
			return nil, fmt.Errorf("could not scan product: %v", err)
		}
		products = append(products, product)
	}
	return products, nil
}

// UpdateProduct updates an existing product in the database
func UpdateProduct(id int, name, description string, price float64) error {
	_, err := db.Exec(`
		UPDATE products SET name = $1, description = $2, price = $3 WHERE id = $4`, name, description, price, id)
	if err != nil {
		return fmt.Errorf("could not update product: %v", err)
	}
	return nil
}

// DeleteProduct removes a product from the database
func DeleteProduct(id int) error {
	_, err := db.Exec(`DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("could not delete product: %v", err)
	}
	return nil
}
