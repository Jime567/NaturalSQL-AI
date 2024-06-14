package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jime567/NaturalSQL-AI/structs"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// Credentials
	dbUser := os.Getenv("AI_SQL_USER")
	if dbUser == "" {
		log.Fatal("SQL_USER environment variable not set")
	}
	dbPass := os.Getenv("AI_SQL_PASS")
	if dbPass == "" {
		log.Fatal("SQL_PASS environment variable not set")
	}
	AIKEY := os.Getenv("OPENAI_KEY")
	if AIKEY == "" {
		log.Fatal("OPENAI_KEY environment variable not set")
	}
	AIORG := os.Getenv("OPENAI_ORG_ID")
	if AIORG == "" {
		log.Fatal("OPENAI_ORG_ID environment variable not set")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/cyclists_db", dbUser, dbPass)

	// Open a database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {

		log.Fatal("Unable to open connection: ", err)
	}
	defer db.Close()

	// Verify the connection
	if err := db.Ping(); err != nil {
		log.Fatal("Connection Failed: ", err)
	}

	fmt.Println("Successfully connected to the database!")

	// Register HTTP handlers
	http.HandleFunc("/getCyclists", func(w http.ResponseWriter, r *http.Request) {
		cyclists := getCyclists(db)

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Write cyclists data as JSON response
		if err := json.NewEncoder(w).Encode(cyclists); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Start the HTTP server
	port := ":8080"
	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func executeSQL(db *sql.DB, query string) {
	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}
}

func getBikes(db *sql.DB) string {
	var bikes string
	rows, err := db.Query("SELECT * FROM Bikes")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	// Iterate over the rows
	for rows.Next() {
		var id int
		var nickname, serialNumber, year, model, make string
		var mileage int

		if err := rows.Scan(&id, &nickname, &serialNumber, &year, &model, &make, &mileage); err != nil {
			log.Fatal(err)
		}
		bikes += fmt.Sprintf("ID: %d, Nickname: %s, Serial Number: %s, Year: %s, Model: %s, Make: %s, Mileage: %d\n", id, nickname, serialNumber, year, model, make, mileage)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return bikes
}

func getAddresses(db *sql.DB) string {
	var addresses string
	rows, err := db.Query("SELECT * FROM Addresses")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		var id int
		var street, zip, state string

		if err := rows.Scan(&id, &street, &zip, &state); err != nil {
			log.Fatal(err)
		}
		addresses += fmt.Sprintf("ID: %d, Street: %s, Zip: %s, State: %s\n", id, street, zip, state)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return addresses
}

func getCyclists(db *sql.DB) []structs.Cyclist {
	var cyclists []structs.Cyclist

	rows, err := db.Query("SELECT * FROM Cyclists")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		var cyclist structs.Cyclist
		if err := rows.Scan(&cyclist.ID, &cyclist.FirstName, &cyclist.LastName, &cyclist.Email, &cyclist.AddressID, &cyclist.BikeID); err != nil {
			log.Fatal(err)
		}
		cyclists = append(cyclists, cyclist)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return cyclists
}
