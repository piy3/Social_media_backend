package main

import (
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/piy3/social/internal/env"
)

func main() {
	// Load .env file from project root
	godotenv.Load(filepath.Join("..", "..", ".env"))

	cfg := config{
		addr: env.GetString("ADDR", "localhost:8081"),
	}
	app := &application{
		config: cfg,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
