FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build \
    -ldflags="-s -w -extldflags '-static'" \
    -trimpath \
    -o /oidc-mock-issuer .

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /oidc-mock-issuer /oidc-mock-issuer

EXPOSE 8080

ENTRYPOINT ["/oidc-mock-issuer"]
