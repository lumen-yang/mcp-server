package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"postgres_server/security"
)

// ToolHandler 执行 API 调用并返回 JSON 响应字符串
type ToolHandler func(client *postgres.Client, args map[string]interface{}) (string, error)

// registerTool 创建并注册一个 MCP Tool，封装公共的 handler 样板代码
// 参数:
//   - s: MCP Server
//   - cred: 腾讯云凭证
//   - g: Guard 安全中间件
//   - name: 工具名（不含 postgres- 前缀，自动添加）
//   - description: 工具描述
//   - guardLevel: guard 级别
//   - params: 工具参数选项（不含 region，自动添加）
//   - handler: API 调用处理函数
func registerTool(
	s *server.MCPServer,
	cred *common.Credential,
	g *security.Guard,
	name string,
	description string,
	guardLevel security.GuardLevel,
	params []mcp.ToolOption,
	handler ToolHandler,
) {
	// 为需要确认的操作添加 confirm 参数
	if security.NeedsConfirm(guardLevel) {
		params = append(params, mcp.WithBoolean("confirm",
			mcp.Description("确认真执行此操作。设为 true 以确认。")))
	}

	// 构建工具选项：description + region + 用户自定义参数
	toolOpts := []mcp.ToolOption{
		mcp.WithDescription(description),
		mcp.WithString("region", mcp.Required(), mcp.Description("地域")),
	}
	toolOpts = append(toolOpts, params...)

	tool := mcp.NewTool("postgres-"+name, toolOpts...)

	// 该工具是否本身声明了 DBInstanceId 参数：只有声明了此参数的工具，
	// 才允许 guard 在未显式传参时自动注入 scope 限定的实例ID，避免向
	// 不认识该字段的 SDK 请求结构体（如查规格/查版本等全局目录查询）盲注入。
	_, hasInstanceIdField := tool.InputSchema.Properties["DBInstanceId"]

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 1. 提取 region（默认 ap-guangzhou）
		region := "ap-guangzhou"
		arguments := request.GetArguments()
		if v, ok := arguments["region"].(string); ok && v != "" {
			region = v
		}
		delete(arguments, "region")

		// 2. Guard 校验（只读模式/资源 scoping）
		if g != nil {
			if err := g.Check(name, region, arguments, guardLevel, hasInstanceIdField); err != nil {
				return mcp.NewToolResultText(fmt.Sprintf(`{"code":403,"error":"%s"}`, err.Error())), nil
			}
		}

		// 3. 确认检查（写操作需要 confirm=true）
		if security.NeedsConfirm(guardLevel) {
			confirmed := false
			if v, ok := arguments["confirm"].(bool); ok {
				confirmed = v
			}
			delete(arguments, "confirm")

			if !confirmed {
				warning := security.GetGuardWarning(guardLevel, name)
				return mcp.NewToolResultText(fmt.Sprintf(`{"code":403,"warning":"%s","require_confirm":true}`, warning)), nil
			}
		}

		// 4. 创建 postgres client
		cpf := profile.NewClientProfile()
		cpf.Debug = getDebugFlag()
		client, err := postgres.NewClient(cred, region, cpf)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf(`{"code":500,"error":"create client failed: %s"}`, err.Error())), nil
		}

		// 5. 执行 handler
		result, err := handler(client, arguments)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf(`{"code":500,"error":"%s"}`, err.Error())), nil
		}
		return mcp.NewToolResultText(result), nil
	})
}

// marshalArgs 将参数 map 序列化为 JSON 字符串
func marshalArgs(args map[string]interface{}) string {
	if len(args) == 0 {
		return "{}"
	}
	jsonstr, _ := json.Marshal(args)
	return string(jsonstr)
}

// getDebugFlag 从环境变量读取 Debug 标志
func getDebugFlag() bool {
	return os.Getenv("MCP_DEBUG") == "true"
}

// Log 工具注册日志
func Log(format string, args ...interface{}) {
	log.Printf(format, args...)
}
