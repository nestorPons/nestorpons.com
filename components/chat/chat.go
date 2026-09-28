package chat

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"sync"
	"time"
)

//go:embed chat.html
var chatHTML string

// Config desde entorno. La API key nunca sale al frontend:
// el navegador habla con este handler y es el backend quien llama a OpenAI.
func apiKey() string { return os.Getenv("OPENAI_API_KEY") }
func promptID() string {
	return firstNonEmpty(os.Getenv("OPENAI_PROMPT_ID"), "pmpt_6aa92a8690c881959522b0836077156d0e8ed2595823acc7")
}
func promptVersion() string { return strings.TrimSpace(os.Getenv("OPENAI_PROMPT_VERSION")) }

// Variables del prompt almacenado (p. ej. {"pregunta":"..."} si la plantilla
// usa {{variables}}). El playground las rellena desde su UI; por API hay que
// enviarlas explícitamente o la API rechaza la petición si falta alguna.
// Formato: OPENAI_PROMPT_VARIABLES='{"nombre":"valor"}'
func promptVariables() map[string]any {
	raw := strings.TrimSpace(os.Getenv("OPENAI_PROMPT_VARIABLES"))
	if raw == "" {
		return nil
	}
	var vars map[string]any
	if err := json.Unmarshal([]byte(raw), &vars); err != nil {
		log.Printf("chat: OPENAI_PROMPT_VARIABLES inválido, se ignora: %v", err)
		return nil
	}
	return vars
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

const (
	maxMessageRunes = 2000
	maxContactName  = 120
	maxContactEmail = 254
	maxContactBody  = 4000
	maxBodyBytes    = 16 << 10
	openAITimeout   = 90 * time.Second
	rateWindow      = time.Minute
	rateMax         = 20 // peticiones/min por IP
)

// ── Rate limit en memoria ─────────────────────────────────────────────

var rateMu sync.Mutex
var rateHits = map[string][]time.Time{}

func allow(ip string) bool {
	now := time.Now()
	rateMu.Lock()
	defer rateMu.Unlock()
	hits := rateHits[ip]
	kept := hits[:0]
	for _, t := range hits {
		if now.Sub(t) < rateWindow {
			kept = append(kept, t)
		}
	}
	if len(kept) >= rateMax {
		rateHits[ip] = kept
		return false
	}
	rateHits[ip] = append(kept, now)
	return true
}

func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		if ip, _, err := net.SplitHostPort(strings.TrimSpace(strings.Split(fwd, ",")[0]) + ":0"); err == nil {
			return ip
		}
		return strings.TrimSpace(strings.Split(fwd, ",")[0])
	}
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return ip
	}
	return r.RemoteAddr
}

// ── Handlers ──────────────────────────────────────────────────────────

// StatusHandler diagnostica la configuración efectiva del chat (GET).
// No expone secretos: de la API key solo dice si está definida o no.
// Incluye el JSON exacto que se enviaría a OpenAI para un mensaje de ejemplo,
// para compararlo con un curl que funcione.
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	sample := strings.TrimSpace(r.URL.Query().Get("message"))
	if sample == "" {
		sample = "ping"
	}
	if len([]rune(sample)) > 200 {
		sample = string([]rune(sample)[:200])
	}
	lang := normalizeLang(r.URL.Query().Get("lang"))
	example := newResponsesRequest(sample, "", lang, false)

	version := promptVersion()
	versionNote := "fija"
	if version == "" {
		versionNote = "omitida -> la API usa la default del dashboard"
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"api_key_set":         apiKey() != "",
		"prompt_id":           promptID(),
		"prompt_version":      version,
		"prompt_version_note": versionNote,
		"lang_var":            langVarName(),
		"lang":                lang,
		"example_request":     example,
	})
}

