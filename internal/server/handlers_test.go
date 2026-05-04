package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func TestCreateAccountHandler(t *testing.T) {
	db := setupTestDB()
	api := &API{DB: db}

	payload := []byte(`{"document_number": "123456789"}`)

	req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(payload))

	rr := executeRequest(req, api.CreateAccount)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response Account
	json.NewDecoder(rr.Body).Decode(&response)
	if response.DocumentNumber != "123456789" {
		t.Errorf("handler returned unexpected body: got %v", response.DocumentNumber)
	}

	t.Run("Failure - Duplicate Document Number", func(t *testing.T) {
		payload := []byte(`{"document_number": "12345"}`)
		http.Post("/accounts", "application/json", bytes.NewBuffer(payload))

		db.Exec("INSERT INTO accounts (document_number) VALUES ('12345')")

		req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(payload))

		rr := executeRequest(req, api.CreateAccount)

		if rr.Code != http.StatusConflict {
			t.Errorf("expected status 409, got %d", rr.Code)
		}
	})
}

func TestGetAccountHandler(t *testing.T) {
	db := setupTestDB()
	db.MustExec("INSERT INTO accounts (account_id, document_number) VALUES (1, '12345')")
	api := &API{DB: db}

	t.Run("Success - Account Found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/accounts/1", nil)
		req.SetPathValue("accountId", "1")

		rr := executeRequest(req, api.GetAccount)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rr.Code)
		}

		var acc Account
		json.NewDecoder(rr.Body).Decode(&acc)
		if acc.DocumentNumber != "12345" {
			t.Errorf("expected document 12345, got %s", acc.DocumentNumber)
		}
	})

	t.Run("Failure - Account Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/accounts/99", nil)
		req.SetPathValue("accountId", "99")

		rr := executeRequest(req, api.GetAccount)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rr.Code)
		}
	})
}

func TestCreateTransactionHandler(t *testing.T) {
	db := setupTestDB()
	db.MustExec("INSERT INTO accounts (account_id, document_number) VALUES (1, '123')")
	api := &API{DB: db}

	tests := []struct {
		name           string
		payload        Transaction
		expectedStatus int
	}{
		{
			name: "Success - Credit Voucher",
			payload: Transaction{
				AccountID:       1,
				OperationTypeID: OpCreditVoucher,
				Amount:          100.0,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Failure - Negative Purchase",
			payload: Transaction{
				AccountID:       1,
				OperationTypeID: OpNormalPurchase,
				Amount:          50.0, // Should be negative
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Failure - Account Does Not Exist",
			payload: Transaction{
				AccountID:       999, // FK constraint will trigger
				OperationTypeID: OpNormalPurchase,
				Amount:          -50.0,
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))

			rr := executeRequest(req, api.CreateTransaction)

			if rr.Code != tc.expectedStatus {
				t.Errorf("%s: expected status %d, got %d", tc.name, tc.expectedStatus, rr.Code)
			}
		})
	}
}

func setupTestDB() *sqlx.DB {
	db, _ := sqlx.Connect("sqlite", ":memory:")
	db.MustExec("PRAGMA foreign_keys = ON;")
	db.MustExec(Schema)
	db.MustExec(OperationTypesSeedQuery)
	return db
}

func executeRequest(req *http.Request, handler func(http.ResponseWriter, *http.Request)) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	h := http.HandlerFunc(handler)
	h.ServeHTTP(rr, req)

	return rr
}
