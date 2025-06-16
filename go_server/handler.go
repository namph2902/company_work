package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type User struct {
	ID	int    `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

func GetUsersHandler(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, email FROM users")
		  if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			   return
	}  
	  defer rows.Close()

	  var users []User
	  for rows.Next() {
		var u User
		  if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			   http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	  return
	}
	users = append(users, u)
}
    w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
}