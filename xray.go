package main

import (
	"context"
	"encoding/json/v2"
	"os"
	"text/template"
)

func InitXray() error {
	tmpl, err := template.ParseFS(tmpls, "tmpls/config.json.tmpl")

	if err != nil {
		return err
	}

	f, err := os.Create("/usr/local/etc/xray/config.json")

	if err != nil {
		return err
	}

	clients, err := FindClients(context.Background())

	if err != nil {
		return err
	}

	b, err := json.Marshal(clients)

	if err != nil {
		return err
	}

	return tmpl.Execute(f, string(b))
}
