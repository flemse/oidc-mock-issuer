package main

import (
	"log"
	"net/http"
	"os"

	"github.com/flemse/oidc-mock-issuer/internal/issuer"
	"github.com/flemse/oidc-mock-issuer/internal/server"
	"github.com/flemse/oidc-mock-issuer/internal/tokens"
	"github.com/spf13/cobra"
)

func main() {
	var issuerURL string
	var port string

	rootCmd := &cobra.Command{
		Use:   "oidc-mock-issuer",
		Short: "A mock OIDC issuer for testing",
		Run: func(cmd *cobra.Command, args []string) {
			ks, err := tokens.NewKeySet()
			if err != nil {
				log.Fatalf("failed to generate key set: %v", err)
			}
			mux := http.NewServeMux()
			mux.HandleFunc("/.well-known/openid-configuration", issuer.DiscoveryHandler(issuerURL))
			mux.HandleFunc("/.well-known/jwks.json", issuer.JWKSHandler(ks))
			mux.HandleFunc("/token", issuer.TokenHandler(ks, issuerURL))
			server.Start(":"+port, mux)
		},
	}
	rootCmd.Flags().StringVar(&issuerURL, "issuer", "http://localhost:8080", "Issuer URL")
	rootCmd.Flags().StringVar(&port, "port", "8080", "Port to listen on")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
