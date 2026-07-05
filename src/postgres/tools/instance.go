package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"postgres_server/security"
)

// RegisterInstanceTools 注册实例相关工具（12个：2现有迁移 + 4查询 + 6管理）
func RegisterInstanceTools(s *server.MCPServer, cred *common.Credential, g *security.Guard) {
	// ===== 现有迁移（2个）=====

	// DescribeDBInstanceAttribute - 查询实例详情（只读）
	registerTool(s, cred, g, "DescribeDBInstanceAttribute", "查询实例详情",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Description("实例ID")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeDBInstanceAttributeRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeDBInstanceAttribute(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// UpgradeDBInstanceKernelVersion - 升级实例内核版本号（写，L2业务确认）
	registerTool(s, cred, g, "UpgradeDBInstanceKernelVersion", "升级实例内核版本号",
		security.LevelBusiness,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Description("实例ID")),
			mcp.WithString("TargetDBKernelVersion", mcp.Description("升级的目标内核版本号")),
			mcp.WithNumber("SwitchTag", mcp.Description("指定切换时间。可选值: 0-立即切换, 1-维护时间切换, 2-指定时间切换")),
			mcp.WithString("SwitchStartTime", mcp.Description("切换开始时间，格式HH:MM:SS")),
			mcp.WithString("SwitchEndTime", mcp.Description("切换截止时间，格式HH:MM:SS")),
			mcp.WithBoolean("DryRun", mcp.Description("是否执行预检查")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewUpgradeDBInstanceKernelVersionRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.UpgradeDBInstanceKernelVersion(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// ===== 实例查询组（4个，只读）=====

	// DescribeDBInstances - 查询实例列表
	registerTool(s, cred, g, "DescribeDBInstances", "查询实例列表",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Description("实例ID")),
			mcp.WithArray("Filters", mcp.Description("过滤条件")),
			mcp.WithNumber("Limit", mcp.Description("每页返回数目，默认20")),
			mcp.WithNumber("Offset", mcp.Description("数据偏移量，从0开始")),
			mcp.WithString("OrderBy", mcp.Description("排序字段")),
			mcp.WithString("OrderByType", mcp.Description("排序方式: asc|desc")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			// 注意：DescribeDBInstances 接口没有 DBInstanceId 字段，仅支持 Filters 过滤，
			// 当设置了 InstanceScope 时，强制用 db-instance-id 过滤条件覆盖用户传入的 Filters，
			// 防止越权列出 scope 之外的实例。
			if g.InstanceScopeActive() {
				args["Filters"] = []map[string]interface{}{
					{"Name": "db-instance-id", "Values": []string{g.InstanceScope}},
				}
			}
			req := postgres.NewDescribeDBInstancesRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeDBInstances(req)
			if err != nil {
				return "", err
			}
			// 兜底二次过滤：即使服务端 Filters 未生效，也不把 scope 之外的实例返回给调用方
			if g.InstanceScopeActive() && rsp.Response != nil {
				filtered := make([]*postgres.DBInstance, 0, len(rsp.Response.DBInstanceSet))
				for _, inst := range rsp.Response.DBInstanceSet {
					if inst != nil && inst.DBInstanceId != nil && *inst.DBInstanceId == g.InstanceScope {
						filtered = append(filtered, inst)
					}
				}
				rsp.Response.DBInstanceSet = filtered
				total := uint64(len(filtered))
				rsp.Response.TotalCount = &total
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeClasses - 查询可用规格
	registerTool(s, cred, g, "DescribeClasses", "查询可用规格列表",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("Zone", mcp.Description("可用区ID")),
			mcp.WithString("DBEngine", mcp.Description("数据库引擎，默认postgresql")),
			mcp.WithString("DBMajorVersion", mcp.Description("数据库主版本号")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeClassesRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeClasses(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeDBVersions - 查询可用数据库版本
	registerTool(s, cred, g, "DescribeDBVersions", "查询可用数据库版本",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBEngine", mcp.Description("数据库引擎，默认postgresql")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeDBVersionsRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeDBVersions(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DescribeTasks - 查询异步任务状态
	registerTool(s, cred, g, "DescribeTasks", "查询异步任务状态",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Description("实例ID")),
			mcp.WithNumber("Limit", mcp.Description("每页返回数目")),
			mcp.WithNumber("Offset", mcp.Description("数据偏移量")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeTasksRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeTasks(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// ===== 实例管理组（6个，写操作需guard）=====

	// CreateInstances - 创建实例（L1费用确认）
	registerTool(s, cred, g, "CreateInstances", "创建实例",
		security.LevelFee,
		[]mcp.ToolOption{
			mcp.WithString("Zone", mcp.Required(), mcp.Description("可用区ID")),
			mcp.WithString("DBVersion", mcp.Description("数据库版本")),
			mcp.WithString("InstanceSpec", mcp.Description("实例规格")),
			mcp.WithString("InstanceName", mcp.Description("实例名称")),
			mcp.WithNumber("Volume", mcp.Description("磁盘容量(GB)")),
			mcp.WithNumber("Memory", mcp.Description("内存(GB)")),
			mcp.WithString("DBCharset", mcp.Description("数据库字符集")),
			mcp.WithString("InstanceChargeType", mcp.Description("计费类型: POSTPAID_BY_HOUR|PREPAID")),
			mcp.WithNumber("Period", mcp.Description("购买时长(月)")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewCreateInstancesRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.CreateInstances(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// ModifyDBInstanceName - 修改实例名称（L4审计）
	registerTool(s, cred, g, "ModifyDBInstanceName", "修改实例名称",
		security.LevelAudit,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithString("InstanceName", mcp.Required(), mcp.Description("新实例名称")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewModifyDBInstanceNameRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.ModifyDBInstanceName(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// ModifyDBInstanceSpec - 变更实例规格（L1费用确认）
	registerTool(s, cred, g, "ModifyDBInstanceSpec", "变更实例规格(扩缩容)",
		security.LevelFee,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithNumber("Memory", mcp.Description("内存(GB)")),
			mcp.WithNumber("Volume", mcp.Description("磁盘容量(GB)")),
			mcp.WithString("InstanceType", mcp.Description("实例类型")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewModifyDBInstanceSpecRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.ModifyDBInstanceSpec(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// RestartDBInstance - 重启实例（L2业务确认）
	registerTool(s, cred, g, "RestartDBInstance", "重启实例",
		security.LevelBusiness,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewRestartDBInstanceRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.RestartDBInstance(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// IsolateDBInstances - 隔离实例（L2业务确认）
	registerTool(s, cred, g, "IsolateDBInstances", "隔离实例",
		security.LevelBusiness,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewIsolateDBInstancesRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.IsolateDBInstances(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DisIsolateDBInstances - 解除隔离（L4审计）
	registerTool(s, cred, g, "DisIsolateDBInstances", "解除隔离实例",
		security.LevelAudit,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDisIsolateDBInstancesRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DisIsolateDBInstances(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	Log("Instance tools registered: 12")
}
