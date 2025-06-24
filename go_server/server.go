package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Password string `json:"password"`
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetUsers handler called")
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, email, age FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Age); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}
	fmt.Printf("Fetched users: %+v\n", users)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func emailExists(db *sql.DB, email string, excludeID int) (bool, error) {
	var id int
	query := "SELECT id FROM users WHERE email = ?"
	args := []interface{}{email}
	if excludeID > 0 {
		query += " AND id != ?"
		args = append(args, excludeID)
	}
	if err := db.QueryRow(query, args...).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateUser handler called")
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("Inserting user: %+v\n", user)

	exists, err := emailExists(db, user.Email, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "email already exists", http.StatusBadRequest)
		return
	}

	if !isValidEmail(user.Email) {
		http.Error(w, "invalid email format", http.StatusBadRequest)
		return
	}
	if user.Name == "" || user.Email == "" || user.Age <= 0 || user.Password == "" {
		http.Error(w, "all fields are required and age must be positive", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO users(name, email, age, password) VALUES(?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.Name, user.Email, user.Age, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	fmt.Fprintf(w, "User created with ID: %d", id)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UpdateUser handler called")
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exists, err := emailExists(db, user.Email, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "email already exists", http.StatusBadRequest)
		return
	}

	if !isValidEmail(user.Email) {
		http.Error(w, "invalid email format", http.StatusBadRequest)
		return
	}
	if user.Name == "" || user.Email == "" || user.Age <= 0 || user.Password == "" {
		http.Error(w, "all fields are required and age must be positive", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("UPDATE users SET name=?, email=?, age=?, password=? WHERE id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.Name, user.Email, user.Age, user.Password, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	fmt.Fprintf(w, "User updated: %d rows affected", rowsAffected)
}

func DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeleteUserByID handler called")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM users WHERE id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	if _, err := stmt.Exec(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func main() {
	fmt.Println("Starting server on :8080")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Root handler called")
		w.Write([]byte("server is running"))
	})

	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("/api/users handler called, method: %s\n", r.Method)
		switch r.Method {
		case http.MethodGet:
			GetUsers(w, r)
		case http.MethodPost:
			CreateUser(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("/api/users/ handler called, method: %s\n", r.Method)
		switch r.Method {
		case http.MethodPut:
			UpdateUser(w, r)
		case http.MethodDelete:
			DeleteUserByID(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", nil)
}
