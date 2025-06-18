package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./company.db")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS employees (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			position TEXT NOT NULL,
			salary REAL NOT NULL
		);`)
	if err != nil {
		return nil, err
	}

	// Insert initial data if table is empty
	row := db.QueryRow("SELECT COUNT(*) FROM employees")
	var count int
	if err := row.Scan(&count); err != nil {
		return nil, err
	}
	if count == 0 {
		_, err = db.Exec(
			`INSERT INTO employees (name, position, salary) VALUES
			('Alice', 'Developer', 60000),
			('Bob', 'Manager', 80000),
			('Charlie', 'Designer', 70000);`)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func main() {
	db, err := CreateDatabase()
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/employees", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, position, salary FROM employees")
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		type Employee struct {
			ID       int     `json:"id"`
			Name     string  `json:"name"`
			Position string  `json:"position"`
			Salary   float64 `json:"salary"`
		}
		var employees []Employee
		for rows.Next() {
			var emp Employee
			if err := rows.Scan(&emp.ID, &emp.Name, &emp.Position, &emp.Salary); err != nil {
				http.Error(w, "Scan error", http.StatusInternalServerError)
				return
			}
			employees = append(employees, emp)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(employees)
	})

	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
