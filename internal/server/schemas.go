package server

const Schema = `
CREATE TABLE IF NOT EXISTS accounts (
    account_id INTEGER PRIMARY KEY AUTOINCREMENT,
    document_number TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS operation_types (
    operation_type_id INTEGER PRIMARY KEY,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    transaction_id INTEGER PRIMARY KEY AUTOINCREMENT,
    account_id INTEGER NOT NULL,
    operation_type_id INTEGER NOT NULL,
    amount REAL NOT NULL,
    event_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES accounts(account_id)
    FOREIGN KEY (operation_type_id) REFERENCES operation_types(operation_type_id)
);
`

const OperationTypesSeedQuery = `
INSERT OR IGNORE INTO operation_types (operation_type_id, description) VALUES 
(1, 'NORMAL PURCHASE'),
(2, 'PURCHASE WITH INSTALLMENTS'),
(3, 'WITHDRAWAL'),
(4, 'CREDIT VOUCHER');
`