// Handler sirve la UI del chat (GET) o responde una pregunta sin streaming (POST JSON).
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(chatHTML))
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !allow(clientIP(r)) {
		writeJSON(w, http.StatusTooManyRequests, map[string]string{"error": "Demasiadas peticiones, espera un momento e inténtalo de nuevo."})
		return
	}

	msg, prevID, lang, err := parseInput(w, r)
	if err != nil {
		return // parseInput ya respondió
	}

	result, err := complete(r.Context(), msg, prevID, lang)
	if err != nil {
		log.Printf("chat complete error: %v", err)
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "El asistente no responde ahora mismo. Inténtalo de nuevo."})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// StreamHandler hace proxy SSE de OpenAI al navegador:
// POST /api/components/chat/stream -> text/event-stream con
// data: {"delta":"..."} ... data: {"done":true,"response_id":"resp_..."}
func StreamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !allow(clientIP(r)) {
		writeJSON(w, http.StatusTooManyRequests, map[string]string{"error": "Demasiadas peticiones, espera un momento e inténtalo de nuevo."})
		return
	}

	msg, prevID, lang, err := parseInput(w, r)
	if err != nil {
		return
	}

	upstream, err := openStream(r.Context(), newResponsesRequest(msg, prevID, lang, true))
	if err != nil {
		log.Printf("chat stream error: %v", err)
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "El asistente no responde ahora mismo. Inténtalo de nuevo."})
		return
	}
	defer upstream.Body.Close()

	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	first, _, err := readStreamResponse(upstream, func(delta string) {
		emit(w, flusher, map[string]string{"delta": delta})
	})
	if err != nil {
		log.Printf("chat stream read error: %v", err)
		emit(w, flusher, map[string]string{"error": "El asistente no responde ahora mismo. Inténtalo de nuevo."})
		return
	}
	call := findContactToolCall(first.Output)
	if call == nil {
		emit(w, flusher, map[string]any{"done": true, "response_id": first.ID})
		return
	}

	toolOutput, draft := contactToolOutput(call.Arguments)
	followup, err := openStream(r.Context(), newToolOutputRequest(first.ID, call.CallID, toolOutput, first.Model, true))
	if err != nil {
		log.Printf("chat tool follow-up error: %v", err)
		emit(w, flusher, map[string]string{"error": "No pude completar la preparación del formulario. Inténtalo de nuevo."})
		return
	}
	defer followup.Body.Close()
	last, _, err := readStreamResponse(followup, func(delta string) {
		emit(w, flusher, map[string]string{"delta": delta})
	})
	if err != nil {
		log.Printf("chat tool follow-up stream error: %v", err)
		emit(w, flusher, map[string]string{"error": "No pude completar la preparación del formulario. Inténtalo de nuevo."})
		return
	}
	if draft != nil {
		emit(w, flusher, map[string]any{"contact_draft": draft})
	}
	emit(w, flusher, map[string]any{"done": true, "response_id": last.ID})
}

func emit(w http.ResponseWriter, f http.Flusher, v any) {
	b, _ := json.Marshal(v)
	w.Write([]byte("data: " + string(b) + "\n\n"))
	f.Flush()
}

type chatRequest struct {
	Message          string `json:"message"`
	PreviousResponse string `json:"previous_response_id"`
	Lang             string `json:"lang"`
}

func parseInput(w http.ResponseWriter, r *http.Request) (msg, prevID, lang string, err error) {
	if apiKey() == "" {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Chat no configurado (falta OPENAI_API_KEY)."})
		return "", "", "", io.EOF
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	var in chatRequest
	ct := r.Header.Get("Content-Type")
	if strings.Contains(ct, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Petición inválida."})
			return "", "", "", err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Petición inválida."})
			return "", "", "", err
		}
		in.Message = r.FormValue("message")
		in.PreviousResponse = r.FormValue("previous_response_id")
		in.Lang = r.FormValue("lang")
	}

	msg = strings.TrimSpace(in.Message)
	if msg == "" || len([]rune(msg)) > maxMessageRunes {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "El mensaje debe tener entre 1 y 2000 caracteres."})
		return "", "", "", io.EOF
	}
	prevID = strings.TrimSpace(in.PreviousResponse)
	if prevID != "" && !strings.HasPrefix(prevID, "resp_") {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Identificador de conversación inválido."})
		return "", "", "", io.EOF
	}
	return msg, prevID, normalizeLang(in.Lang), nil
}

