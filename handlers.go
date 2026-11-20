package main

import (
	"html/template"
	"log"
	"net/http"
	"uuid"
)

func internalError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func render(path string, w http.ResponseWriter, data any) error {
	tmpl, err := template.ParseFS(tmpls, "tmpls/base.html.tmpl", path)

	if err != nil {
		return err
	}

	return tmpl.Execute(w, data)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "POST" {
		action := r.PostFormValue("action")

		if action == "create" {
			email := NewField(r, "email", []Validator{
				NewMinLength(2),
				NewMaxLength(16),
			})

			if !email.Valid() {
				sessionManager.Put(r.Context(), email.Name, email.Errors)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			client := Client{
				ID:    uuid.New().String(),
				Email: email.Value,
			}

			err := NewClient(r.Context(), client)

			if err != nil {
				log.Print(err)
				internalError(w)
				return
			}

			err = InitXray()

			if err != nil {
				log.Print(err)
				internalError(w)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else if action == "delete" {
			err := RemoveClient(r.Context(), r.PostFormValue("email"))

			if err != nil {
				log.Print(err)
				internalError(w)
				return
			}

			err = InitXray()

			if err != nil {
				log.Print(err)
				internalError(w)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			internalError(w)
			return
		}
	}

	clients, err := FindClients(r.Context())

	if err != nil {
		log.Print(err)
		internalError(w)
		return
	}

	err = render("tmpls/home.html.tmpl", w, map[string]any{
		"User":        sessionManager.Get(r.Context(), "user"),
		"Path":        r.URL.Path,
		"Title":       "Home",
		"Clients":     clients,
		"EmailErrors": sessionManager.Pop(r.Context(), "email"),
	})

	if err != nil {
		log.Print(err)
		internalError(w)
		return
	}
}

func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		domain := NewField(r, "domain", []Validator{
			NewMinLength(2),
			NewMaxLength(16),
		})

		if !domain.Valid() {
			sessionManager.Put(r.Context(), domain.Name, domain.Errors)
			http.Redirect(w, r, "/config", http.StatusSeeOther)
			return
		}

		config := &Config{
			Domain: domain.Value,
		}

		err := SaveConfig(r.Context(), config)

		if err != nil {
			log.Print(err)
			internalError(w)
			return
		}

		err = InitCaddy()

		if err != nil {
			log.Print(err)
			internalError(w)
			return
		}

		http.Redirect(w, r, "/config", http.StatusSeeOther)
		return
	}

	config, err := FindConfig(r.Context())

	if err != nil {
		log.Print(err)
		internalError(w)
		return
	}

	err = render("tmpls/config.html.tmpl", w, map[string]any{
		"User":         sessionManager.Get(r.Context(), "user"),
		"Path":         r.URL.Path,
		"Title":        "Config",
		"Config":       config,
		"DomainErrors": sessionManager.Pop(r.Context(), "domain"),
	})

	if err != nil {
		log.Print(err)
		internalError(w)
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := NewField(r, "name", []Validator{
			NewMinLength(2),
			NewMaxLength(16),
		})

		password := NewField(r, "password", []Validator{
			NewMinLength(8),
		})

		if !name.Valid() || !password.Valid() {
			sessionManager.Put(r.Context(), name.Name, name.Errors)
			sessionManager.Put(r.Context(), password.Name, password.Errors)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := FindUser(r.Context(), name.Value)

		if err != nil {
			log.Print(err)
			internalError(w)
			return
		}

		err = sessionManager.RenewToken(r.Context())

		if err != nil {
			log.Print(err)
			internalError(w)
			return
		}

		sessionManager.Put(r.Context(), "user", user)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := render("tmpls/login.html.tmpl", w, map[string]any{
		"User":           sessionManager.Get(r.Context(), "user"),
		"Path":           r.URL.Path,
		"Title":          "Login",
		"NameErrors":     sessionManager.Pop(r.Context(), "name"),
		"PasswordErrors": sessionManager.Pop(r.Context(), "password"),
	})

	if err != nil {
		log.Print(err)
		internalError(w)
		return
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionManager.Clear(r.Context())
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := NewField(r, "name", []Validator{
			NewMinLength(2),
			NewMaxLength(16),
		})

		password := NewField(r, "password", []Validator{
			NewMinLength(8),
		})

		if !name.Valid() || !password.Valid() {
			sessionManager.Put(r.Context(), "name", name.Errors)
			sessionManager.Put(r.Context(), "password", password.Errors)
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		user := &User{
			Name:     name.Value,
			Password: password.Value,
		}

		err := NewUser(r.Context(), user)

		if err != nil {
			log.Print(err)
			internalError(w)
			return
		}

		sessionManager.Put(r.Context(), "user", user)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := render("tmpls/register.html.tmpl", w, map[string]any{
		"User":           sessionManager.Get(r.Context(), "user"),
		"Path":           r.URL.Path,
		"Title":          "Register",
		"NameErrors":     sessionManager.Pop(r.Context(), "name"),
		"PasswordErrors": sessionManager.Pop(r.Context(), "password"),
	})

	if err != nil {
		log.Print(err)
		internalError(w)
		return
	}
}
