package config

import (
	"context"
	"os"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v4"
)

func InitPgSQL() (context.Context, *pgx.ConnConfig) {
	ctx := context.Background()
	connstring := "postgresql://"

	if env := os.Getenv("POSTGRES_USER"); env == "" {
		log.Fatal("Bad 'POSTGRES_USER' parameter env")
	} else {
		connstring += env
	}

	if env := os.Getenv("POSTGRES_PASSWORD"); env == "" {
		log.Warn("'POSTGRES_PASSWORD' not set parameter env")
	} else {
		connstring += ":" + env

	}

	if env := os.Getenv("POSTGRES_HOST"); env == "" {
		log.Fatal("Bad 'POSTGRES_HOST' parameter env")
		os.Exit(1)
	} else {
		connstring += "@" + env
	}

	if env := os.Getenv("POSTGRES_DB"); env == "" {
		log.Fatal("Bad 'POSTGRES_DB' parameter env")
	} else {
		connstring += "/" + env
	}

	connstring += "?sslmode=disable"

	connConf, err := pgx.ParseConfig(connstring)
	if err != nil {
		log.Fatalf("Parse error : %s", err)
	}

	sqlCo, err := pgx.ConnectConfig(ctx, connConf)
	if err != nil {
		log.Fatalf("error connect psql : %s", err)
		return ctx, connConf
	}
	defer sqlCo.Close(ctx)

	query := `
	CREATE EXTENSION IF NOT EXISTS pgcrypto;

	CREATE TABLE IF NOT EXISTS account (
		id 							SERIAL,
		email 						TEXT NOT NULL UNIQUE,
		firstname 					TEXT NOT NULL,
		lastname 					TEXT NOT NULL,
		password 					TEXT NOT NULL,
		recover_token 				TEXT,
		admin 						boolean DEFAULT FALSE,
		enable						boolean DEFAULT TRUE,
		PRIMARY KEY(id)
	);

	CREATE TABLE IF NOT EXISTS chat_session (
		id 							SERIAL,
		user_id 					INTEGER NOT NULL REFERENCES account(id) ON DELETE CASCADE,
		title 						TEXT,
		created_at 					TIMESTAMP DEFAULT NOW(),
		updated_at 					TIMESTAMP DEFAULT NOW(),
		PRIMARY KEY(id),
		FOREIGN KEY(user_id) 		REFERENCES account(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS chat_messages (
		id 							SERIAL,
		chat_session_id				INTEGER NOT NULL,
		user_message				TEXT NOT NULL,
		bot_message					TEXT NOT NULL,
		created_at 					TIMESTAMP DEFAULT NOW(),
		PRIMARY KEY(id),
		FOREIGN KEY(chat_session_id) REFERENCES chat_session(id) ON DELETE CASCADE
	);
	`

	_, err = sqlCo.Exec(ctx, query)
	if err != nil {
		log.Fatal("During postgresql setup", "query", query, "error", err)
	}
	return ctx, connConf
}
