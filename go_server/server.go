package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"regexp"

	_ "github.com/mattn/go-sqlite3"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetUsers handler called")
	PrintAllUsers()
	db, err := sql.Open(("sqlite3"), "./data.db")
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

	type User struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Age      int    `json:"age"`
		Password string `json:"password"`
	}
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
func CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateUser handler called")
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	type User struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Age      int    `json:"age"`
		Password string `json:"password"`
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("Inserting user: %+v\n", user)

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
	if !isValidEmail(user.Email) {
		http.Error(w, "invalid email format", http.StatusBadRequest)
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

	type User struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Age      int    `json:"age"`
		Password string `json:"password"`
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	if !isValidEmail(user.Email) {
    http.Error(w, "invalid email format", http.StatusBadRequest)
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

func PrintAllUsers() {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		fmt.Println("Error opening DB:", err)
		return
	}
	defer db.Close()
	rows, err := db.Query("SELECT id, name, email, age, password FROM users")
	if err != nil {
		fmt.Println("Error querying users:", err)
		return
	}
	defer rows.Close()
	fmt.Println("All users in DB:")
	for rows.Next() {
		var id int
		var name, email, password string
		var age int
		if err := rows.Scan(&id, &name, &email, &age, &password); err != nil {
			fmt.Println("Error scanning user:", err)
			return
		}
		fmt.Printf("ID: %d, Name: %s, Email: %s, Age: %d, Password: %s\n", id, name, email, age, password)
	}
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func main() {
	fmt.Println("Starting server on :8080")
	PrintAllUsers()
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
