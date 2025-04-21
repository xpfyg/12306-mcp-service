package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"12306-mcp-service/services"

	"github.com/mark3labs/mcp-go/mcp"
)

// TicketHandler 车票处理器
type TicketHandler struct {
	ticketService *services.TicketService
}

// NewTicketHandler 创建车票处理器实例
func NewTicketHandler() *TicketHandler {
	return &TicketHandler{
		ticketService: services.NewTicketService(),
	}
}

// QueryTickets 处理车票查询请求
func (h *TicketHandler) QueryTickets(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 获取并验证参数
	fromStation, toStation, date, err := h.validateAndExtractParams(request)
	if err != nil {
		return nil, err
	}

	// 调用服务层查询车票
	tickets, err := h.ticketService.QueryTickets(fromStation, toStation, date)
	if err != nil {
		return nil, fmt.Errorf("获取车票数据失败: %v", err)
	}

	// 转换为JSON字符串
	result, err := json.MarshalIndent(tickets, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("序列化车票数据失败: %v", err)
	}

	return mcp.NewToolResultText(string(result)), nil
}

// validateAndExtractParams 验证并提取请求参数
func (h *TicketHandler) validateAndExtractParams(request mcp.CallToolRequest) (fromStation, toStation, date string, err error) {
	fromStation, ok := request.Params.Arguments["from_station"].(string)
	if !ok || fromStation == "" {
		return "", "", "", errors.New("出发站不能为空")
	}

	toStation, ok = request.Params.Arguments["to_station"].(string)
	if !ok || toStation == "" {
		return "", "", "", errors.New("到达站不能为空")
	}

	date, ok = request.Params.Arguments["date"].(string)
	if !ok || date == "" {
		return "", "", "", errors.New("出发日期不能为空")
	}

	// 验证日期格式
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return "", "", "", fmt.Errorf("日期格式无效，应为YYYY-MM-DD: %v", err)
	}

	return fromStation, toStation, date, nil
}
