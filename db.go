package main

import (
	"context"
	"database/sql"
)

func InitDB() error {
	var err error

	db, err = sql.Open("sqlite3", "/var/lib/nanami/db")

	if err != nil {
		return err
	}

	ctx := context.Background()

	_, err = db.ExecContext(
		ctx,
		"CREATE TABLE IF NOT EXISTS clients(id TEXT PRIMARY KEY, email TEXT UNIQUE)",
	)

	if err != nil {
		return err
	}

	_, err = db.ExecContext(
		ctx,
		"CREATE TABLE IF NOT EXISTS config(id INTEGER PRIMARY KEY, domain TEXT)",
	)

	if err != nil {
		return err
	}

	_, err = db.ExecContext(
		ctx,
		"CREATE TABLE IF NOT EXISTS sessions(token TEXT PRIMARY KEY, data BLOB NOT NULL, expiry REAL NOT NULL)",
	)

	if err != nil {
		return err
	}

	_, err = db.ExecContext(
		ctx,
		"CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions(expiry)",
	)

	if err != nil {
		return err
	}

	_, err = db.ExecContext(
		ctx,
		"CREATE TABLE IF NOT EXISTS users(name TEXT PRIMARY KEY, password TEXT)",
	)

	if err != nil {
		return err
	}

	return nil
}

type Client struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func FindClients(ctx context.Context) ([]Client, error) {
	clients := []Client{}

	rows, err := db.QueryContext(ctx, "SELECT id, email FROM clients")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var client Client

		if err := rows.Scan(&client.ID, &client.Email); err != nil {
			return nil, err
		}

		clients = append(clients, client)
	}

	return clients, nil
}

func NewClient(ctx context.Context, client Client) error {
	_, err := db.ExecContext(
		ctx,
		"INSERT INTO clients(id, email) VALUES($1, $2)",
		client.ID,
		client.Email,
	)

	if err != nil {
		return err
	}

	return nil
}

func RemoveClient(ctx context.Context, email string) error {
	_, err := db.ExecContext(ctx, "DELETE FROM clients WHERE email = $1", email)

	if err != nil {
		return err
	}

	return nil
}

type Config struct {
	Domain string
}

func FindConfig(ctx context.Context) (*Config, error) {
	config := &Config{}

	rows, err := db.QueryContext(ctx, "SELECT domain FROM config WHERE id = 1")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&config.Domain); err != nil {
			return nil, err
		}
	}

	return config, nil
}

func SaveConfig(ctx context.Context, config *Config) error {
	currentConfig, err := FindConfig(ctx)

	if err != nil {
		return err
	}

	if currentConfig.Domain == "" {
		_, err = db.ExecContext(ctx, "INSERT INTO config(domain) VALUES($1)", config.Domain)

		if err != nil {
			return err
		}

		return nil
	}

	_, err = db.ExecContext(ctx, "UPDATE config SET domain = $1 WHERE id = 1", config.Domain)

	if err != nil {
		return err
	}

	return nil
}

type User struct {
	Name     string
	Password string
}

func FindUser(ctx context.Context, name string) (*User, error) {
	user := &User{}

	rows, err := db.QueryContext(ctx, "SELECT name, password FROM users WHERE name = $1", name)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&user.Name, &user.Password); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func NewUser(ctx context.Context, user *User) error {
	_, err := db.ExecContext(
		ctx,
		"INSERT INTO users(name, password) VALUES($1, $2)",
		user.Name,
		user.Password,
	)

	if err != nil {
		return err
	}

	return nil
}