// normalizeLang deja solo 'es' o 'en' (idiomas del sitio). Todo lo demás -> 'es'.
func normalizeLang(raw string) string {
	if strings.ToLower(strings.TrimSpace(raw)) == "en" {
		return "en"
	}
	return "es"
}

// Nombre de la variable de la plantilla que recibe el idioma.
// Vacío = no inyectar (por si el prompt no define {{idioma}}).
func langVarName() string { return strings.TrimSpace(os.Getenv("OPENAI_LANG_VAR")) }

func langDisplayName(code string) string {
	if code == "en" {
		return "inglés"
	}
	return "español"
}

func currentDateTime() string {
	return time.Now().Format("02/01/2006 15:04")
}

// buildPromptVars combina las variables fijas (OPENAI_PROMPT_VARIABLES) con el
// idioma de la petición. El idioma de la petición tiene prioridad.
func buildPromptVars(base map[string]any, lang string) map[string]any {
	vars := map[string]any{}
	for k, v := range base {
		vars[k] = v
	}
	// Estas variables son requeridas por el prompt publicado actualmente.
	// También se envía "lang" aunque OPENAI_LANG_VAR apunte a un nombre
	// histórico distinto, para que el prompt no falle por configuración antigua.
	vars["lang"] = langDisplayName(normalizeLang(lang))
	vars["datetime"] = currentDateTime()
	if name := langVarName(); name != "" {
		vars[name] = langDisplayName(normalizeLang(lang))
	}
	return vars
}

// ── OpenAI Responses API (equivalente Go del snippet JS del usuario) ──

type responsesRequest struct {
	Prompt       map[string]any  `json:"prompt,omitempty"`
	Instructions string          `json:"instructions,omitempty"`
	Input        []any           `json:"input"`
	PrevID       string          `json:"previous_response_id,omitempty"`
	Model        string          `json:"model,omitempty"`
	Reason       reasoningConfig `json:"reasoning"`
	Store        bool            `json:"store"`
	Stream       bool            `json:"stream"`
	Include      []string        `json:"include,omitempty"`
	Tools        []functionTool  `json:"tools,omitempty"`
	ToolChoice   string          `json:"tool_choice,omitempty"`
}

type reasoningConfig struct {
	Mode    string `json:"mode"`
	Summary string `json:"summary"`
}

