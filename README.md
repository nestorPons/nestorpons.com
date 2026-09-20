# nestorgo.com — Portfolio Néstor Pons

Landing + portfolio personal. Backend mínimo en Go que sirve estático y expone componentes HTML vía htmx.

- URL: https://nestorpons.com
- Stack: Go 1.22 · htmx · Tailwind CSS · Docker + Traefik

## Estructura

- `main.go` — estático (`SERVE_DIR`) + registro `/api/components/*`
- `components/` — `registry.go` + 1 carpeta por componente (`*.go` + `*.json` + `*.html`)
  - `skills` / `projects` / `experience` — leen `*.json`, renderizan `*.html`
  - `contact` — `GET` sirve `form.html`, `POST` exige checkbox privacidad + Turnstile y envía por SMTP
  - `header` — `GET` sirve el header (logo + menú + selector idioma) y `GET /api/components/header/translator.js` sirve su JS de traducción; `main.js` usa delegación al inyectarse por HTMX
  - `chat` — `GET` sirve la UI; `POST /api/components/chat` responde JSON `{reply, response_id}`; `POST /api/components/chat/stream` reenvía SSE `data: {"delta"}` … `{"done", "response_id"}`. Proxy a OpenAI Responses API con el prompt almacenado (`OPENAI_PROMPT_ID`); la key nunca sale al frontend
  - Diagnóstico: `GET /api/components/chat/status?message=hola&lang=es` muestra la config efectiva y el JSON exacto enviado a OpenAI (sin secretos)
- `public_html/` — `index.html`, `privacidad.html` + `assets/` (`main.js`, `styles.css`, `images/`, `cv.pdf`)
- `Dockerfile` — multi-stage `golang:1.22-alpine` → `alpine:3.19`
- `docker-compose.yml` — servicio + labels Traefik

## Requisitos

- Go 1.22+ (solo dev local)
- Docker + Compose (deploy)
- Red externa `traefik-net`

## Uso local

```bash
go run .                          # PORT=8080 SERVE_DIR=./public_html go run .
```

- Web: http://localhost:80 (o `PORT` definido)
- Componentes: `GET /api/components/{skills,projects,experience,contact,chat}` + `POST /api/components/chat{,/stream}`

## Uso con Docker

```bash
docker compose up -d --build
docker compose logs -f
```

- Desarrollo local (puerto 8001, sin Traefik): `docker compose -f docker-compose.dev.yml up --build` → http://localhost:8001. Fichero ignorado por git y sftp; no toca producción.

## Env vars

| Var | Uso | Defecto |
|---|---|---|
| `PORT` | puerto escucha | `80` |
| `SERVE_DIR` | raíz estática | `/var/www/html` |
| `SMTP_HOST` / `SMTP_PORT` / `SMTP_USER` / `SMTP_PASS` / `SMTP_TO` | envío formulario contacto | — |
| `TURNSTILE_SECRET` | verificación anti-spam contacto | — |
| `RECAPTCHA_SKIP=true` | bypassea Turnstile (solo dev) | — |
| `OPENAI_API_KEY` | proxy del chat (Responses API) | — (chat desactivado sin ella) |
| `OPENAI_PROMPT_ID` | prompt almacenado del chat | `pmpt_6aa9…2acc7` |
| `OPENAI_PROMPT_VERSION` | versión fija del prompt; vacío = versión default del dashboard | — |
| `OPENAI_PROMPT_VARIABLES` | JSON con variables `{{...}}` de la plantilla (solo si las usa) | — |
| `OPENAI_LANG_VAR` | variable adicional que recibe el idioma del visitante; el backend siempre envía también `lang` y `datetime` requeridas por el prompt publicado | `idioma` |

## Añadir componente

1. Crear `components/<nombre>/` con handler `Handler(w, r)`.
2. Registrar en `components/registry.go`: `"nombre": nombre.Handler`.
3. Consumir desde HTML: `hx-get="/api/components/<nombre>"`.

## Deploy

- Traefik: host `nestorpons.com`, entrypoints `web,websecure`, `certresolver=myresolver`.
- Volúmenes: `./public_html:/var/www/html`.

## Seguridad

- Secretos (`SMTP_*`, `TURNSTILE_SECRET`) solo en `.env` local (ignorado); el repo lleva `.env.example`.
- Nunca commitear `.env`.
