package main

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"text/template"
)

func InitCaddy() error {
	tmpl, err := template.ParseFS(tmpls, "tmpls/Caddyfile.tmpl")

	if err != nil {
		return err
	}

	var b bytes.Buffer

	config, err := FindConfig(context.Background())

	if err != nil {
		return err
	}

	err = tmpl.Execute(&b, config.Domain)

	if err != nil {
		return err
	}

	resp, err := http.Post("http://127.0.0.1:2019/load", "text/caddyfile", &b)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("caddy load failed")
	}

	return nil
}
