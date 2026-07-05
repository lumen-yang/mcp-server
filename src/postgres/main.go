package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"postgres_server/config"
	"postgres_server/security"
	"postgres_server/tools"
)

func main() {
	// 初始化功能组开关
	config.Init()

	// 创建 MCP Server
	mcpServerName := "mcp-server-postgres"
	mcpsvr := server.NewMCPServer(
		"腾讯云 Postgres MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	// 创建腾讯云凭证
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	// 创建 Guard 安全中间件
	guard := security.NewGuard()

	// 按功能组注册工具
	toolCount := 0

	if config.IsEnabled("instance") {
		tools.RegisterInstanceTools(mcpsvr, credential, guard)
		toolCount += 12
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
		toolCount += 4
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
	ssePort := getEnv("MCP_SERVER_SSE_PORT", "9000")

	sseServer := server.NewSSEServer(mcpsvr,
		server.WithSSEEndpoint(sseEndpoint),
		server.WithMessageEndpoint(messageEndpoint),
		server.WithAppendQueryToMessageEndpoint())

	log.Printf("SSE server listening on :%s", ssePort)
	serverURL := fmt.Sprintf("http://127.0.0.1:%s%s", ssePort, sseEndpoint)
	outputMCPServerConfig(mcpServerName, serverURL)
	if err := sseServer.Start(":" + ssePort); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func outputMCPServerConfig(mcpServerName, serverURL string) {
	config := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			mcpServerName: map[string]interface{}{
				"url":  serverURL,
				"type": "sse",
			},
		},
	}

	jsonOutput, _ := json.MarshalIndent(config, "", " ")
	fmt.Println("=== MCP Server Configuration ===")
	fmt.Println("Copy the following configuration to your MCP client:")
	fmt.Println()
	fmt.Println(string(jsonOutput))
	fmt.Println()
	fmt.Println("The server is now ready to accept connections.")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
