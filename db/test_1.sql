create TABLE users (
    id INT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(50) NOT NULL UNIQUE
);

INSERT INTO users (id, name, email) VALUES
(1, 'Alice Smith', 'alice@example.com'),
(2, 'Bob Johnson', 'bob@example.com'),
(3, 'Charlie Lee', 'charlie@example.com'),
(4, 'Diana Prince', 'diana@example.com'),
(5, 'Ethan Hunt', 'ethan@example.com'),
(6, 'Fiona Gallagher', 'fiona@example.com'),
(7, 'George Costanza', 'george@example.com');
