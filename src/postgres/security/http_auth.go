package security

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/client/transport"
)

const (
	// APIAuthTokenEnv 控制 MCP SSE / message 入口鉴权。
	APIAuthTokenEnv = "MCP_API_TOKEN"

	// APIAuthFallbackHeader 是兼容用的自定义头；优先推荐 Authorization: Bearer。
	APIAuthFallbackHeader = "X-MCP-API-Token"
)

// APIAuthTokenFromEnv 读取当前进程配置的 MCP API Token。
func APIAuthTokenFromEnv() string {
	return strings.TrimSpace(os.Getenv(APIAuthTokenEnv))
}

// WrapHTTPAuth 为 HTTP Handler 增加基于共享 Token 的入口鉴权。
// 当 expectedToken 为空时，返回原始 handler（即不启用鉴权）。
func WrapHTTPAuth(expectedToken string, next http.Handler) http.Handler {
	expectedToken = strings.TrimSpace(expectedToken)
	if expectedToken == "" || next == nil {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !authorizedRequest(r, expectedToken) {
			w.Header().Set("WWW-Authenticate", `Bearer realm="mcp-server-postgres"`)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// MCPClientOptionsFromEnv 为本地验证/写测客户端自动注入鉴权 header。
func MCPClientOptionsFromEnv() []transport.ClientOption {
	token := APIAuthTokenFromEnv()
	if token == "" {
		return nil
	}
	return []transport.ClientOption{
		transport.WithHeaders(map[string]string{
			"Authorization": "Bearer " + token,
		}),
	}
}

func authorizedRequest(r *http.Request, expectedToken string) bool {
	if r == nil {
		return false
	}
	return secureTokenEqual(extractRequestToken(r), expectedToken)
}

func extractRequestToken(r *http.Request) string {
	if r == nil {
		return ""
	}
	if bearerToken := extractBearerToken(r.Header.Get("Authorization")); bearerToken != "" {
		return bearerToken
	}
	return strings.TrimSpace(r.Header.Get(APIAuthFallbackHeader))
}

func extractBearerToken(authHeader string) string {
	fields := strings.Fields(strings.TrimSpace(authHeader))
	if len(fields) != 2 {
		return ""
	}
	if !strings.EqualFold(fields[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(fields[1])
}

func secureTokenEqual(got, expected string) bool {
	got = strings.TrimSpace(got)
	expected = strings.TrimSpace(expected)
	if got == "" || expected == "" || len(got) != len(expected) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(got), []byte(expected)) == 1
}
