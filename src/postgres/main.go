package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"postgres_server/config"
	"postgres_server/security"
	"postgres_server/tools"
)

func main() {
	config.Init()
	startedAt := time.Now().UTC()

	authMode := security.MCPAuthModeFromEnv()
	adminEnabled := authMode == security.AuthModeIssuedToken && security.AdminAPITokenFromEnv() != ""
	exchangeCfg := security.TencentCloudTokenExchangeConfigFromEnv(authMode)
	needsTokenStore := authMode == security.AuthModeIssuedToken || adminEnabled || exchangeCfg.Enabled

	var (
		sqliteStore        *security.SQLiteTokenStore
		tokenStore         security.TokenStore
		tokenIssuer        *security.TokenIssuer
		issuerCfg          security.TokenIssuerConfig
		credentialCipher   *security.CredentialCipher
		credentialProvider security.CredentialProvider
		tokenExchange      *security.TencentCloudTokenExchangeService
		err                error
	)

	if needsTokenStore {
		issuerCfg = security.TokenIssuerConfigFromEnv()
		sqliteStore, err = security.NewSQLiteTokenStoreFromEnv()
		if err != nil {
			log.Fatalf("init token store failed: %v", err)
		}
		tokenStore = sqliteStore
		tokenIssuer = security.NewTokenIssuer(sqliteStore, issuerCfg)
	}

	switch authMode {
	case security.AuthModeIssuedToken:
		credentialCipher, err = security.NewCredentialCipherFromEnv()
		if err != nil {
			log.Fatal(err)
		}
		credentialProvider, err = security.NewTokenBoundCredentialProvider(sqliteStore, credentialCipher)
		if err != nil {
			log.Fatal(err)
		}
	case security.AuthModeRequestCredential:
		credentialProvider = security.NewRequestHeaderCredentialProvider()
	default:
		credentialProvider, err = security.NewStaticCredentialProviderFromEnv()
		if err != nil {
			log.Fatal(err)
		}
	}

	if exchangeCfg.Enabled {
		if credentialCipher == nil {
			credentialCipher, err = security.NewCredentialCipherFromEnv()
			if err != nil {
				log.Fatal(err)
			}
		}
		tokenExchange, err = security.NewTencentCloudTokenExchangeService(tokenIssuer, sqliteStore, credentialCipher, exchangeCfg)
		if err != nil {
			log.Fatal(err)
		}
	}

	guard := security.NewGuard()
	mcpServerName := "mcp-server-postgres"
	mcpsvr := server.NewMCPServer(
		"腾讯云 Postgres MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	toolCount := 0
	if config.IsEnabled("instance") {
		tools.RegisterInstanceTools(mcpsvr, credentialProvider, guard)
		toolCount += 15
	}
	if config.IsEnabled("account") {
		tools.RegisterAccountTools(mcpsvr, credentialProvider, guard)
		toolCount += 6
	}
	if config.IsEnabled("database") {
		tools.RegisterDatabaseTools(mcpsvr, credentialProvider, guard)
		toolCount += 4
	}
	if config.IsEnabled("parameter") {
		tools.RegisterParameterTools(mcpsvr, credentialProvider, guard)
		toolCount += 5
	}
	if config.IsEnabled("backup") {
		tools.RegisterBackupTools(mcpsvr, credentialProvider, guard)
		toolCount += 8
	}
	if config.IsEnabled("monitoring") {
		tools.RegisterMonitoringTools(mcpsvr, credentialProvider, guard)
		toolCount += 3
	}
	if config.IsEnabled("network") {
		tools.RegisterNetworkTools(mcpsvr, credentialProvider, guard)
		toolCount += 4
	}
	if config.IsEnabled("readonly") {
		tools.RegisterReadonlyTools(mcpsvr, credentialProvider, guard)
		toolCount += 2
	}
	tools.RegisterSSLTools(mcpsvr, credentialProvider, guard)
	toolCount += 1
	log.Printf("Total tools registered: %d", toolCount)

	authenticator, err := security.NewMCPAuthenticator(authMode, tokenStore, issuerCfg.Pepper)
	if err != nil {
		log.Fatalf("init auth failed: %v", err)
	}

	httpEndpoint := transportEndpointFromEnv()
	serverPort := getEnv("MCP_SERVER_PORT", getEnv("MCP_SERVER_SSE_PORT", "9000"))
	bindHost := getEnv("MCP_SERVER_BIND_HOST", "127.0.0.1")
	listenAddr := fmt.Sprintf("%s:%s", bindHost, serverPort)
	serverURL := buildServerURL(bindHost, serverPort, httpEndpoint)
	statelessHTTP := getEnvBool("MCP_STREAMABLE_HTTP_STATELESS", true)

	transportServer := server.NewStreamableHTTPServer(mcpsvr,
		server.WithEndpointPath(httpEndpoint),
		server.WithStateLess(statelessHTTP),
	)

	mux := http.NewServeMux()
	security.RegisterHealthRoutes(mux, security.HealthStatus{
		Service:   mcpServerName,
		Version:   "1.0.0",
		StartedAt: startedAt,
	})
	if adminEnabled {
		security.RegisterAdminRoutes(mux, tokenIssuer, tokenStore)
	}
	if tokenExchange != nil {
		security.RegisterTokenExchangeRoutes(mux, tokenExchange, security.TokenExchangeBootstrapConfig{
			MCPServerName: mcpServerName,
			ServerURL:     serverURL,
		})
	}
	mux.Handle(httpEndpoint, security.WrapMCPAuth(authenticator, transportServer))

	httpServer := &http.Server{
		Addr:              listenAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       2 * time.Minute,
	}

	log.Printf("MCP streamable-http server listening on %s%s", listenAddr, httpEndpoint)
	log.Printf("MCP transport: streamable-http (stateless=%t)", statelessHTTP)
	switch authMode {
	case security.AuthModeSharedToken:
		log.Printf("MCP auth mode: shared-token (%s)", security.APIAuthFallbackHeader)
	case security.AuthModeIssuedToken:
		log.Printf("MCP auth mode: issued-token")
		log.Printf("Token store: %s", security.TokenStorePathFromEnv())
		log.Printf("Credential source: token-bound dynamic credentials")
	case security.AuthModeRequestCredential:
		log.Printf("MCP auth mode: request-credential")
		log.Printf("Credential source: %s / %s / %s", security.RequestSecretIDHeader, security.RequestSecretKeyHeader, security.RequestSessionTokenHeader)
		if security.RequestCredentialValidateIdentityFromEnv() {
			log.Printf("Request credential identity validation: enabled via STS GetCallerIdentity")
		} else {
			log.Printf("warning: request credential identity validation is disabled")
		}
	default:
		log.Printf("MCP auth mode: none")
	}
	if adminEnabled {
		log.Printf("Admin token API enabled on /admin/tokens")
	} else if authMode == security.AuthModeIssuedToken {
		log.Printf("warning: issued-token mode is enabled but %s is empty; admin token API is disabled", security.AdminAPITokenEnv)
	}
	if tokenExchange != nil {
		log.Printf("TencentCloud token exchange enabled on /auth/token-exchange/tencentcloud (%s)", exchangeCfg.Mode)
	}
	log.Printf("Health check endpoints enabled on /healthz and /readyz")

	outputMCPServerConfig(mcpServerName, serverURL, authMode)
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server error: %v", err)
	}
}

