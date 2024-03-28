package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	// println("Ola Mundo")
	//TODO: register routes
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", indexHandler)
	mux.HandleFunc("/users", listUserHandler)
	//TODO: port default
	http.ListenAndServe(":8080", mux)
}

// TODO: func to list users
func listUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "users.db")
	fmt.Fprintf(w, "Welcome 1!")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Welcome 2!")
	defer db.Close() // TODO: close connection

	// SQL statement to create table
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			email TEXT
		);
		`

	// Execute SQL statement
	_, err = db.Exec(createTableSQL)
	if err != nil {
		//log.Fatalf("Error creating table: %v\n", err)
		fmt.Fprintf(w, "Welcome to the index page!%v\n", err)
	}

	fmt.Println("Table created successfully!")

	fmt.Println(w, "Welcome 3!")
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		// Log the error for debugging
		fmt.Println("Error executing SQL query:", err)

		// Return an HTTP 500 Internal Server Error response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// fmt.Fprintf(w, "Welcome 4!")
	fmt.Println("Table 3!")
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, &u)
	}

	fmt.Println("Table 4!")

	// TODO: transform users

	// w.Header().Set("content-type", "application/json")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Error encoding json", http.StatusInternalServerError)
		return
	}

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the index page!")
}
