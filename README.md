# nestorgo.com — Portfolio Néstor Pons

Landing + portfolio personal. Backend mínimo en Go que sirve estático y expone componentes HTML vía htmx.

- URL: https://nestorpons.com
- Stack: Go 1.22 · htmx · Tailwind CSS · Docker + Traefik

## Estructura

- `main.go` — estático (`SERVE_DIR`) + registro `/api/components/*`
- `components/` — `registry.go` + 1 carpeta por componente (`*.go` + `*.json` + `*.html`)
  - `skills` / `projects` / `experience` — leen `*.json`, renderizan `*.html`
  - `contact` — `GET` sirve `form.html`, `POST` verifica Turnstile y envía por SMTP
- `public_html/` — `index.html` + `assets/` (`main.js`, `translator.js`, `styles.css`, `images/`, `cv.pdf`)
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
- Componentes: `GET /api/components/{skills,projects,experience,contact}`

## Uso con Docker

```bash
docker compose up -d --build
docker compose logs -f
```

## Env vars

| Var | Uso | Defecto |
|---|---|---|
| `PORT` | puerto escucha | `80` |
| `SERVE_DIR` | raíz estática | `/var/www/html` |
| `SMTP_HOST` / `SMTP_PORT` / `SMTP_USER` / `SMTP_PASS` / `SMTP_TO` | envío formulario contacto | — |
| `TURNSTILE_SECRET` | verificación anti-spam contacto | — |
| `RECAPTCHA_SKIP=true` | bypassea Turnstile (solo dev) | — |

## Añadir componente

1. Crear `components/<nombre>/` con handler `Handler(w, r)`.
2. Registrar en `components/registry.go`: `"nombre": nombre.Handler`.
3. Consumir desde HTML: `hx-get="/api/components/<nombre>"`.

## Deploy

- Traefik: host `nestorpons.com`, entrypoints `web,websecure`, `certresolver=myresolver`.
- Volúmenes: `./public_html:/var/www/html`.

## Seguridad antes de subir a GitHub

- `docker-compose.yml` contiene `SMTP_*` y `TURNSTILE_SECRET` en claro → mover a `.env` y no commitear.
- Añadir `.gitignore`: `.env`, `*.log`, binarios (`server`).
