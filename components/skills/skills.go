package skills

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Skill struct {
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Level int    `json:"level"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	jsonData, err := os.ReadFile("components/skills/skills.json")
	if err != nil {
		http.Error(w, "skills json not found", http.StatusInternalServerError)
		return
	}

	var skills []Skill
	if err := json.Unmarshal(jsonData, &skills); err != nil {
		http.Error(w, "invalid skills json", http.StatusInternalServerError)
		return
	}

	tmplData, err := os.ReadFile("components/skills/skills.html")
	if err != nil {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("skills").Parse(string(tmplData))
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, skills); err != nil {
		log.Printf("skills template error: %v", err)
	}
}
