# Codex Project Instructions

## Backend Rules

When modifying code under `backend/`, Codex must follow the backend rules in:

```text
.codex/rules/backend/AI_RULES.md
```

Before making backend changes, read that file and apply its constraints for architecture, routing, handlers, database migrations, permissions, errors, testing, and final reporting.

Key requirements:

- Keep backend business code under `backend/internal`.
- Follow the existing Go + Gin + GORM + PostgreSQL architecture.
- Do not bypass the shared `response`, `errs`, `middleware`, or migration conventions.
- Do not introduce unrelated refactors or broad framework changes unless explicitly requested.
- For backend code changes, run `gofmt` and `go test ./...` from `backend` when applicable, and report the result.
