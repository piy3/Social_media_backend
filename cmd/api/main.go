package main

import (
	"log"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/piy3/social/internal/db"
	"github.com/piy3/social/internal/env"
	"github.com/piy3/social/internal/mailer"
	"github.com/piy3/social/internal/store"
)

const version="0.0.1"

func main() { 
	// Load .env file from project root
	godotenv.Load(filepath.Join("..", "..", ".env"))

	cfg := config{
		addr: env.GetString("ADDR", "localhost:8081"),
		frontendURL: env.GetString("FRONTEND_URL","http://localhost:4000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://piyush:root123@localhost:5433/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env:env.GetString("ENV","development"),
		mail: mailConfig{
			exp: time.Hour *24*3, // 3 days
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
				
			},
			fromEmail: env.GetString("SENDGRID_FROM_EMAIL", ""),
		},
	}
	db,err:=db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Connected to database successfully")
	defer db.Close()
	store := store.NewStorage(db)
	mailer:= mailer.NewSendgrid(cfg.mail.sendGrid.apiKey,cfg.mail.fromEmail)

	app := &application{
		config: cfg,
		store:store,
		mailer: mailer,
	}


	mux := app.mount()
	log.Fatal(app.run(mux))
}
//new cmt