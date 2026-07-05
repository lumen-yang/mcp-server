package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

const region = "ap-chengdu"
const instanceID = "postgres-1lbqykq6"

func main() {
	c, err := client.NewSSEMCPClient("http://127.0.0.1:9000/sse")
	if err != nil {
		fmt.Println("new client error:", err)
		os.Exit(1)
	}
	ctx := context.Background()
	if err := c.Start(ctx); err != nil {
		fmt.Println("start error:", err)
		os.Exit(1)
	}
	defer c.Close()

	_, err = c.Initialize(ctx, mcp.InitializeRequest{})
	if err != nil {
		fmt.Println("initialize error:", err)
		os.Exit(1)
	}

	call := func(name string, args map[string]interface{}) (string, error) {
		req := mcp.CallToolRequest{}
		req.Params.Name = "postgres-" + name
		req.Params.Arguments = args
		cctx, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()
		res, err := c.CallTool(cctx, req)
		if err != nil {
			return "", err
		}
		var sb strings.Builder
		for _, content := range res.Content {
			if tc, ok := content.(mcp.TextContent); ok {
				sb.WriteString(tc.Text)
			}
		}
		return sb.String(), nil
	}

	results := map[string]string{}
	order := []string{}
	record := func(name, out string, err error) {
		order = append(order, name)
		if err != nil {
			results[name] = "TRANSPORT_ERROR: " + err.Error()
			return
		}
		results[name] = out
	}

	base := map[string]interface{}{"region": region, "DBInstanceId": instanceID}
	noID := map[string]interface{}{"region": region}

	// 1-5 instance group
	out, err := call("DescribeDBInstanceAttribute", base)
	record("DescribeDBInstanceAttribute", out, err)

	out, err = call("DescribeDBInstances", noID)
	record("DescribeDBInstances", out, err)

	out, err = call("DescribeClasses", map[string]interface{}{"region": region, "DBEngine": "postgresql"})
	record("DescribeClasses", out, err)

	out, err = call("DescribeDBVersions", noID)
	record("DescribeDBVersions", out, err)

	out, err = call("DescribeTasks", base)
	record("DescribeTasks", out, err)

	// parameter group
	out, err = call("DescribeDBInstanceParameters", base)
	record("DescribeDBInstanceParameters", out, err)

	out, err = call("DescribeParameterTemplates", noID)
	record("DescribeParameterTemplates", out, err)
	templateID := extractFirst(out, "TemplateId")
	fmt.Println("DEBUG extracted templateID:", templateID)

	if templateID != "" {
		out, err = call("DescribeParameterTemplateAttributes", map[string]interface{}{"region": region, "TemplateId": templateID})
	} else {
		out, err = call("DescribeParameterTemplateAttributes", map[string]interface{}{"region": region, "TemplateId": "notfound"})
	}
	record("DescribeParameterTemplateAttributes", out, err)

	// ssl
	out, err = call("DescribeDBInstanceSSLConfig", base)
	record("DescribeDBInstanceSSLConfig", out, err)

	// account
	out, err = call("DescribeAccounts", base)
	record("DescribeAccounts", out, err)

	// network
	out, err = call("DescribeDBInstanceSecurityGroups", base)
	record("DescribeDBInstanceSecurityGroups", out, err)

	// monitoring
	now := time.Now()
	monArgs := map[string]interface{}{
		"region":       region,
		"DBInstanceId": instanceID,
		"StartTime":    now.Add(-24 * time.Hour).Format("2006-01-02 15:04:05"),
		"EndTime":      now.Format("2006-01-02 15:04:05"),
	}
	out, err = call("DescribeSlowQueryList", monArgs)
	record("DescribeSlowQueryList", out, err)

	out, err = call("DescribeSlowQueryAnalysis", monArgs)
	record("DescribeSlowQueryAnalysis", out, err)

	out, err = call("DescribeDBErrlogs", monArgs)
	record("DescribeDBErrlogs", out, err)

	// database
	out, err = call("DescribeDatabases", base)
	record("DescribeDatabases", out, err)
	dbName := extractFirst(out, "DBName")
	if dbName == "" {
		dbName = extractFirst(out, "DatabaseName")
	}

	if dbName != "" {
		out, err = call("DescribeDatabaseObjects", map[string]interface{}{"region": region, "DBInstanceId": instanceID, "DatabaseName": dbName, "ObjectType": "schema"})
	} else {
		out, err = call("DescribeDatabaseObjects", map[string]interface{}{"region": region, "DBInstanceId": instanceID, "DatabaseName": "postgres", "ObjectType": "schema"})
	}
	record("DescribeDatabaseObjects", out, err)

	// backup
	out, err = call("DescribeBackupOverview", base)
	record("DescribeBackupOverview", out, err)

	out, err = call("DescribeBaseBackups", base)
	record("DescribeBaseBackups", out, err)

	out, err = call("DescribeLogBackups", base)
	record("DescribeLogBackups", out, err)

	out, err = call("DescribeAvailableRecoveryTime", base)
	record("DescribeAvailableRecoveryTime", out, err)

	// readonly
	out, err = call("DescribeReadOnlyGroups", base)
	record("DescribeReadOnlyGroups", out, err)

	fmt.Println("========== RESULTS ==========")
	for _, name := range order {
		out := results[name]
		if len(out) > 500 {
			out = out[:500] + "...(truncated)"
		}
		fmt.Printf("\n--- %s ---\n%s\n", name, out)
	}
}

// extractFirst 从JSON文本里粗略提取第一个字段值(仅用于验证脚本内部串联参数,非生产代码)
func extractFirst(jsonText string, field string) string {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(jsonText), &raw); err != nil {
		return ""
	}
	resp, ok := raw["Response"].(map[string]interface{})
	if !ok {
		return ""
	}
	for _, v := range resp {
		arr, ok := v.([]interface{})
		if !ok || len(arr) == 0 {
			continue
		}
		first, ok := arr[0].(map[string]interface{})
		if !ok {
			continue
		}
		if val, ok := first[field].(string); ok {
			return val
		}
	}
	return ""
}
