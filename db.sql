-- Create the users table
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(50) NOT NULL UNIQUE,
                       email VARCHAR(100) NOT NULL UNIQUE,
                       password VARCHAR(100) NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create the products table
CREATE TABLE products (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(100) NOT NULL,
                          description TEXT,
                          price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          user_id INTEGER REFERENCES users(id) ON DELETE SET NULL
);


-- Insert mock users
INSERT INTO users (username, email, password)
VALUES
    ('alice', 'alice@example.com', 'password123'),
    ('bob', 'bob@example.com', 'password123'),
    ('charlie', 'charlie@example.com', 'password123');

-- Insert mock products, associating them with users by user_id
INSERT INTO products (name, description, price, user_id)
VALUES
    ('Laptop', 'High-performance laptop for gaming and work', 1200.00, 1),
    ('Smartphone', 'Latest model smartphone with all features', 800.00, 2),
    ('Headphones', 'Noise-canceling headphones', 150.00, 1),
    ('Tablet', 'Lightweight tablet for on-the-go productivity', 300.00, 3),
    ('Monitor', '4K UHD monitor for immersive viewing', 400.00, 2),
    ('Keyboard', 'Mechanical keyboard with RGB lighting', 90.00, 1);
