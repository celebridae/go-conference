package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
	// _ "github.com/mattn/go-sqlite3"
	// _ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var db *sql.DB

func main() {
	// println("Ola Mundo")
	//TODO: register routes
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", indexHandler)
	mux.HandleFunc("/users", listUserHandler)
	mux.HandleFunc("/users/save", insertHandler)
	//TODO: port default
	http.ListenAndServe(":8080", mux)
}

// TODO: func to list users
func listUserHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize database connection if not already initialized
	if db == nil {
		err := InitializeDB()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Error encoding json", http.StatusInternalServerError)
		return
	}
}

// InitializeDB initializes the database connection
func InitializeDB() error {
	var err error
	// Open database connection
	db, err = sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/goconference_db?sslmode=disable")
	if err != nil {
		return err
	}
	// Ensure the users table exists
	err = createUsersTable()
	if err != nil {
		return err
	}
	return nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the index page!")
}

// createUsersTable creates the users table if it does not exist
func createUsersTable() error {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL  PRIMARY KEY,
			name TEXT,
			email TEXT
		);
	`
	_, err := db.Exec(createTableSQL)
	return err
}

// Handler function for inserting data
func insertHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		err := InitializeDB()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// Parse request body
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		fmt.Println("Failed to parse request body:", err)
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	// Insert data into database
	_, err = db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", newUser.Name, newUser.Email)
	if err != nil {
		http.Error(w, "Failed to insert data into database", http.StatusInternalServerError)
		return
	}
	// Respond with success message
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Data inserted successfully")
}

// Handler function for finding a user by ID
func findByIdHandler(w http.ResponseWriter, r *http.Request) {
	// Parse user ID from request URL
	userID, err := strconv.ParseInt(r.URL.Path[len("/users/"):], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Query the database to retrieve the user by ID
	var user User
	err = db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Name, &user.Email)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		return
	}

	// Serialize the user data into JSON format
	jsonData, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Failed to serialize user data", http.StatusInternalServerError)
		return
	}

	// Write the JSON response back to the client
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
