package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"machin8/core/conf"
)

func main() {
	var username, password string
	flag.StringVar(&username, "username", "", "username")
	flag.StringVar(&username, "u", "", "username (shorthand)")
	flag.StringVar(&password, "password", "", "password")
	flag.StringVar(&password, "p", "", "password (shorthand)")
	flag.Parse()

	if username == "" || password == "" {
		fmt.Fprintln(os.Stderr, "usage: create_user -u <username> -p <password>")
		os.Exit(1)
	}

	if err := conf.InitSettings(); err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := sql.Open("postgres", conf.Settings.DBUrl)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	var id int64
	err = db.QueryRow(
		`INSERT INTO auth.user (username, password, email, role)
		 VALUES ($1, crypt($2, gen_salt('bf')), $1 || '@example.com', 'analist')
		 RETURNING id`,
		username, password,
	).Scan(&id)
	if err != nil {
		log.Fatalf("create user: %v", err)
	}

	fmt.Printf("user created: id=%d username=%s\n", id, username)
}
