package main

import (
	"encoding/base64"
	"net/http"
	"strings"
	// "errors"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	// "net"
	// "os"
	// "os/signal"
	// "strconv"
	// "time"
	//"github.com/lib/pq"
)

var host string

func sql_open(u string, p string) (*sql.DB, error) {
	sql, err := sql.Open("postgres", fmt.Sprintf("host=127.0.0.1 user=%s sslmode=disable dbname=booktown", u, p))
	return sql, err
}

func auth(w http.ResponseWriter, r *http.Request) (string, string, bool) {
	if auth, ok := r.Header["Authorization"]; ok {
		b64Creds := strings.Split(auth[0], " ")[1]
		creds, _ := base64.StdEncoding.DecodeString(b64Creds)
		credStr := string(creds)
		credArray := strings.Split(credStr, ":")
		username := credArray[0]
		password := credArray[1]

		return username, password, true
	} else {
		return "", "", false
	}
}

func sql_query(db *sql.DB) {
	var username string
	id := 25
	err := db.QueryRow("SELECT username FROM users WHERE id=?", id).Scan(&username)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID.")
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Printf("Username is %s\n", username)
	}
}

func request_auth(w http.ResponseWriter) {
	log.Printf("requesting auth")
	w.Header().Set("WWW-Authenticate", "Basic realm=\"Postgresql\"")
	w.WriteHeader(401)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if username, password, ok := auth(w, r); ok {
		db, err := sql_open(username, password)
		if err != nil {
			log.Fatal(err)
			request_auth(w)
		} else {
			sql_query(db)
			fmt.Fprintf(w, "%s username: %s password: %s", "OK", username, password)
		}
	} else {
		request_auth(w)
	}

}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}
