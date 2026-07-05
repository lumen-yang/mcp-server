package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"postgres_server/security"
)

// RegisterBackupTools 注册备份恢复工具（8个，全新增）
func RegisterBackupTools(s *server.MCPServer, cred *common.Credential, g *security.Guard) {
	// ===== 只读（5个）=====

	// DescribeBackupOverview - 查询备份概览（只读）
	registerTool(s, cred, g, "DescribeBackupOverview", "查询备份概览",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeBackupOverviewRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeBackupOverview(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeBaseBackups - 查询基础备份列表（只读）
	registerTool(s, cred, g, "DescribeBaseBackups", "查询基础备份列表",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Description("实例ID")),
			mcp.WithString("StartTime", mcp.Description("开始时间")),
			mcp.WithString("EndTime", mcp.Description("结束时间")),
			mcp.WithNumber("Limit", mcp.Description("每页返回数目")),
			mcp.WithNumber("Offset", mcp.Description("数据偏移量")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			// 注意：DescribeBaseBackups 接口没有 DBInstanceId 字段，仅支持 Filters 过滤，
			// 当设置了 InstanceScope 时，强制用 db-instance-id 过滤条件覆盖调用参数，
			// 防止越权列出 scope 之外其他实例的备份信息。
			if id, ok := args["DBInstanceId"].(string); ok && id != "" {
				args["Filters"] = []map[string]interface{}{
					{"Name": "db-instance-id", "Values": []string{id}},
				}
				delete(args, "DBInstanceId")
			}
			req := postgres.NewDescribeBaseBackupsRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeBaseBackups(req)
			if err != nil {
				return "", err
			}
			// 兜底二次过滤：scope 之外的备份记录不返回给调用方
			if g.InstanceScopeActive() && rsp.Response != nil {
				filtered := make([]*postgres.BaseBackup, 0, len(rsp.Response.BaseBackupSet))
				for _, b := range rsp.Response.BaseBackupSet {
					if b != nil && b.DBInstanceId != nil && *b.DBInstanceId == g.InstanceScope {
						filtered = append(filtered, b)
					}
				}
				rsp.Response.BaseBackupSet = filtered
				total := uint64(len(filtered))
				rsp.Response.TotalCount = &total
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeLogBackups - 查询日志备份列表（只读）
	registerTool(s, cred, g, "DescribeLogBackups", "查询日志备份列表",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Description("实例ID")),
			mcp.WithString("StartTime", mcp.Description("开始时间")),
			mcp.WithString("EndTime", mcp.Description("结束时间")),
			mcp.WithNumber("Limit", mcp.Description("每页返回数目")),
			mcp.WithNumber("Offset", mcp.Description("数据偏移量")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			// 注意：DescribeLogBackups 接口没有 DBInstanceId 字段，仅支持 Filters 过滤，
			// 当设置了 InstanceScope 时，强制用 db-instance-id 过滤条件覆盖调用参数，
			// 防止越权列出 scope 之外其他实例的备份信息。
			if id, ok := args["DBInstanceId"].(string); ok && id != "" {
				args["Filters"] = []map[string]interface{}{
					{"Name": "db-instance-id", "Values": []string{id}},
				}
				delete(args, "DBInstanceId")
			}
			req := postgres.NewDescribeLogBackupsRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeLogBackups(req)
			if err != nil {
				return "", err
			}
			// 兜底二次过滤：scope 之外的备份记录不返回给调用方
			if g.InstanceScopeActive() && rsp.Response != nil {
				filtered := make([]*postgres.LogBackup, 0, len(rsp.Response.LogBackupSet))
				for _, b := range rsp.Response.LogBackupSet {
					if b != nil && b.DBInstanceId != nil && *b.DBInstanceId == g.InstanceScope {
						filtered = append(filtered, b)
					}
				}
				rsp.Response.LogBackupSet = filtered
				total := uint64(len(filtered))
				rsp.Response.TotalCount = &total
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeAvailableRecoveryTime - 查询可恢复时间范围（只读）
	registerTool(s, cred, g, "DescribeAvailableRecoveryTime", "查询可恢复时间范围",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeAvailableRecoveryTimeRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeAvailableRecoveryTime(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeBackupDownloadURL - 获取备份下载链接（L4审计，链接泄露风险）
	registerTool(s, cred, g, "DescribeBackupDownloadURL", "获取备份下载链接",
		security.LevelAudit,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithString("BackupId", mcp.Required(), mcp.Description("备份ID")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeBackupDownloadURLRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeBackupDownloadURL(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// ===== 写操作（3个）=====

	// CreateBaseBackup - 创建基础备份（L2业务确认，消耗IO/存储）
	registerTool(s, cred, g, "CreateBaseBackup", "创建基础备份",
		security.LevelBusiness,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewCreateBaseBackupRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.CreateBaseBackup(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// CloneDBInstance - 克隆实例（L1费用确认）
	registerTool(s, cred, g, "CloneDBInstance", "克隆实例",
		security.LevelFee,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("源实例ID")),
			mcp.WithString("SpecName", mcp.Description("克隆实例规格")),
			mcp.WithString("Zone", mcp.Description("可用区")),
			mcp.WithString("InstanceName", mcp.Description("克隆实例名称")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewCloneDBInstanceRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.CloneDBInstance(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// RestoreDBInstanceObjects - 恢复数据（L3最高级确认，覆盖数据不可逆）
	registerTool(s, cred, g, "RestoreDBInstanceObjects", "恢复数据对象(覆盖)",
		security.LevelCritical,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("目标实例ID")),
			mcp.WithString("BackupId", mcp.Required(), mcp.Description("备份ID")),
			mcp.WithString("RestoreType", mcp.Required(), mcp.Description("恢复类型")),
			mcp.WithString("DBName", mcp.Description("数据库名")),
			mcp.WithString("TableName", mcp.Description("表名")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewRestoreDBInstanceObjectsRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.RestoreDBInstanceObjects(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	Log("Backup tools registered: 8")
}
