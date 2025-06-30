package util

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

// DeferRowsClose attempts to close a sql.Rows and logs an error message if fails
//
// This function is supposed to be deferred, like in the example below:
//
//	rows, err := db.Query(...)
//	if err != nil { ... }
//	defer util.DeferRowsClose(rows)
func DeferRowsClose(rows *sql.Rows) {
	err := rows.Close()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("[rows.Close] | %v", err)
	}
}

// DeferFileClose attempts to close an os.File and logs an error message if fails
//
// This function is supposed to be deferred, like in the example below:
//
//	file, err := os.Open(...)
//	if err != nil { ... }
//	defer util.DeferFileClose(file)
func DeferFileClose(file *os.File) {
	err := file.Close()
	if err != nil && !errors.Is(err, os.ErrClosed) {
		log.Printf("[file.Close] | %v", err)
	}
}

// DeferTxRollback attempts to roll back a sql.Tx and logs an error message if fails
//
// This function is supposed to be deferred, like in the example below:
//
//	tx, err := db.Begin()
//	if err != nil { ... }
//	defer util.DeferTxRollback(tx)
func DeferTxRollback(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil && !errors.Is(err, sql.ErrTxDone) {
		log.Printf("[tx.Rollback] | %v", err)
	}
}
