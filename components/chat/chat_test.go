package chat

import "testing"

func TestNormalizeLang(t *testing.T) {
	cases := map[string]string{
		"es": "es", "en": "en", "EN": "en", "": "es",
		"fr": "es", "español": "es", " en ": "en",
	}
	for in, want := range cases {
		if got := normalizeLang(in); got != want {
			t.Errorf("normalizeLang(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestBuildPromptVars(t *testing.T) {
	t.Setenv("OPENAI_LANG_VAR", "idioma")

	vars := buildPromptVars(nil, "en")
	if vars["lang"] != "inglés" || vars["datetime"] == "" {
		t.Errorf("vars = %v, want lang and datetime", vars)
	}
	if vars["idioma"] != "inglés" {
		t.Errorf("vars = %v, want idioma=inglés", vars)
	}

	// La base fija se conserva y el idioma de la petición tiene prioridad.
	vars = buildPromptVars(map[string]any{"idioma": "francés", "tono": "formal"}, "es")
	if vars["idioma"] != "español" || vars["tono"] != "formal" {
		t.Errorf("vars = %v, want idioma=español tono=formal", vars)
	}

	// Aunque no haya variable configurable, se conservan las variables requeridas
	// por el prompt publicado.
	t.Setenv("OPENAI_LANG_VAR", "")
	if vars := buildPromptVars(nil, "es"); vars["lang"] != "español" || vars["datetime"] == "" {
		t.Errorf("vars = %v, want lang and datetime", vars)
	}
}

func TestPromptVersionOmitted(t *testing.T) {
	t.Setenv("OPENAI_PROMPT_ID", "pmpt_test")
	t.Setenv("OPENAI_PROMPT_VERSION", "")
	t.Setenv("OPENAI_LANG_VAR", "")

	req := newResponsesRequest("hola", "", "es", false)
	if _, ok := req.Prompt["version"]; ok {
		t.Errorf("prompt = %v, version debe omitirse para usar la default del dashboard", req.Prompt)
	}

	t.Setenv("OPENAI_PROMPT_VERSION", "2")
	req = newResponsesRequest("hola", "", "es", false)
	if req.Prompt["version"] != "2" {
		t.Errorf("prompt = %v, want version=2", req.Prompt)
	}
}
