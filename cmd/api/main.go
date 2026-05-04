package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite" // Inject driver as anonymous

	"github.com/DouglasCI/pismo-tech-case-backend/internal/server"
)

func main() {
	// Connect to SQLite
	db, err := sqlx.Connect("sqlite", "transactions.db")
	if err != nil {
		log.Fatalln(err)
	}

	// Close DB connection after clean up
	defer db.Close()

	// Initialize DB tables with fail fast
	db.MustExec("PRAGMA foreign_keys = ON;") // Forces SQLite to respect Foreign Key constraints
	db.MustExec(server.Schema)
	db.MustExec(server.OperationTypesSeedQuery)

	// Initialize API struct with the DB connection
	api := &server.API{DB: db}

	// Initialize router
	mux := http.NewServeMux()

	// Endpoints
	mux.HandleFunc("POST /accounts", api.CreateAccount)
	mux.HandleFunc("GET /accounts/{accountId}", api.GetAccount)
	mux.HandleFunc("POST /transactions", api.CreateTransaction)

	// Start server
	fmt.Println("Server listening in port 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalln(err)
	}
}
