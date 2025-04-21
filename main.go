package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"12306-mcp-service/handlers"
	"12306-mcp-service/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// loadConfig 加载配置文件
func loadConfig() (*models.ServerConfig, error) {
	// 默认配置
	config := &models.ServerConfig{
		Port:    ":8080",
		Version: "1.0.0",
		Name:    "12306-MCP Server 🚀",
	}

	// 如果配置文件存在则加载
	if _, err := os.Stat("config.json"); err == nil {
		file, err := os.ReadFile("config.json")
		if err != nil {
			return nil, fmt.Errorf("读取配置文件失败: %v", err)
		}

		if err := json.Unmarshal(file, config); err != nil {
			return nil, fmt.Errorf("解析配置文件失败: %v", err)
		}
	}

	return config, nil
}

func main() {
	// 加载配置
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 创建MCP服务器
	s := server.NewMCPServer(
		config.Name,
		config.Version,
	)

	// 创建车票处理器
	ticketHandler := handlers.NewTicketHandler()

	// 添加车票查询工具
	ticketTool := mcp.NewTool("query_train_tickets",
		mcp.WithDescription("根据出发地、目的地、出发日期,查询列车车票信息"),
		mcp.WithString("from_station",
			mcp.Required(),
			mcp.Description("出发站"),
		),
		mcp.WithString("to_station",
			mcp.Required(),
			mcp.Description("到达站"),
		),
		mcp.WithString("date",
			mcp.Required(),
			mcp.Description("出发日期 (YYYY-MM-DD)"),
		),
	)

	// 添加工具处理器
	s.AddTool(ticketTool, ticketHandler.QueryTickets)

	log.Printf("启动服务器\n 名称: %s\n 版本: %s\n 端口: %s\n", config.Name, config.Version, config.Port)
	sseServer := server.NewSSEServer(s)
	if err := sseServer.Start(config.Port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