type inputItem struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type functionTool struct {
	Type        string         `json:"type"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
	Strict      bool           `json:"strict"`
}

type contactDraft struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

func contactDraftTool() functionTool {
	return functionTool{
		Type:        "function",
		Name:        "prepare_contact_draft",
		Description: "Prepara un borrador en el formulario de contacto cuando el usuario quiera contactar y haya proporcionado nombre, correo y mensaje. No envía el formulario, no marca la casilla de privacidad y no completa Turnstile. Si falta algún dato, pregunta primero por él.",
		Strict:      true,
		Parameters: map[string]any{
			"type":                 "object",
			"additionalProperties": false,
			"properties": map[string]any{
				"name":    map[string]any{"type": "string", "description": "Nombre de la persona", "maxLength": maxContactName},
				"email":   map[string]any{"type": "string", "description": "Correo electrónico de contacto", "maxLength": maxContactEmail},
				"message": map[string]any{"type": "string", "description": "Mensaje para el formulario de contacto", "maxLength": maxContactBody},
			},
			"required": []string{"name", "email", "message"},
		},
	}
}

func parseContactDraft(raw string) (contactDraft, error) {
	var draft contactDraft
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&draft); err != nil {
		return draft, err
	}
	draft.Name = strings.TrimSpace(draft.Name)
	draft.Email = strings.TrimSpace(draft.Email)
	draft.Message = strings.TrimSpace(draft.Message)
	if draft.Name == "" || len([]rune(draft.Name)) > maxContactName {
		return contactDraft{}, errors.New("nombre de contacto inválido")
	}
	if len(draft.Email) > maxContactEmail {
		return contactDraft{}, errors.New("correo de contacto demasiado largo")
	}
	parsedEmail, err := mail.ParseAddress(draft.Email)
	if err != nil || parsedEmail.Address != draft.Email {
		return contactDraft{}, errors.New("correo de contacto inválido")
	}
	if draft.Message == "" || len([]rune(draft.Message)) > maxContactBody {
		return contactDraft{}, errors.New("mensaje de contacto inválido")
	}
	return draft, nil
}

func newResponsesRequest(msg, prevID, lang string, stream bool) responsesRequest {
	// Sin versión -> la API usa la versión marcada como "default" en el dashboard.
	// No existe el valor literal "default": hay que omitir el campo.
	prompt := map[string]any{"id": promptID()}
	if v := promptVersion(); v != "" {
		prompt["version"] = v
	}
	if vars := buildPromptVars(promptVariables(), lang); vars != nil {
		prompt["variables"] = vars
	}
	return responsesRequest{
		Prompt:       prompt,
		Instructions: contactInstructions(),
		Input:        []any{inputItem{Role: "user", Content: msg}},
		PrevID:       prevID,
		Reason:       reasoningConfig{Mode: "standard", Summary: "auto"},
		Store:        true,
		Stream:       stream,
		Include: []string{
			"reasoning.encrypted_content",
			"web_search_call.action.sources",
		},
		Tools:      []functionTool{contactDraftTool()},
		ToolChoice: "auto",
	}
}

func contactInstructions() string {
	return "Responde en el idioma del usuario. Si quiere contactar, recopila nombre, correo y mensaje; pregunta solo por los datos que falten. Cuando los tres estén presentes y sean válidos, llama a prepare_contact_draft. La herramienta solo prepara el formulario: nunca afirmes que el mensaje se envió. Indica que debe revisar los datos, aceptar la política con la casilla, completar Turnstile manualmente y pulsar Enviar. No marques ni envíes el formulario."
}

func openAIClient() *http.Client {
	return &http.Client{Timeout: openAITimeout}
}

func postResponses(ctx context.Context, body responsesRequest) (*http.Response, error) {
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/responses", bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey())
	req.Header.Set("Content-Type", "application/json")
	return openAIClient().Do(req)
}

// complete: llamada sin streaming, devuelve texto + id de respuesta.
func complete(ctx context.Context, msg, prevID, lang string) (chatResult, error) {
	resp, err := postResponses(ctx, newResponsesRequest(msg, prevID, lang, false))
	if err != nil {
		return chatResult{}, err
	}
	first, err := readResponse(resp)
	if err != nil {
		return chatResult{}, err
	}
	call := findContactToolCall(first.Output)
	if call == nil {
		return chatResult{Reply: responseText(first.Output), ResponseID: first.ID}, nil
	}

	toolOutput, draft := contactToolOutput(call.Arguments)
	resp, err = postResponses(ctx, newToolOutputRequest(first.ID, call.CallID, toolOutput, first.Model, false))
	if err != nil {
		return chatResult{}, err
	}
	last, err := readResponse(resp)
	if err != nil {
		return chatResult{}, err
	}
	return chatResult{Reply: responseText(last.Output), ResponseID: last.ID, ContactDraft: draft}, nil
}

type responseEnvelope struct {
	ID     string               `json:"id"`
	Model  string               `json:"model"`
	Output []responseOutputItem `json:"output"`
}

type responseOutputItem struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	CallID    string `json:"call_id"`
	Arguments string `json:"arguments"`
	Content   []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}

type functionCallOutput struct {
	Type   string `json:"type"`
	CallID string `json:"call_id"`
	Output string `json:"output"`
}

func responseText(items []responseOutputItem) string {
	var text strings.Builder
	for _, item := range items {
		for _, content := range item.Content {
			if content.Type == "output_text" || content.Type == "refusal" {
				text.WriteString(content.Text)
			}
		}
	}
	return text.String()
}

func findContactToolCall(items []responseOutputItem) *responseOutputItem {
	for i := range items {
		if items[i].Type == "function_call" && items[i].Name == "prepare_contact_draft" {
			return &items[i]
		}
	}
	return nil
}

func contactToolOutput(raw string) (string, *contactDraft) {
	draft, err := parseContactDraft(raw)
	if err != nil {
		return `{"ok":false,"message":"Faltan datos válidos: solicita nombre, correo y mensaje antes de preparar el formulario."}`, nil
	}
	encoded, _ := json.Marshal(map[string]any{"ok": true, "message": "El borrador está preparado. Indica al usuario que revise el formulario, marque la casilla de privacidad, complete Turnstile y pulse el botón Enviar. No afirmes que el mensaje ya se envió."})
	return string(encoded), &draft
}

func newToolOutputRequest(previousID, callID, output, model string, stream bool) responsesRequest {
	return responsesRequest{
		Instructions: contactInstructions(),
		Input:        []any{functionCallOutput{Type: "function_call_output", CallID: callID, Output: output}},
		PrevID:       previousID,
		Model:        model,
		Reason:       reasoningConfig{Mode: "standard", Summary: "auto"},
		Store:        true,
		Stream:       stream,
	}
}

func readResponse(resp *http.Response) (responseEnvelope, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return responseEnvelope{}, err
	}
	if resp.StatusCode >= 300 {
		return responseEnvelope{}, &apiError{status: resp.StatusCode, body: strings.TrimSpace(string(body))}
	}
	var out responseEnvelope
	if err := json.Unmarshal(body, &out); err != nil {
		return responseEnvelope{}, err
	}
	if out.ID == "" {
		return responseEnvelope{}, errors.New("OpenAI response has no id")
	}
	return out, nil
}

func readStreamResponse(resp *http.Response, onDelta func(string)) (responseEnvelope, string, error) {
	defer resp.Body.Close()
	sc := bufio.NewScanner(resp.Body)
	sc.Buffer(make([]byte, 64*1024), 1024*1024)
	var text strings.Builder
	var out responseEnvelope
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if payload == "[DONE]" {
			break
		}
		var evt struct {
			Type     string           `json:"type"`
			Delta    string           `json:"delta"`
			Response responseEnvelope `json:"response"`
			Error    struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal([]byte(payload), &evt); err != nil {
			continue
		}
		switch evt.Type {
		case "response.output_text.delta", "response.refusal.delta":
			text.WriteString(evt.Delta)
			if onDelta != nil && evt.Delta != "" {
				onDelta(evt.Delta)
			}
		case "response.completed":
			out = evt.Response
		case "response.failed":
			if evt.Error.Message == "" {
				evt.Error.Message = "OpenAI response failed"
			}
			return responseEnvelope{}, "", errors.New(evt.Error.Message)
		}
	}
	if err := sc.Err(); err != nil {
		return responseEnvelope{}, "", err
	}
	if out.ID == "" {
		return responseEnvelope{}, "", errors.New("OpenAI stream ended without a completed response")
	}
	return out, text.String(), nil
}

// openStream: llamada con streaming, devuelve el body SSE para reenviar.
func openStream(ctx context.Context, request responsesRequest) (*http.Response, error) {
	resp, err := postResponses(ctx, request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
		resp.Body.Close()
		return nil, &apiError{status: resp.StatusCode, body: strings.TrimSpace(string(body))}
	}
	return resp, nil
}

type apiError struct {
	status int
	body   string
}

func (e *apiError) Error() string { return http.StatusText(e.status) + ": " + e.body }

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

type chatResult struct {
	Reply        string        `json:"reply"`
	ResponseID   string        `json:"response_id"`
	ContactDraft *contactDraft `json:"contact_draft,omitempty"`
}
