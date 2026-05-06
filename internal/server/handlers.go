package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/DouglasCI/pismo-tech-case-backend/internal/domain"
)

type API struct {
	DB *sqlx.DB
}

func (api *API) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var newAccount domain.Account

	if err := json.NewDecoder(r.Body).Decode(&newAccount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := newAccount.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	query := `INSERT INTO accounts (document_number) VALUES (?) RETURNING account_id`
	if err := api.DB.QueryRow(query, newAccount.DocumentNumber).Scan(&newAccount.ID); err != nil {
		http.Error(w, "Failed to create account (document number might already exist)", http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newAccount)
}

func (api *API) GetAccount(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("accountId")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid account ID format", http.StatusBadRequest)
		return
	}

	var account domain.Account
	query := `SELECT account_id, document_number FROM accounts WHERE account_id = ?`
	if err := api.DB.Get(&account, query, id); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

func (api *API) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var newTransaction domain.Transaction

	if err := json.NewDecoder(r.Body).Decode(&newTransaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := newTransaction.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	newTransaction.EventDate = time.Now().Format(time.RFC3339Nano)

	query := `
		INSERT INTO transactions (account_id, operation_type_id, amount, event_date) 
		VALUES (?, ?, ?, ?) 
		RETURNING transaction_id
	`
	if err := api.DB.QueryRow(
		query,
		newTransaction.AccountID,
		newTransaction.OperationTypeID,
		newTransaction.Amount,
		newTransaction.EventDate,
	).Scan(&newTransaction.ID); err != nil {
		http.Error(w, "Failed to create transaction (check if account exists)", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTransaction)
}
