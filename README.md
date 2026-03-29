# OIDC Mock Issuer

An OIDC issuer mock for testing OIDC clients. Provides discovery, JWKS, and token endpoints. Supports issuing valid and invalid tokens.

## Features

- OIDC discovery and JWKS endpoints
- Token issuance with configurable claims and validity
- CLI flags for port and issuer URL
- Build for amd64 and arm64

## Usage

### Build

```go
go build -o oidc-mock-issuer .
```

### Run

```go
go run . --port 8080 --issuer http://localhost:8080
```

### Endpoints

- `/.well-known/openid-configuration` — OIDC discovery
- `/.well-known/jwks.json` — JWKS
- `/token` — POST a JSON body to issue a token (see below)

### Token endpoint

`POST /token`

**Request body**

| Field | Type | Description |
|-------|------|-------------|
| `valid` | bool | `true` issues a signed RS256 JWT; `false` returns a static invalid token string. Defaults to `true`. |
| `claims` | object | Additional claims merged into the JWT payload. Any JSON-serialisable key/value pairs are accepted. |

The request body is optional. Omitting it is equivalent to `{"valid": true}`.

The issued JWT always includes `iat` (now), `exp` (now + 1 h), and `iss` (the configured `--issuer` URL). Custom claims in `claims` override those defaults if the same key is provided.

**Response**

```json
{ "access_token": "<jwt>", "token_type": "Bearer" }
```

**Examples**

Issue a valid token with standard OIDC claims:
```bash
curl -s -X POST http://localhost:8080/token \
  -H "Content-Type: application/json" \
  -d '{"valid": true, "claims": {"sub": "user-123", "iss": "http://localhost:8080", "aud": "my-client"}}' \
  | jq .
```

Issue a valid token with arbitrary custom claims:

```bash
curl -s -X POST http://localhost:8080/token \
  -H "Content-Type: application/json" \
  -d '{"valid": true, "claims": {"sub": "svc-account", "roles": ["admin", "reader"], "tenant": "acme"}}' \
  | jq .
```

Issue an invalid (unsigned) token to test rejection scenarios:

```bash
curl -s -X POST http://localhost:8080/token \
  -H "Content-Type: application/json" \
  -d '{"valid": false}' \
  | jq .
```

Extract just the token value for use in subsequent requests:

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/token \
  -H "Content-Type: application/json" \
  -d '{"valid": true, "claims": {"sub": "user-123"}}' \
  | jq -r .access_token)

curl -H "Authorization: Bearer $TOKEN" http://your-service/protected
```
