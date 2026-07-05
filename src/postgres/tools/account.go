package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"postgres_server/security"
)

// RegisterAccountTools 注册账号管理工具（6个：1现有迁移 + 5新增）
func RegisterAccountTools(s *server.MCPServer, cred *common.Credential, g *security.Guard) {
	// ===== 现有迁移（1个）=====

	// DescribeAccounts - 查询数据库账号列表（只读）
	registerTool(s, cred, g, "DescribeAccounts", "查询实例的数据库账号列表",
		security.LevelNone,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Description("实例ID，形如postgres-6fego161")),
			mcp.WithNumber("Limit", mcp.Description("每页返回数目，默认20，取值1-100")),
			mcp.WithNumber("Offset", mcp.Description("数据偏移量，从0开始")),
			mcp.WithString("OrderBy", mcp.Description("排序字段: createTime|name|updateTime")),
			mcp.WithString("OrderByType", mcp.Description("排序方式: desc|asc")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDescribeAccountsRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DescribeAccounts(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// ===== 新增（5个）=====

	// CreateAccount - 创建账号（L4审计）
	registerTool(s, cred, g, "CreateAccount", "创建数据库账号",
		security.LevelAudit,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithString("UserName", mcp.Required(), mcp.Description("账号名")),
			mcp.WithString("Password", mcp.Required(), mcp.Description("账号密码")),
			mcp.WithString("Remarks", mcp.Description("备注")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewCreateAccountRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.CreateAccount(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// DeleteAccount - 删除账号（L2业务确认，误删导致应用断连）
	registerTool(s, cred, g, "DeleteAccount", "删除数据库账号",
		security.LevelBusiness,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithString("UserName", mcp.Required(), mcp.Description("账号名")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewDeleteAccountRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.DeleteAccount(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// ModifyAccountPrivileges - 修改账号权限（L4审计，提权审计）
	registerTool(s, cred, g, "ModifyAccountPrivileges", "修改账号权限",
		security.LevelAudit,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithString("UserName", mcp.Required(), mcp.Description("账号名")),
			mcp.WithString("DBName", mcp.Required(), mcp.Description("数据库名")),
			mcp.WithString("Privileges", mcp.Required(), mcp.Description("权限: rw|r|ddl|owner")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewModifyAccountPrivilegesRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.ModifyAccountPrivileges(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// ResetAccountPassword - 重置账号密码（L4审计，旧连接断开）
	registerTool(s, cred, g, "ResetAccountPassword", "重置账号密码",
		security.LevelAudit,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithString("UserName", mcp.Required(), mcp.Description("账号名")),
			mcp.WithString("Password", mcp.Required(), mcp.Description("新密码")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewResetAccountPasswordRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.ResetAccountPassword(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	// LockAccount - 锁定账号（L2业务确认，导致应用报错）
	registerTool(s, cred, g, "LockAccount", "锁定数据库账号",
		security.LevelBusiness,
		[]mcp.ToolOption{
			mcp.WithString("DBInstanceId", mcp.Required(), mcp.Description("实例ID")),
			mcp.WithString("UserName", mcp.Required(), mcp.Description("账号名")),
		},
		func(client *postgres.Client, args map[string]interface{}) (string, error) {
			req := postgres.NewLockAccountRequest()
			req.FromJsonString(marshalArgs(args))
			rsp, err := client.LockAccount(req)
			if err != nil {
				return "", err
			}
			return rsp.ToJsonString(), nil
		})

	Log("Account tools registered: 6")
}
