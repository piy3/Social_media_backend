package main

import (
	"log"

	"github.com/piy3/social/internal/db"
	"github.com/piy3/social/internal/env"
	"github.com/piy3/social/internal/store"
)

func main() {
	addr:=env.GetString("DB_ADDR","postgres://piyush:root123@localhost:5433/socialnetwork?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err) 
	}
	
	defer conn.Close()

	store := store.NewStorage(conn)
	db.Seed(store,	conn)
}