# Go OIDC Mock Issuer: Project Plan

## Overview
A single-binary OIDC issuer mock written in Go, providing:
- OIDC discovery (/.well-known/openid-configuration)
- JWKS endpoint (/.well-known/jwks.json)
- Token issuance endpoint (configurable for valid/invalid tokens)
- CLI interface (cobra-cli)
- Build targets: amd64 and arm64
- CI/CD: GitHub Actions for build, test, lint (golangci-lint), and release
- Documentation

---


## 1. Project Structure
- `/cmd/oidc-mock-issuer/` — Main CLI entrypoint
- `/internal/issuer/` — OIDC logic (discovery, JWKS, token issuance)
- `/internal/server/` — HTTP server and routing
- `/internal/tokens/` — Token generation/validation helpers
- `/pkg/` — (Optional) Public Go packages
- `.github/workflows/` — GitHub Actions workflows
- `README.md` — Project documentation

**Testing:**
- Tests should be stored alongside the corresponding `pkg` and `internal` files (e.g., `issuer_test.go` next to `issuer.go`).

---

## 2. Features
- Serve OIDC discovery and JWKS endpoints
- Issue tokens via HTTP endpoint, with parameters for claims, expiry, validity, etc.
- Option to issue intentionally invalid tokens (e.g., bad signature, wrong claims)
- CLI flags for port, issuer URL, keys, etc.
- Self-contained: no external dependencies except cobra and golangci-lint

---

## 3. CLI & Binary Builds
- Use cobra-cli for command-line interface
- Build static binaries for:
  - linux/amd64
  - linux/arm64
  - darwin/amd64
  - darwin/arm64
- Store built binaries as GitHub Releases artifacts

---

## 4. Endpoints
- `/.well-known/openid-configuration` — OIDC discovery
- `/.well-known/jwks.json` — JWKS endpoint
- `/token` — Issue tokens (accepts POST with parameters for claims, validity, etc.)

---

## 5. CI/CD (GitHub Actions)
- **Build:** Cross-compile for all targets
- **Test:** Run Go tests
- **Lint:** Run golangci-lint
- **Release:** On tag, upload binaries as release assets

---

## 6. Linting
- Use golangci-lint (configurable via `.golangci.yml`)

---

## 7. Documentation
- `README.md` with:
  - Project overview
  - Usage instructions (CLI, endpoints)
  - Example requests/responses
  - Build and release process
  - Contribution guidelines

---

## 8. Next Steps
1. Scaffold Go project structure
2. Implement CLI and HTTP server
3. Implement OIDC endpoints and token logic
4. Add tests
5. Set up GitHub Actions workflows
6. Write documentation

---

*This plan will guide the implementation of the Go OIDC mock issuer project, ensuring robust features, CI/CD, and clear documentation.*
