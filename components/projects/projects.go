package projects

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Project struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Url         string   `json:"url"`
	Logo        Logo     `json:"logo"`
	Tecs        []string `json:"tecs"`
}

type Logo struct {
	Src   string `json:"src"`
	Alt   string `json:"alt"`
	Class string `json:"class"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	jsonData, err := os.ReadFile("components/projects/projects.json")
	if err != nil {
		http.Error(w, "projects json not found", http.StatusInternalServerError)
		return
	}

	var projects []Project
	if err := json.Unmarshal(jsonData, &projects); err != nil {
		http.Error(w, "invalid projects json", http.StatusInternalServerError)
		return
	}

	tmplData, err := os.ReadFile("components/projects/projects.html")
	if err != nil {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("projects").Parse(string(tmplData))
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.Execute(w, projects); err != nil {
		log.Printf("projects template error: %v", err)
	}
}
