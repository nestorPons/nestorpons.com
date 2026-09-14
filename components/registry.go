package components

import (
	"net/http"

	"nestorgo-landing/components/contact"
	"nestorgo-landing/components/experience"
	"nestorgo-landing/components/projects"
	"nestorgo-landing/components/skills"
)

// Registry maps component names to their HTTP handlers.
// To add a new component, just add an entry here.
var Registry = map[string]http.HandlerFunc{
	"skills":     skills.Handler,
	"projects":   projects.Handler,
	"experience": experience.Handler,
	"contact":    contact.Handler,
}
