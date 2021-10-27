package transact

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Anonymously import the driver package
	"github.com/renanvicente/CloudNativeGo/Chapter10/hexarch/core"
)

type PostgresTransactionLogger struct {
	events chan<- core.Event // Write-only channel for sending events
	errors <-chan error      // Read-only channel for receiving errors
	db     *sql.DB           // The database access interface
	params PostgresDBParams
}

type PostgresDBParams struct {
	dbName           string
	host             string
	user             string
	password         string
	transactionTable string
}

func (l *PostgresTransactionLogger) WriteDelete(key string) {
	l.events <- core.Event{
		EventType: core.EventDelete,
		Key:       key,
	}
}

func (l *PostgresTransactionLogger) WritePut(key, value string) {
	l.events <- core.Event{
		EventType: core.EventPut,
		Key:       key,
		Value:     value,
	}
}

func (l *PostgresTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *PostgresTransactionLogger) ReadEvents() (<-chan core.Event, <-chan error) {
	outEvent := make(chan core.Event) // An unbuffered events channel
	outError := make(chan error, 1)   // A buffered errors channel
	var table = l.params.transactionTable

	go func() {
		defer close(outEvent) // Close the channels when the
		defer close(outError) // goroutine ends

		query := fmt.Sprintf(`SELECT sequence, event_type, key, value FROM %s
                  ORDER BY sequence`, table)

		rows, err := l.db.Query(query) // Run query; get result set
		if err != nil {
			outError <- fmt.Errorf("sql query error: %w", err)
			return
		}
		defer rows.Close() // This is important!

		e := core.Event{} // Create an empty Event

		for rows.Next() { // Iterate over the rows
			err = rows.Scan( // Read the values from the
				&e.Sequence, &e.EventType, &e.Key, &e.Value) // row into the Event.
			if err != nil {
				outError <- fmt.Errorf("error reading row: %w", e)
				return
			}
			outEvent <- e // Send e to the channel
		}
		err = rows.Err()
		if err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
		}
	}()
	return outEvent, outError
}

func (l *PostgresTransactionLogger) Run() {
	events := make(chan core.Event, 16) // Make an events channel
	l.events = events

	errors := make(chan error, 1) // Make an errors channel
	l.errors = errors

	go func() {
		query := `INSERT INTO transactions
            (event_type, key, value)
            VALUES ($1, $2, $3)`

		for e := range events { // Retrieve the next Event
			_, err := l.db.Exec( // Execute the INSERT query
				query,
				e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
			}
		}

	}()
}

func (l *PostgresTransactionLogger) verifyTableExists() (bool, error) {
	//const table = "transactions"
	var table = l.params.transactionTable
	var result string
	rows, err := l.db.Query(fmt.Sprintf("SELECT to_regclass('public.%s');", table))
	defer rows.Close()
	if err != nil {
		return false, err
	}

	for rows.Next() && result != table {
		rows.Scan(&result)
	}
	return result == table, rows.Err()
}

func (l *PostgresTransactionLogger) createTable() error {
	var err error
	var table = l.params.transactionTable
	createQuery := fmt.Sprintf(`CREATE TABLE %s (
		sequence      BIGSERIAL PRIMARY KEY,
		event_type    SMALLINT,
		key 		  TEXT,
		value         TEXT
	  );`, table)
	_, err = l.db.Exec(createQuery)
	if err != nil {
		return err
	}
	return nil

}

func NewPostgresTransactionLogger(config PostgresDBParams) (core.TransactionLogger, error) {
	if config.transactionTable == "" {
		config.transactionTable = "transactions"
	}
	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable",
		config.host, config.dbName, config.user, config.password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}
	err = db.Ping() // Test the database connection
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}
	logger := &PostgresTransactionLogger{db: db, params: config}
	exists, err := logger.verifyTableExists()
	if err != nil {
		return nil, fmt.Errorf("failed to verify table exists: %w", err)
	}
	if !exists {
		if err = logger.createTable(); err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}
	return logger, nil
}
