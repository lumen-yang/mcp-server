package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"postgres_server/security"
)

// RegisterNetworkTools 注册网络和安全组工具（4个，全新增）
func RegisterNetworkTools(s *server.MCPServer, cred *common.Credential, g *security.Guard) {
	// ===== 网络组（2个）=====

	// OpenDBExtranetAccess - 开启公网访问（L2业务确认，暴露公网）
	registerTool(s, cred, g, "OpenDBExtranetAccess", "开启实例公网访问",
		security.LevelBusiness,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithNumber("WanPort", mcp.Description("公网端口")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewOpenDBExtranetAccessRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.OpenDBExtranetAccess(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// CloseDBExtranetAccess - 关闭公网访问（L2业务确认，中断公网连接）
	registerTool(s, cred, g, "CloseDBExtranetAccess", "关闭实例公网访问",
		security.LevelBusiness,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewCloseDBExtranetAccessRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.CloseDBExtranetAccess(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// ===== 安全组组（2个）=====

	// DescribeDBInstanceSecurityGroups - 查询实例安全组（只读）
	registerTool(s, cred, g, "DescribeDBInstanceSecurityGroups", "查询实例安全组",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeDBInstanceSecurityGroupsRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeDBInstanceSecurityGroups(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// ModifyDBInstanceSecurityGroups - 修改实例安全组（L2业务确认，误改断连）
	registerTool(s, cred, g, "ModifyDBInstanceSecurityGroups", "修改实例安全组",
		security.LevelBusiness,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithArray("SecurityGroupIds", mcp.Required(), mcp.Description("安全组ID列表")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewModifyDBInstanceSecurityGroupsRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.ModifyDBInstanceSecurityGroups(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	Log("Network tools registered: 4")
}
