package skills

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
)

type SkillItem struct {
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Use   string `json:"use"`
	Badge string `json:"badge,omitempty"`
}

type SkillGroup struct {
	Group string      `json:"group"`
	Items []SkillItem `json:"items"`
}

type WheelGroup struct {
	Group    string
	Items    []SkillItem
	Arc      string
	LabelX   float64
	LabelY   float64
	LabelRot float64
}

const (
	wheelCX   = 260.0
	wheelCY   = 260.0
	wheelR    = 220.0
	arcSpan   = 62.0
	wheelView = 520
)

func pointOnCircle(deg, radius float64) (float64, float64) {
	rad := deg * math.Pi / 180
	return wheelCX + radius*math.Sin(rad), wheelCY - radius*math.Cos(rad)
}

func buildWheel(groups []SkillGroup) []WheelGroup {
	n := len(groups)
	out := make([]WheelGroup, 0, n)
	if n == 0 {
		return out
	}
	step := 360.0 / float64(n)
	for i, g := range groups {
		center := float64(i) * step
		x1, y1 := pointOnCircle(center-arcSpan/2, wheelR)
		x2, y2 := pointOnCircle(center+arcSpan/2, wheelR)
		lx, ly := pointOnCircle(center, wheelR)
		out = append(out, WheelGroup{
			Group:    g.Group,
			Items:    g.Items,
			Arc:      fmt.Sprintf("M %.2f %.2f A %.0f %.0f 0 0 1 %.2f %.2f", x1, y1, wheelR, wheelR, x2, y2),
			LabelX:   math.Round(lx*100) / 100,
			LabelY:   math.Round(ly*100) / 100,
			LabelRot: center,
		})
	}
	return out
}

func Handler(w http.ResponseWriter, r *http.Request) {
	jsonData, err := os.ReadFile("components/skills/skills.json")
	if err != nil {
		http.Error(w, "skills json not found", http.StatusInternalServerError)
		return
	}

	var groups []SkillGroup
	if err := json.Unmarshal(jsonData, &groups); err != nil {
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
	if err := tmpl.Execute(w, buildWheel(groups)); err != nil {
		log.Printf("skills template error: %v", err)
	}
}
