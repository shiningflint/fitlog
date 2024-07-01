package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
	"os"
)

var db *sql.DB

func GetConnection() *sql.DB {
	if db != nil {
		return db
	}

	var dbName string

	if os.Getenv("APP_ENV") == "testing" {
		dbName = "app_test.db"
	} else {
		dbName = "app.db"
	}

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		fmt.Printf("🔥 failed to connect to the database: %s", err.Error())
	}

	fmt.Println("🚀 Connected Successfully to the Database")

	return db
}

// Users struct to hold retrieved data
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func usersIndexHandler(w http.ResponseWriter, r *http.Request) {
	db = GetConnection()

	// Define the SQL query
	rows, err := db.Query("SELECT id, name, email FROM users ORDER BY created_at DESC")
	if err != nil {
		panic(err)
	}

	defer rows.Close() // Close the rows after iterating

	// Scan the results
	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}

	// Check for any errors during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	tmpl := template.Must(template.ParseFiles("templates/users/index.html", "templates/users/_newButton.html", "templates/users/_item.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "index.html", users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Handler function for POST requests on /admin/users
func usersCreateHandler(w http.ResponseWriter, r *http.Request) {
	db = GetConnection()

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	name := r.Form.Get("name")
	email := r.Form.Get("email")
	user := User{
		Name:  name,
		Email: email,
	}

	// Insert user into database
	stmt, err := db.Prepare("INSERT INTO users (name, email) VALUES (?, ?)")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error preparing statement: %v", err)
		return
	}
	defer stmt.Close() // Close prepared statement on exit

	// Execute insert statement
	_, err = stmt.Exec(user.Name, user.Email) // Replace with hashedPassword if implemented
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error inserting user: %v", err)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/users/create.html", "templates/users/_newButton.html", "templates/users/_item.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "create.html", user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func usersEditHanlder(w http.ResponseWriter, r *http.Request) {
	db = GetConnection()

	userID := r.PathValue("id")
	var user User

	err := db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		panic(err)
	}

	tmpl := template.Must(template.ParseFiles("templates/users/edit.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tmpl.ExecuteTemplate(w, "edit.html", user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func usersUpdateHandler(w http.ResponseWriter, r *http.Request) {
	db = GetConnection()

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	id := r.Form.Get("id")
	name := r.Form.Get("name")
	email := r.Form.Get("email")
	user := User{
		ID:    id,
		Name:  name,
		Email: email,
	}

	stmt, err := db.Prepare("UPDATE users SET name = ?, email = ? WHERE id = ?")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error preparing statement: %v", err)
		return
	}
	defer stmt.Close() // Close prepared statement on exit

	// Execute insert statement
	_, err = stmt.Exec(user.Name, user.Email, user.ID) // Replace with hashedPassword if implemented
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error inserting user: %v", err)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/users/_item.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "_item.html", user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func usersDeleteHandler(w http.ResponseWriter, r *http.Request) {
	db = GetConnection()

	userID := r.PathValue("id")

	stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error preparing statement: %v", err)
		return
	}
	defer stmt.Close() // Close prepared statement on exit

	// Execute insert statement
	_, err = stmt.Exec(userID) // Replace with hashedPassword if implemented
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting user: %v", err)
		return
	}
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	// db = GetConnection()

	data := struct{}{} // empty data for now

	tmpl := template.Must(template.ParseFiles("templates/me/index.html"))

	if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func serveStatic() {
	// Define the static directory
	assetsDir := http.Dir("assets")

	// Create a file server for the static directory (excluding index.html)
	fileServer := http.StripPrefix("/assets/", http.FileServer(assetsDir))

	// Serve static files under the "/assets/" prefix
	http.Handle("/assets/", fileServer)
}

func main() {
	os.Setenv("APP_ENV", "production")

	serveStatic()

	// New feature from go 1.22, using NewServeMux will enable you to do URL patterns such as: "/path/{id}/edit"
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})
	mux.HandleFunc("GET /admin/users", usersIndexHandler)
	mux.HandleFunc("POST /admin/users", usersCreateHandler)
	mux.HandleFunc("GET /admin/users/new", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/users/new.html")
	})
	mux.HandleFunc("GET /admin/users/cancel", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/users/_newButton.html")
	})
	mux.HandleFunc("GET /admin/users/{id}/edit", usersEditHanlder)
	mux.HandleFunc("PATCH /admin/users/{id}", usersUpdateHandler)
	mux.HandleFunc("DELETE /admin/users/{id}", usersDeleteHandler)
	mux.HandleFunc("GET /me", meHandler)

	http.ListenAndServe(":5000", mux)
}
