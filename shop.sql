CREATE DATABASE shop;
DROP DATABASE shop;

CREATE TABLE customer (
    id serial PRIMARY KEY,
    firstname varchar(50) NOT NULL,
    lastname varchar(50) NOT NULL,
    phone_number varchar(30) NOT NULL,
    gender BOOLEAN NOT NULL,
    birth_date DATE NOT NULL,
    balance numeric(18, 2) DEFAULT 0,
    created_at TIMESTAMP  DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP 
);

CREATE TABLE orders (
    id serial PRIMARY KEY,
    customer_id int NOT NULL REFERENCES customer(id),
    total_amount DECIMAL(18, 2)
);

CREATE table order_items (
    id serial NOT NULL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(id),
    product_name VARCHAR,
    product_id INT NOT NULL REFERENCES products(id),
    count INT NOT NULL,
    total_price DECIMAL(18, 2) NOT NULL,
    status BOOLEAN DEFAULT false NOT NULL
);

CREATE TABLE categories (
    id SERIAL NOT NULL PRIMARY KEY,
    name VARCHAR NOT NULL,
    image_url varchar NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
    id SERIAL NOT NULL PRIMARY KEY,
    category_id INT NOT NULL REFERENCES categories(id),
    name VARCHAR NOT NULL,
    price DECIMAL(18, 2) NOT NULL,
    image_url VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE product_images (
    id SERIAL NOT NULL PRIMARY KEY,
    image_url VARCHAR NOT NULL,
    sequence_number INT NOT NULL,
    product_id INT NOT NULL REFERENCES products(id)
);

SELECT * FROM customer;
SELECT * FROM orders;
SELECT * FROM order_items;
SELECT * FROM categories;
SELECT * FROM products;
SELECT * FROM product_images;

DROP TABLE product_images;
DROP TABLE products;
DROP TABLE categories;
DROP TABLE order_items;
DROP TABLE customer;
DROP TABLE orders;