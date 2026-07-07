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
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"postgres_server/config"
	"postgres_server/security"
	"postgres_server/tools"
)

func main() {
	// 初始化功能组开关
	config.Init()

	secretID, secretKey, usingLegacyEnv := resolveCloudCredentials()
	if secretID == "" || secretKey == "" {
		log.Fatal("missing credentials: set MCP_SECRET_ID/MCP_SECRET_KEY")
	}
	if usingLegacyEnv {
		log.Printf("warning: using legacy TENCENTCLOUD_SECRET_* envs; please migrate to MCP_SECRET_*")
	}
	authToken := security.APIAuthTokenFromEnv()

	// 创建 MCP Server
	mcpServerName := "mcp-server-postgres"
	mcpsvr := server.NewMCPServer(
		"腾讯云 Postgres MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	// 创建腾讯云凭证
	credential := common.NewCredential(secretID, secretKey)

	// 创建 Guard 安全中间件
	guard := security.NewGuard()

	// 按功能组注册工具
	toolCount := 0

	if config.IsEnabled("instance") {
		tools.RegisterInstanceTools(mcpsvr, credential, guard)
		toolCount += 15
	}
	if config.IsEnabled("account") {
		tools.RegisterAccountTools(mcpsvr, credential, guard)
		toolCount += 6
	}
	if config.IsEnabled("database") {
		tools.RegisterDatabaseTools(mcpsvr, credential, guard)
		toolCount += 4
	}
	if config.IsEnabled("parameter") {
		tools.RegisterParameterTools(mcpsvr, credential, guard)
		toolCount += 5
	}
	if config.IsEnabled("backup") {
		tools.RegisterBackupTools(mcpsvr, credential, guard)
		toolCount += 8
	}
	if config.IsEnabled("monitoring") {
		tools.RegisterMonitoringTools(mcpsvr, credential, guard)
		toolCount += 3
	}
	if config.IsEnabled("network") {
		tools.RegisterNetworkTools(mcpsvr, credential, guard)
		toolCount += 4
	}
	if config.IsEnabled("readonly") {
		tools.RegisterReadonlyTools(mcpsvr, credential, guard)
		toolCount += 2
	}
	// SSL 工具始终注册（1个）
	tools.RegisterSSLTools(mcpsvr, credential, guard)
	toolCount += 1

	log.Printf("Total tools registered: %d", toolCount)

	// 启动 SSE Server
	sseEndpoint := getEnv("MCP_SERVER_SSE_ENDPOINT", "/sse")
	messageEndpoint := getEnv("MCP_SERVER_MESSAGE_ENDPOINT", "/message")
	ssePort := getEnv("MCP_SERVER_PORT", getEnv("MCP_SERVER_SSE_PORT", "9000"))
	bindHost := getEnv("MCP_SERVER_BIND_HOST", "127.0.0.1")
	listenAddr := fmt.Sprintf("%s:%s", bindHost, ssePort)

	sseServer := server.NewSSEServer(mcpsvr,
		server.WithSSEEndpoint(sseEndpoint),
		server.WithMessageEndpoint(messageEndpoint),
		server.WithAppendQueryToMessageEndpoint())

	httpServer := &http.Server{
		Addr:              listenAddr,
		Handler:           security.WrapHTTPAuth(authToken, sseServer),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       2 * time.Minute,
	}

	log.Printf("SSE server listening on %s", listenAddr)
	if authToken != "" {
		log.Printf("MCP API token auth enabled (Authorization: Bearer <token> or %s)", security.APIAuthFallbackHeader)
	}
	serverURL := buildServerURL(bindHost, ssePort, sseEndpoint)
	outputMCPServerConfig(mcpServerName, serverURL, authToken != "")
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server error: %v", err)
	}
}

func outputMCPServerConfig(mcpServerName, serverURL string, authEnabled bool) {
	serverConfig := map[string]interface{}{
		"url":  serverURL,
		"type": "sse",
	}
	if authEnabled {
		serverConfig["headers"] = map[string]string{
			"Authorization": "Bearer <MCP_API_TOKEN>",
		}
	}

	config := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			mcpServerName: serverConfig,
		},
	}

	jsonOutput, _ := json.MarshalIndent(config, "", " ")
	fmt.Println("=== MCP Server Configuration ===")
	fmt.Println("Copy the following configuration to your MCP client:")
	fmt.Println()
	fmt.Println(string(jsonOutput))
	fmt.Println()
	if authEnabled {
		fmt.Println("Authentication is enabled for this server.")
		fmt.Printf("Replace <MCP_API_TOKEN> with the shared token, or use header %s if your client does not support Bearer config.\n", security.APIAuthFallbackHeader)
		fmt.Println()
	}
	fmt.Println("The server is now ready to accept connections.")
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

func resolveCloudCredentials() (secretID, secretKey string, usingLegacyEnv bool) {
	secretID = strings.TrimSpace(os.Getenv("MCP_SECRET_ID"))
	secretKey = strings.TrimSpace(os.Getenv("MCP_SECRET_KEY"))
	if secretID == "" {
		if legacy := strings.TrimSpace(os.Getenv("TENCENTCLOUD_SECRET_ID")); legacy != "" {
			secretID = legacy
			usingLegacyEnv = true
		}
	}
	if secretKey == "" {
		if legacy := strings.TrimSpace(os.Getenv("TENCENTCLOUD_SECRET_KEY")); legacy != "" {
			secretKey = legacy
			usingLegacyEnv = true
		}
	}
	return secretID, secretKey, usingLegacyEnv
}

func getEnv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
