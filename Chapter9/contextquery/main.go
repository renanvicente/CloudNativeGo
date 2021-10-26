package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // Anonymously import the driver package
	"log"
	"net/http"
	"strconv"
	"time"
)

var db *sql.DB

type PostgresDBParams struct {
	dbName				string
	host				string
	user				string
	password			string
}

func  verifyTableExists() (bool,error) {
	//const table = "transactions"
	var table = "users"
	var result string
	rows, err := db.Query(fmt.Sprintf("SELECT to_regclass('public.%s');", table))
	defer rows.Close()
	if err != nil {
		return false, err
	}

	for rows.Next() && result != table {
		rows.Scan(&result)
	}
	return result == table, rows.Err()
}

func  createTable() error {
	var err error
	var table = "users"
	createQuery := fmt.Sprintf(`CREATE TABLE %s (
		id      BIGSERIAL PRIMARY KEY,
		username 		  TEXT
	  );`,table)
	_, err = db.Exec(createQuery)
	if err != nil {
		return err
	}
	return nil

}

func insertSampleUser() error {
	var err error
	query := `INSERT INTO users
            (username)
            VALUES ('testing')`
	_, err = db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func UserName(ctx context.Context, id int) (string, error) {
	const query = "SELECT username FROM users WHERE id=$1"

	dctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var username string
	err := db.QueryRowContext(dctx, query, id).Scan(&username)

	return username,err
}

func UserGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Get the request's context. This context is canceled when
	// the client's connection closes, the request is canceled
	// (with HTTP/2), or when the ServeHTTP method returns.
	rctx := r.Context()

	ctx, cancel := context.WithTimeout(rctx, 10*time.Second)
	defer cancel()
	idi, err := strconv.Atoi(id)
	username, err := UserName(ctx, idi)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		http.Error(w, "no such user", http.StatusNotFound)
	case errors.Is(err, context.DeadlineExceeded):
		http.Error(w, "database timeout", http.StatusGatewayTimeout)
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	default:
		w.Write([]byte(username))
	}
}

func main() {
	config := PostgresDBParams{
		host: "127.0.0.1",
		dbName: "kvs",
		user: "test",
		password: "kvstest",
	}
	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable",
		config.host, config.dbName, config.user,config.password)
	db1, err := sql.Open("postgres",connStr)
	if err != nil {
		fmt.Errorf("failed to open db: %w", err)
	}
	err = db1.Ping()		// Test the database connection
	db = db1
	exists, err := verifyTableExists()
	if err != nil {
		fmt.Errorf("failed to verify table exists: %w", err)
	}
	if !exists {
		if err = createTable(); err != nil {
			fmt.Errorf("failed to create table: %w", err)
		}
		insertSampleUser()

	}

	r := mux.NewRouter()
	r.HandleFunc("/user/{id}",UserGetHandler )
	log.Fatal(http.ListenAndServe(":8080", r))


}