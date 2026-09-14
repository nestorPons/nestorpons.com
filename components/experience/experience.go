package experience

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Experience struct {
	Type        string   `json:"type"`
	Role        string   `json:"role"`
	Company     string   `json:"company"`
	Period      string   `json:"period"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	jsonData, err := os.ReadFile("components/experience/experience.json")
	if err != nil {
		http.Error(w, "experience json not found", http.StatusInternalServerError)
		return
	}

	var expList []Experience
	if err := json.Unmarshal(jsonData, &expList); err != nil {
		http.Error(w, "invalid experience json", http.StatusInternalServerError)
		return
	}

	tmplData, err := os.ReadFile("components/experience/experience.html")
	if err != nil {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("experience").Parse(string(tmplData))
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, expList); err != nil {
		log.Printf("experience template error: %v", err)
	}
}
