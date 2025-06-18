package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDatabase() {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			age INTEGER NOT NULL,
			password TEXT NOT NULL
		);`)
	if err != nil {
		panic(err)
	}

	mockUsers := []struct {
		name     string
		email    string
		age      int
		password string
	}{
		{"John Doe", "doe@gmail.com", 30, "password123"},
		{"Alice Smith", "alice.smith@example.com", 25, "alicepass"},
		{"Bob Johnson", "bob.johnson@example.com", 28, "bobpass"},
		{"Carol White", "carol.white@example.com", 32, "carolpass"},
		{"David Brown", "david.brown@example.com", 27, "davidpass"},
		{"Eve Black", "eve.black@example.com", 29, "evepass"},
		{"Frank Green", "frank.green@example.com", 31, "frankpass"},
		{"Grace Lee", "grace.lee@example.com", 26, "gracepass"},
		{"Henry King", "henry.king@example.com", 33, "henrypass"},
		{"Ivy Scott", "ivy.scott@example.com", 24, "ivypass"},
		{"Jack Turner", "jack.turner@example.com", 35, "jackpass"},
		{"Kathy Adams", "kathy.adams@example.com", 28, "kathypass"},
		{"Leo Clark", "leo.clark@example.com", 30, "leopass"},
		{"Mona Lewis", "mona.lewis@example.com", 27, "monapass"},
		{"Nina Walker", "nina.walker@example.com", 29, "ninapass"},
		{"Oscar Hall", "oscar.hall@example.com", 31, "oscarpass"},
		{"Pam Young", "pam.young@example.com", 26, "pampass"},
		{"Quinn Allen", "quinn.allen@example.com", 32, "quinnpass"},
		{"Rita Wright", "rita.wright@example.com", 25, "ritapass"},
		{"Sam Hill", "sam.hill@example.com", 34, "sampass"},
		{"Tina Baker", "tina.baker@example.com", 28, "tinapass"},
		{"Uma Nelson", "uma.nelson@example.com", 30, "umapass"},
		{"Vince Carter", "vince.carter@example.com", 27, "vincepass"},
		{"Wendy Mitchell", "wendy.mitchell@example.com", 29, "wendypass"},
		{"Xander Perez", "xander.perez@example.com", 31, "xanderpass"},
		{"Yara Roberts", "yara.roberts@example.com", 26, "yarapass"},
		{"Zane Evans", "zane.evans@example.com", 33, "zanepass"},
		{"Amy Foster", "amy.foster@example.com", 24, "amypass"},
		{"Ben Graham", "ben.graham@example.com", 35, "benpass"},
		{"Cindy Hughes", "cindy.hughes@example.com", 28, "cindypass"},
	}
	for _, user := range mockUsers {
		_, err = db.Exec("INSERT OR IGNORE INTO users (name, email, age, password) VALUES (?, ?, ?, ?)",
			user.name, user.email, user.age, user.password)
		if err != nil {
			panic(err)
		}
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
func main() {
	CreateDatabase()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server is running"))
	})
	http.HandleFunc("/api/users", GetUsers)
	http.ListenAndServe(":8080", nil)
}
