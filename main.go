package main

import (
	"log"
	"net/http"
	"os"

	"nestorgo-landing/components"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	dir := os.Getenv("SERVE_DIR")
	if dir == "" {
		dir = "/var/www/html"
	}

	for name, handler := range components.Registry {
		http.HandleFunc("/api/components/"+name, handler)
		log.Printf("📦 Component registered: /api/components/%s", name)
	}

	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	log.Printf("🚀 Server running on port %s, serving %s", port, dir)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
