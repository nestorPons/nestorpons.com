package contact

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"strings"
	"text/template"
)

//go:embed form.html
var formHTML string

//go:embed success.html
var successHTML string

var successTmpl = template.Must(template.New("success").Parse(successHTML))

type form struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Message        string `json:"message"`
	Privacy        string `json:"privacy"`
	RecaptchaToken string `json:"cf-turnstile-response"`
}

type response struct {
	Ok    bool
	Error string
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 1. SI ES GET: Servimos el HTML embebido del formulario
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(formHTML))
		return
	}

	// 2. SI ES POST: Procesamos el envío (lo que ya tienes hecho)
	if r.Method != http.MethodPost {
		htmlError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	if r.Method != http.MethodPost {
		htmlError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	f, err := parseForm(r)
	if err != nil {
		log.Printf("contact form parse error: %v", err)
		htmlError(w, http.StatusBadRequest, "Solicitud inválida")
		return
	}

	f.Name = strings.TrimSpace(f.Name)
	f.Email = strings.TrimSpace(f.Email)
	f.Message = strings.TrimSpace(f.Message)

	if f.Name == "" || f.Email == "" || f.Message == "" {
		htmlError(w, http.StatusBadRequest, "Todos los campos son obligatorios")
		return
	}

	if f.Privacy != "on" && f.Privacy != "true" && f.Privacy != "1" && f.Privacy != "yes" {
		htmlError(w, http.StatusBadRequest, "Debes aceptar la política de privacidad")
		return
	}

	ok, err := verifyRecaptcha(f.RecaptchaToken)
	if err != nil || !ok {
		log.Printf("reCaptcha failed: %v", err)
		msg := "Verificación de seguridad fallida"
		if err != nil {
			msg = err.Error()
		}
		htmlError(w, http.StatusForbidden, msg)
		return
	}

	if err := sendEmail(f); err != nil {
		log.Printf("Error sending email: %v", err)
		email := os.Getenv("SMTP_USER")
		htmlError(w, http.StatusInternalServerError, "Error al enviar. Escríbeme a "+email)
		return
	}

	render(w, http.StatusOK, response{Ok: true})
}

func parseForm(r *http.Request) (form, error) {
	var f form
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return f, err
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	contentType := r.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/json") {
		if err := json.Unmarshal(body, &f); err == nil {
			return f, nil
		}
		log.Printf("contact: JSON decode failed, falling back to form data")
		r.Body = io.NopCloser(bytes.NewReader(body))
	}

	if err := r.ParseForm(); err != nil {
		return f, err
	}

	f.Name = r.FormValue("name")
	f.Email = r.FormValue("email")
	f.Message = r.FormValue("message")
	f.Privacy = r.FormValue("privacy")
	f.RecaptchaToken = r.FormValue("cf-turnstile-response")
	return f, nil
}

func render(w http.ResponseWriter, status int, data response) {
	tmplData, err := os.ReadFile("components/contact/success.html")
	if err != nil {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("success").Parse(string(tmplData))
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("contact template error: %v", err)
	}
}

func htmlError(w http.ResponseWriter, status int, msg string) {
	render(w, status, response{Ok: false, Error: msg})
}

func verifyRecaptcha(token string) (bool, error) {
	if os.Getenv("RECAPTCHA_SKIP") == "true" {
		return true, nil
	}

	secret := os.Getenv("TURNSTILE_SECRET")
	if secret == "" {
		return false, errors.New("TURNSTILE_SECRET not set")
	}

	if strings.TrimSpace(token) == "" {
		return false, errors.New("turnstile token is empty")
	}

	resp, err := http.PostForm("https://challenges.cloudflare.com/turnstile/v0/siteverify", url.Values{
		"secret":   {secret},
		"response": {token},
	})
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, err
	}
	if !result.Success {
		return false, errors.New("turnstile rejected the token")
	}

	return true, nil
}

func sendEmail(f form) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	to := os.Getenv("SMTP_TO")

	if port == "" {
		port = "587"
	}

	auth := smtp.PlainAuth("", user, pass, host)

	body := "From: " + user + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: Contacto web - " + f.Name + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n\r\n" +
		"Nombre: " + f.Name + "\r\n" +
		"Email: " + f.Email + "\r\n\r\n" +
		f.Message

	return smtp.SendMail(host+":"+port, auth, user, []string{to}, []byte(body))
}