func outputMCPServerConfig(mcpServerName, serverURL string, authMode security.AuthMode) {
	serverConfig := map[string]any{
		"url":  serverURL,
		"type": "streamable-http",
	}
	if headers, ok := clientHeaderPlaceholders(authMode); ok {
		serverConfig["headers"] = headers
	}

	configPayload := map[string]any{
		"mcpServers": map[string]any{
			mcpServerName: serverConfig,
		},
	}

	jsonOutput, _ := json.MarshalIndent(configPayload, "", " ")
	fmt.Println("=== MCP Server Configuration ===")
	fmt.Println("Copy the following configuration to your MCP client:")
	fmt.Println()
	fmt.Println(string(jsonOutput))
	fmt.Println()
	switch authMode {
	case security.AuthModeSharedToken:
		fmt.Println("Authentication is enabled for this server.")
		fmt.Printf("Replace <MCP_API_TOKEN> with the shared token, or use header %s if your client does not support Bearer config.\n", security.APIAuthFallbackHeader)
		fmt.Println()
	case security.AuthModeIssuedToken:
		fmt.Println("Issued-token auth is enabled for this server.")
		fmt.Println("Use POST /auth/token-exchange/tencentcloud to exchange TencentCloud credentials for a local MCP access token, or let an admin create one via POST /admin/tokens.")
		fmt.Println()
	case security.AuthModeRequestCredential:
		fmt.Println("Request-credential auth is enabled for this server.")
		fmt.Printf("Every MCP request must carry %s and %s headers. %s is optional for STS temporary credentials.\n", security.RequestSecretIDHeader, security.RequestSecretKeyHeader, security.RequestSessionTokenHeader)
		fmt.Println("Use HTTPS in production and never place secret material in URL query parameters.")
		fmt.Println()
	}
	fmt.Println("The server is now ready to accept connections.")
}

func clientHeaderPlaceholders(authMode security.AuthMode) (map[string]string, bool) {
	switch authMode {
	case security.AuthModeSharedToken:
		return map[string]string{"Authorization": "Bearer <MCP_API_TOKEN>"}, true
	case security.AuthModeIssuedToken:
		return map[string]string{"Authorization": "Bearer <MCP_ACCESS_TOKEN>"}, true
	case security.AuthModeRequestCredential:
		return security.RequestCredentialPlaceholderHeaders(), true
	default:
		return nil, false
	}
}

func buildServerURL(bindHost, port, endpoint string) string {
	if publicURL := strings.TrimSpace(os.Getenv("MCP_SERVER_PUBLIC_URL")); publicURL != "" {
		return publicURL
	}
	host := strings.TrimSpace(bindHost)
	switch host {
	case "", "0.0.0.0", "::":
		host = "127.0.0.1"
	}
	return fmt.Sprintf("http://%s:%s%s", host, port, endpoint)
}

func transportEndpointFromEnv() string {
	if endpoint := getEnv("MCP_SERVER_HTTP_ENDPOINT", ""); endpoint != "" {
		return normalizeEndpointPath(endpoint)
	}
	if endpoint := getEnv("MCP_SERVER_SSE_ENDPOINT", ""); endpoint != "" {
		return normalizeEndpointPath(endpoint)
	}
	return "/mcp"
}

func normalizeEndpointPath(path string) string {
	trimmed := strings.Trim(strings.TrimSpace(path), "/")
	if trimmed == "" {
		return "/mcp"
	}
	return "/" + trimmed
}

func getEnv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	switch strings.ToLower(v) {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return fallback
	}
}
