# Repository Guidance

## Structure

- This is one Go 1.22 module (`nestorgo-landing`), not a frontend package workspace.
- `main.go` serves `SERVE_DIR` and mounts component endpoints at `/api/components/{skills,projects,experience,contact}`.
- `public_html/` is the static HTMX frontend; component endpoints return HTML fragments rendered from files under `components/`.
- New component endpoints must be added to `components/registry.go`; data-driven components keep their JSON and HTML beside the handler.

## Local Development

- Run from the repository root: `SERVE_DIR=public_html PORT=8080 go run .` (handlers read `components/...` using relative paths).
- Contact POSTs require `TURNSTILE_SECRET` and `SMTP_*`; use `RECAPTCHA_SKIP=true` only for local form testing and never use production credentials locally.
- Verify changes with `go test ./...` (there are currently no test files) and, when changing Go, `gofmt -w main.go components/*.go components/*/*.go`.
- Smoke-test the running server with `curl http://localhost:8080/` and the relevant `/api/components/<name>` endpoint.

## Docker

- `docker compose up --build -d` is the deployment path. It requires the external Docker network `traefik-net`.
- The final image contains the Go binary but not the website assets; Compose bind-mounts `public_html/` and `components/` into the container, so preserve those mounts when changing deployment.
- Keep SMTP and Turnstile values in environment configuration; do not copy credentials into source, documentation, or frontend assets.
