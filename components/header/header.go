package header

import (
	_ "embed"
	"net/http"
)

//go:embed header.html
var headerHTML string

//go:embed translator.js
var translatorJS string

// Handler sirve el header (logo + menú + selector de idioma). Solo GET.
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(headerHTML))
}

// TranslatorHandler sirve el JS de traducción del header. Solo GET.
func TranslatorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Write([]byte(translatorJS))
}
