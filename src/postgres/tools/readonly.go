package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"postgres_server/security"
)

// RegisterReadonlyTools 注册只读实例工具（2个，全新增）
func RegisterReadonlyTools(s *server.MCPServer, cred *common.Credential, g *security.Guard) {
	// DescribeReadOnlyGroups - 查询只读组列表（只读）
	registerTool(s, cred, g, "DescribeReadOnlyGroups", "查询只读组列表",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Description("主实例ID")),
			mcp.WithString("ReadOnlyGroupId", mcp.Description("只读组ID")),
			mcp.WithNumber("Limit", mcp.Description("每页返回数目")),
			mcp.WithNumber("Offset", mcp.Description("数据偏移量")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			// 注意：DescribeReadOnlyGroups 接口没有 DBInstanceId 字段，仅支持
			// Filters(db-master-instance-id) 过滤，当传入了 DBInstanceId（或被 guard
			// 自动注入 scope）时，强制转换为 Filters，防止越权列出其他主实例的只读组。
			if id, ok := args["DBInstanceId"].(string); ok && id != "" {
				args["Filters"] = []map[string]interface{}{
					{"Name": "db-master-instance-id", "Values": []string{id}},
				}
				delete(args, "DBInstanceId")
			}
			req := postgres.NewDescribeReadOnlyGroupsRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeReadOnlyGroups(req)
			if err != nil {
				return "", err
			}
			// 兜底二次过滤：scope 之外主实例的只读组不返回给调用方
			if g.InstanceScopeActive() && rsp.Response != nil {
				filtered := make([]*postgres.ReadOnlyGroup, 0, len(rsp.Response.ReadOnlyGroupList))
				for _, group := range rsp.Response.ReadOnlyGroupList {
					if group != nil && group.MasterDBInstanceId != nil && *group.MasterDBInstanceId == g.InstanceScope {
						filtered = append(filtered, group)
					}
				}
				rsp.Response.ReadOnlyGroupList = filtered
			}
			return rsp.ToJsonString(), nil
		})

	// CreateReadOnlyDBInstance - 创建只读实例（L1费用确认）
	registerTool(s, cred, g, "CreateReadOnlyDBInstance", "创建只读实例",
		security.LevelFee,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("主实例ID")),
			mcp.WithString("SpecName", mcp.Required(), mcp.Description("实例规格")),
			mcp.WithString("Zone", mcp.Description("可用区")),
			mcp.WithString("InstanceName", mcp.Description("只读实例名称")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewCreateReadOnlyDBInstanceRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.CreateReadOnlyDBInstance(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	Log("Readonly tools registered: 2")
}
