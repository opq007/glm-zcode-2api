# glm-zcode-2api

Local reverse proxy that turns a ZCode plan into an OpenAI-compatible API. The gateway accepts OpenAI `/v1/chat/completions` traffic and translates it to the upstream Anthropic Messages API (SSE streaming). Entry: `cmd/server -config <file>`, default listener `0.0.0.0:7864`.

## Hard rules

- **Pure Go standard library only.** No third-party dependencies — that is a stated project constraint (single binary).
- **Never commit or print real credentials** (API keys, JWTs, account IDs, `client.key`, gateway tokens). Keys are never written to logs or to OMP config. `.gitignore` protects `config.local*.json`, `*.local.json`, `device.key`, `bin/`.
- **Do not re-add an `Authorization: Bearer` mirror alongside `x-api-key`.** It was measured to cost 5x credits and was removed; `x-api-key` alone is the correct auth channel to the upstream.
- Keep `README.md` and `README.zh-CN.md` in sync when changing user-facing docs.

## Commands

```sh
go build ./...                                   # compile everything
go run ./cmd/server -config config.local.json    # run the gateway
gofmt -w . && go vet ./... && go test ./...      # pre-submission check (required)
go test ./internal/convert/ -run TestName        # focused test
```

`scripts/omp-gateway.py` is the OMP (oh-my-pi) launcher: `start/stop/status/restart/token`; it builds to `bin/glm-zcode-2api` and keeps state under `~/.local/state/glm-zcode-2api/`. It uses POSIX `fcntl`/`0o600` file semantics — verify behavior before changing its locking or secret handling.

## Package map (respect the boundaries)

- `internal/openai/`, `internal/anthropic/` — wire types only. Never leak Anthropic types to clients or OpenAI types upstream; all cross-translation lives in `internal/convert/`.
- `internal/convert/` — request/response/SSE translation plus `ReplayCache`: signed thinking blocks returned by upstream must be replayed on the next tool-use turn (Anthropic protocol requirement). Don't drop signature replay when touching translate logic.
- `internal/upstream/` — HTTP client: SSE parsing, idle watchdog, error classes.
- `internal/server/` — the only HTTP surface (routes, `x-api-key` auth, request logging). API key is mandatory because the listener is LAN-facing and carries account credentials.
- `internal/credential/` — read-only discovery of apiKey/baseURL from `~/.zcode/v2/config.json` (cached, re-read on file change). It must never modify ZCode app files.
- `internal/config/` — JSON file first, then `Z2A_*` env overrides (only non-empty vars win).
- `docs/` — GitHub Pages site; `skills/glm-zcode-2api-setup/` — agent setup skill.

## Gotchas

- `upstream.mimic_client` sends the ZCode client's attribution headers so plan promotions (off-peak discounts, free flash windows) apply to proxied traffic; `gateway_origin` routes coding-plan traffic through the ZCode platform gateway where entitlements are validated. Empty `gateway_origin` = talk to the provider endpoint directly. Changing either affects billing/entitlement behavior — check `AI_SETUP.md` and git history first.
- Targets macOS / Linux / Windows (pure-Go server); the Python launcher is the Unix-side tooling.
- Restart the gateway after re-signing into ZCode or switching plans (credentials are read from the live ZCode config).
- Structured log lines use the `event=<name> key=value` shape (`log` package to stderr, UTC). Match it when adding logging in `internal/server/`.

## Read before sensitive changes

- `AI_SETUP.md` — install, health check, and OMP integration (`~/.omp/agent/models.yml` skeleton, token command).
- `README.md` — user-facing config table, endpoints, roadmap.
