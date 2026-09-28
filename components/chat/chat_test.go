package chat

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

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
	if vars["idioma"] != "inglés" {
		t.Errorf("vars = %v, want idioma=inglés", vars)
	}

	// La base fija se conserva y el idioma de la petición tiene prioridad.
	vars = buildPromptVars(map[string]any{"idioma": "francés", "tono": "formal"}, "es")
	if vars["idioma"] != "español" || vars["tono"] != "formal" {
		t.Errorf("vars = %v, want idioma=español tono=formal", vars)
	}

	// Sin variable de idioma ni base -> nil (no se envía "variables").
	t.Setenv("OPENAI_LANG_VAR", "")
	if vars := buildPromptVars(nil, "es"); vars != nil {
		t.Errorf("vars = %v, want nil", vars)
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

func TestContactDraftToolIsStrictAndDoesNotSubmit(t *testing.T) {
	req := newResponsesRequest("quiero contactar", "", "es", false)
	if req.ToolChoice != "auto" || len(req.Tools) != 1 {
		t.Fatalf("request tools = %v, choice = %q", req.Tools, req.ToolChoice)
	}
	tool := req.Tools[0]
	if tool.Name != "prepare_contact_draft" || !tool.Strict {
		t.Fatalf("tool = %+v, want strict prepare_contact_draft", tool)
	}
	if len(tool.Parameters["required"].([]string)) != 3 {
		t.Fatalf("required = %v, want name, email, message", tool.Parameters["required"])
	}
	if !strings.Contains(tool.Description, "No envía el formulario") || !strings.Contains(tool.Description, "no marca la casilla") || !strings.Contains(tool.Description, "no completa Turnstile") {
		t.Fatalf("tool description must explicitly forbid submission and consent actions: %q", tool.Description)
	}
	if !strings.Contains(req.Instructions, "pregunta solo por los datos que falten") || !strings.Contains(req.Instructions, "nunca afirmes que el mensaje se envió") {
		t.Fatalf("instructions do not define the safe contact flow: %q", req.Instructions)
	}
}

func TestParseContactDraft(t *testing.T) {
	draft, err := parseContactDraft(`{"name":" Ada ","email":"ada@example.com","message":" Hola "}`)
	if err != nil {
		t.Fatal(err)
	}
	if draft.Name != "Ada" || draft.Email != "ada@example.com" || draft.Message != "Hola" {
		t.Fatalf("draft = %+v, want trimmed values", draft)
	}

	invalid := []string{
		`{"name":"","email":"ada@example.com","message":"Hola"}`,
		`{"name":"Ada","email":"ada@example.com\r\nBcc:x@example.com","message":"Hola"}`,
		`{"name":"Ada","email":"not-an-email","message":"Hola"}`,
		`{"name":"Ada","email":"ada@example.com","message":" "}`,
		`{"name":"Ada","email":"ada@example.com","message":"Hola","extra":"no"}`,
	}
	for _, raw := range invalid {
		t.Run(raw, func(t *testing.T) {
			if _, err := parseContactDraft(raw); err == nil {
				t.Fatalf("parseContactDraft(%s) unexpectedly succeeded", raw)
			}
		})
	}
}

func TestInvalidDraftReturnsToolErrorWithoutContactDraft(t *testing.T) {
	output, draft := contactToolOutput(`{"name":"Ada","email":"bad-email","message":"Hola"}`)
	if draft != nil {
		t.Fatalf("draft = %+v, want nil for invalid email", draft)
	}
	if !strings.Contains(output, "correo") && !strings.Contains(output, "válidos") {
		t.Fatalf("tool output should request valid missing data: %q", output)
	}
}

func TestResponsesRequestSerializesContactTool(t *testing.T) {
	encoded, err := json.Marshal(newResponsesRequest("hola", "", "es", false))
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	tools, ok := decoded["tools"].([]any)
	if !ok || len(tools) != 1 {
		t.Fatalf("serialized tools = %v", decoded["tools"])
	}
	tool := tools[0].(map[string]any)
	if tool["name"] != "prepare_contact_draft" || tool["strict"] != true {
		t.Fatalf("serialized tool = %v", tool)
	}
}

func TestReadStreamResponseCollectsFunctionCall(t *testing.T) {
	stream := strings.Join([]string{
		`data: {"type":"response.output_text.delta","delta":"Preparando"}`,
		`data: {"type":"response.completed","response":{"id":"resp_tool","output":[{"type":"function_call","name":"prepare_contact_draft","call_id":"call_contact","arguments":"{\"name\":\"Ada\",\"email\":\"ada@example.com\",\"message\":\"Hola\"}"}]}}`,
		"data: [DONE]",
		"",
	}, "\n")
	response := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(stream))}

	var streamed strings.Builder
	out, text, err := readStreamResponse(response, func(delta string) { streamed.WriteString(delta) })
	if err != nil {
		t.Fatal(err)
	}
	if out.ID != "resp_tool" || text != "Preparando" {
		t.Fatalf("response = %+v, text = %q", out, text)
	}
	if streamed.String() != text {
		t.Fatalf("streamed text = %q, want %q", streamed.String(), text)
	}
	call := findContactToolCall(out.Output)
	if call == nil || call.CallID != "call_contact" {
		t.Fatalf("function call = %+v", call)
	}
	draft, err := parseContactDraft(call.Arguments)
	if err != nil || draft.Email != "ada@example.com" {
		t.Fatalf("draft = %+v, err = %v", draft, err)
	}
}

func TestToolOutputRequestContinuesWithoutPromptOrTools(t *testing.T) {
	request := newToolOutputRequest("resp_tool", "call_contact", `{"ok":true}`, "gpt-4o-mini", false)
	if request.PrevID != "resp_tool" || len(request.Input) != 1 {
		t.Fatalf("request = %+v", request)
	}
	if request.Model != "gpt-4o-mini" {
		t.Fatalf("request model = %q, want gpt-4o-mini", request.Model)
	}
	if request.Prompt != nil || len(request.Tools) != 0 {
		t.Fatalf("continuation should use previous response context without redeclaring prompt/tools: %+v", request)
	}
	if !strings.Contains(request.Instructions, "nunca afirmes que el mensaje se envió") {
		t.Fatalf("continuation instructions = %q", request.Instructions)
	}
	encoded, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal(encoded, &body); err != nil {
		t.Fatal(err)
	}
	if body["model"] != "gpt-4o-mini" {
		t.Fatalf("serialized model = %v, want gpt-4o-mini", body["model"])
	}
	items := body["input"].([]any)
	item := items[0].(map[string]any)
	if item["type"] != "function_call_output" || item["call_id"] != "call_contact" || item["output"] != `{"ok":true}` {
		t.Fatalf("serialized function output = %v", item)
	}
}
